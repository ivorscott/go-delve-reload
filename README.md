# go-delve-reload

![Minion](docs/demo.png)

## The Ultimate Go and React Development Setup with Docker (Part 2)

### Building A Better API

This repository is paired with a [blog post](https://blog.ivorscott.com/ultimate-go-react-development-setup-with-docker-part2). Many improvements in Part 2 originate from [Ardan labs service training](https://github.com/ardanlabs/service-training). I highly recommend their [courses](https://education.ardanlabs.com/).

[Previous blog post](https://blog.ivorscott.com/ultimate-go-react-development-setup-with-docker)

## Contents

- Improvements to Part 1
- Graceful Shutdown
- Seeding & Migrations (With Go-Migrate)
- Package Oriented Design
- Fluent SQL Generation (With Squirrel)
- Error Handling
- Cancellation
- Request Validation
- Request Logging
- Integration Testing (With TestContainers-Go)

### Improvements to Part 1

<details>
  <summary>See changes</summary>

  <br/>

1 - Removed Traefik from development

If you recall, the previous post used Traefik for self-signed certificates. We're no longer using Traefik in development. We'll use it in production. `create-react-app` and the `net/http` packages both have mechanisms to use self signed-certificates. This cleans up our docker-compose file and speeds up the workflow. Now we don't need to pull the Traefik image or run the container. In the client app we enable self-signed certificates by adding `HTTPS=true` to the `package.json`.

```json
// package.json

"scripts": {
    "start": "HTTPS=true node scripts/start.js",
    "build": "node scripts/build.js",
    "test": "node scripts/test.js"
  },
```

In Go, we use the `crypto` package to generate a cert with `make cert`. Running make also works because cert is the first target in the makefile and thus the default. This is intentional. Someone might think executing `make` initializes the project (like in Part 1). In that case, they end up generating required API certs instead. This makes the Makefile usage less error-prone. Generating certs more than once replaces existing certs without issue.

```makefile
# makefile

cert:
	mkdir -p ./api/tls
	@go run $(GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	@mv *.pem ./api/tls
```

The following demonstrates how we can switch between self-signed certificates and Traefik. When `cfg.Web.Production` is true, we are using Traefik. In Part 3 ("Docker Swarm and Traefik"), we will have a separate compose file for production.

```go
// main.go

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		if cfg.Web.Production {
			serverErrors <- api.ListenAndServe()
		} else {
			serverErrors <- api.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		}
	}()
```

2 - Cleaner terminal logging

Using self-signed certificates produces ugly logs.

![Minion](docs/ugly-self-signed-cert-logs.png)

We can avoid the `tls: unknown certificate` message by disabling server error logging. It's ok to do this in development. The things we do care about print from logging and error middleware. When `cfg.Web.Production` is false, a new error logger will discard server logs.

```go
// main.go

	var errorLog *log.Logger

	if !cfg.Web.Production {
		// Prevent the HTTP server from logging stuff on its own.
		// The things we care about we log ourselves.
		// Prevents "tls: unknown certificate" errors caused by self-signed certificates.
		errorLog = log.New(ioutil.Discard, "", 0)
	}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      c.Handler(mux),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		ErrorLog:     errorLog,
	}
```

[CompileDaemon](https://github.com/githubnemo/CompileDaemon) created ugly logs as well. CompileDaemon prefixes all child process output with stdout or stderr labels.

![Minion](docs/ugly-compile-daemon-logs.png)

```yaml
# docker-compose.yml

command: CompileDaemon --build="go build -o main ./cmd/api" -log-prefix=false --command=./main
```

3 - Added the Ardan Labs configuration package

In Part 1, API configuration came from environment variables in the docker-compose file. But the were dependent on docker secret values, making it harder to opt out of docker in development. Reserve docker secrets for production and adopt the [Ardan Labs configuration package](https://github.com/ardanlabs/conf). The package supports both environment variables and command line arguments. Now we can out out of docker if we want a more idiomatic Go API development workflow. I copied and paste the package under: `/api/internal/platform/conf`.

The struct field `cfg.Web.Production` can in cli form would be `--web-production`. In environment variable form it is `API_WEB_PRODUCTION`. Notice, as an environment variable there's an extra namespace. This ensures we only parse the vars we expect. This also reduces name conflicts. In our case that namespace is `API`.

```go
// main.go

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address            string        `conf:"default:localhost:4000"`
			Production         bool          `conf:"default:false"`
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:5s"`
			ShutdownTimeout    time.Duration `conf:"default:5s"`
			FrontendAddress    string        `conf:"default:https://localhost:3000"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("API", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

```

The configuration package requires a nested struct describing the configuration fields. Each field has a type and default value supplied in a struct tag. To parse the arguments, as environment variables or the command line flags we do: `conf.Parse(os.Args[1:], "API", &cfg)`. If there's an error we either reveal usage instructions or throw a fatal error.The next snippet shows the same vars referenced in our compose file with the API namespace:

```yaml
# docker-compose.yml

services:
  api:
    build:
      context: ./api
      target: dev
    container_name: api
    environment:
      CGO_ENABLED: 0
      API_DB_HOST: db
      API_WEB_PRODUCTION: "false"
      API_WEB_ADDRESS: :$API_PORT
      API_WEB_READ_TIMEOUT: 7s
      API_WEB_WRITE_TIMEOUT: 7s
      API_WEB_SHUTDOWN_TIMEOUT: 7s
      API_WEB_FRONTEND_ADDRESS: https://localhost:$CLIENT_PORT
    ports:
      - $API_PORT:$API_PORT
```

Store default environment variables in a `.env` file. When an `.env` file exists in the same directory as the docker-compose file, we can reference it. To do this, prefix a dollar sign before the environment variable name. For example: `$API_PORT` or `$CLIENT_PORT`. This allows us to maintain default values in a separate file for docker-compose.

4 - Removed Docker Secrets from Development

Docker secrets are a Swarm specific construct. They aren't secret in docker-compose anyway [PR #4368](https://github.com/docker/compose/pull/4368). This only works because docker-compose isn't complaining when it sees them. Now Docker secrets are only supported when `cfg.Web.Production` is true. When this happens we swap out the default database configuration with secrets.

```go
// main.go

	// =========================================================================
	// Enabled Docker Secrets

	if cfg.Web.Production {
		dockerSecrets, err := secrets.NewDockerSecrets()
		if err != nil {
			log.Fatalf("error : retrieving docker secrets failed : %v", err)
		}

		cfg.DB.Name = dockerSecrets.Get("postgres_db")
		cfg.DB.User = dockerSecrets.Get("postgres_user")
		cfg.DB.Host = dockerSecrets.Get("postgres_host")
		cfg.DB.Password = dockerSecrets.Get("postgres_passwd")
	}

```

More on Docker secrets when we get to production (discussed in Part 3, _"Docker Swarm and Traefik"_).

5 - Removed PgAdmin4

PgAdmin4 is one of many Postgres editors available. For example, I've enjoyed using SQLPro Studio at work. If you're going to use PgAdmin4 or any other editor, use it on your host machine without a container. Reason being, it's more reliable. Importing and exporting sql files is difficult in a PgAdmin4 container.

6 - Enabled Idiomatic Go development

Containerizing the Go API is now optional. This makes our development workflow even more flexible. This tweet made me consider the consequences of having the API too coupled to Docker:

"Folks, keep docker out of your edit/compile/test inner loop."

-- https://twitter.com/davecheney/status/1232078682287591425

In the end Docker should be optional and you should know your reasons for using it. My reasons are:

1. Custom Workflows
2. Predictability Across Machines
3. Isolated Environments
4. Optional Live Reloading
5. Optional Delve Debugging
6. Integration Testing In CI
7. Preparation For Deployments

</details>

### Requirements

- VSCode
- Postman
- Docker

## Getting Started

```
git clone https://github.com/ivorscott/go-delve-reload
cd go-delve-reload
git checkout part2
```

[Setup VSCode](https://blog.ivorscott.com/ultimate-go-react-development-setup-with-docker#setting-up-vscode)

### Usage

1 - Copy .env.sample and rename it to .env

The contents of .env should look like this:

```bash
# DEVELOPMENT ENVIRONMENT VARIABLES

API_PORT=4000
CLIENT_PORT=3000

POSTGRES_DB=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=db
POSTGRES_NET=postgres-net
```

That's it. No further modifications required.

2 - Unblock port 5432 for postgres

Both the Makefile and docker-compose file reference the standard Postgres port: 5432. Before continuing, close any existing Postgres connections.
For example, I have a pre-existing installation of Postgres installed with Homebrew. Executing `brew info postgresql@10` generates info on how to start/stop the service. If I didn't know what version I installed I would run `brew list`.
With Homebrew I do:

```
pg_ctl -D /usr/local/var/postgresql@10 stop
killall postgresql
```

3 - Create self-signed certificates

The next command moves generated certificates to the `./api/tls/` directory.

```makefile
make cert
```

4 - Setup up the Postgres container

Run the database in the background.

```makefile
make db
```

#### Create your first migration

Make a migration to create the products table.

```makefile
make migration create_products_table
```

Add sql to both `up` & `down` migrations files found at: `./api/internal/schema/migrations/`.

**Up**

```sql
-- 000001_create_products_table.up.sql

CREATE TABLE products (
    id UUID not null unique,
    name varchar(100) not null,
    price real not null,
    description varchar(100) not null,
    created timestamp without time zone default (now() at time zone 'utc')
);
```

**Down**

```sql
-- 000001_create_products_table.down.sql

DROP TABLE IF EXISTS products;
```

#### Create your second migration

Make another migration to add tags to products:

```
make migration add_tags_to_products
```

**Up**

```sql

-- 000002_add_tags_to_products.up.sql

ALTER TABLE products
ADD COLUMN tags varchar(255);
```

**Down**

```sql
-- 000002_add_tags_to_products.down.sql

ALTER TABLE products
DROP Column tags;
```

Migrate up to the latest migration

```makefile
make up # you can migrate down with "make down"
```

Display which version you have selected. Expect it two print `2` since you created 2 migrations.

```makefile
make version
```

[Learn more about my go-migrate Postgres helper](https://github.com/ivorscott/go-migrate-postgres-helper)

#### Seeding the database

Create a seed file of the appropriate name matching the table name you wish to seed.

```makefile
make seed products
```

This adds an empty products.sql seed file found under `./api/internal/schema/seeds`. Add the following sql content:

```sql
-- ./api/internal/schema/seeds/products.sql

INSERT INTO products (id, name, price, description, created) VALUES
('cbef5139-323f-48b8-b911-dc9be7d0bc07','Xbox One X', 499.00, 'Eighth-generation home video game console developed by Microsoft.','2019-01-01 00:00:01.000001+00'),
('ce93a886-3a0e-456b-b7f5-8652d2de1e8f','Playstation 4', 299.00, 'Eighth-generation home video game console developed by Sony Interactive Entertainment.','2019-01-01 00:00:01.000001+00'),
('faa25b57-7031-4b37-8a89-de013418deb0','Nintendo Switch', 299.00, 'Hybrid console that can be used as a stationary and portable device developed by Nintendo.','2019-01-01 00:00:01.000001+00')
ON CONFLICT DO NOTHING;
```

Conflicts may arise when you execute the seed file more than once to the database. Appending "ON CONFLICT DO NOTHING;" to the end prevents this. This functionality depends on at least one table column having a unique constraint. In our case id is unique.

Finally, add the products seed to the database.

```
make insert products
```

Enter the database and examine its state.

```makefile
make debug-db
```

![Minion](docs/debug-db.png)

If the database gets deleted, you don't need to repeat every instruction. Simply run:

```
make db
make up
make insert products
```

5 - Run the api and client containers

#### Run the Go API container with live reload enabled

`make api`

#### Run the React TypeScript app container

`make client`

![Minion](docs/run.png)

First, navigate to the API in the browser at: <https://localhost:4000/v1/products>.

Then navigate to the client app at: <https://localhost:3000> in a separate tab.
This approach to development uses containers entirely.

**Note:**

To replicate the production environment as much as possible locally, we use self-signed certificates.

In your browser, you may see a warning and need to click a link to proceed to the requested page. This is common when using self-signed certificates.

6 - **Optional Idiomatic Go development** (container free Go API)

Another approach is to containerize only the client and database. Work with the API in an idiomatic fashion. This means without a container and with live reloading disabled. To configure the API, use command line flags or export environment variables.

```makefile
export API_DB_DISABLE_TLS=true
cd api
go run ./cmd/api
# go run ./cmd/api --db-disable-tls=true
```

### Try it in Postman

**List products**

GET <https://localhost:4000/v1/products>

**Retrieve one product**

GET <https://localhost:4000/v1/products/:id>

**Create a product**

POST <https://localhost:4000/v1/products>

```
{
	"name": "Game Cube",
	"price": 74,
	"description": "The GameCube is the first Nintendo console to use optical discs as its primary storage medium.",
	"tags": null
}
```

**Update a product**

PUT <https://localhost:4000/v1/products/:id>

```
{
	"name": "Nintendo Rich!"
}
```

**Delete a product**

DELETE <https://localhost:4000/v1/products/:id>

#### Commands

```makefile

make api # develop api with live reload

make cert # generate self-signed certificates

make client # develop client react app

make db # start the database in the background

make debug-api # use delve on the same api in a separate container (no live reload)

make debug-db # use pgcli to inspect postgres db

make rm # remove all containers

make rmi # remove all images

make exec user="..." service="..." cmd="..." # execute command in running container

make tidy # clean up unused api dependencies

make test-api # run api tests

make test-client # run client tests

make migration <name> # create a migration

make version # print current migration version

make up <number> # migrate up a number (optional number, defaults to latest migration)

make down <number> # migrate down a number (optional number, defaults to 1)

make force <version> # Set version but don't run migration (ignores dirty state)

make seed <name> # create seed filename

make insert <name> # insert seed file to database
```

#### Using the debugger in VSCode

If you wish to debug with Delve you can do this in a separate container instance on port 8888 automatically.

```makefile
make debug-api
```

Set a break point on a route handler. Click 'Launch remote' then visit the route in the browser.

[Read previous tutorial about delve debugging](https://blog.ivorscott.com/ultimate-go-react-development-setup-with-docker#delve-debugging-a-go-api)

#### VSCode launch.json

```
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch remote",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "cwd": "${workspaceFolder}/api",
      "remotePath": "/api",
      "port": 2345,
      "showLog": true,
      "trace": "verbose"
    }
  ]
}

```

## The Ultimate Go and React Series

### Building A Workflow

[The Ultimate Go and React Development Setup with Docker (Part 1)](https://blog.ivorscott.com/ultimate-go-react-development-setup-with-docker)

<details>

  <summary>See content</summary>

  <br/>

- VSCode Setup
- Docker Basics
- Multi-stage Builds
- Docker Compose
- Using Makefiles
- Using Postgres
- Using Traefik
- Live Reloading a Go API
- Delve Debugging a Go API
- Testing

</details>

### Building A Better API

The Ultimate Go and React Development Setup with Docker (Part 2)

<details>
  <summary>See content</summary>

  <br/>

- Initial Changes From Part 1
- Graceful Shutdown
- Seeding & Migrations (With Go-Migrate)
- Package Oriented Design
- Fluent SQL Generation (With Squirrel)
- Error Handling
- Cancellation
- Request Validation
- Request Logging
- Integration Testing (With TestContainers-Go)

</details>

### Security and Awareness: OAuth, Observability, And Profiling

The Ultimate Go and React Development Setup with Docker (Part 3)

<details>
  <summary>See content</summary>

  <br/>

- Health checks
- Profiling
- Open Telemetry
- OAuth & Auth0
- Authentication
- Authorization

</details>

### Docker Swarm and Traefik

The Ultimate Go and React Production Setup with Docker (Part 4)

<details>
  <summary>See content</summary>

  <br/>
  
- Docker Hub
- Docker Swarm
- Swarm Secrets
- Digital Ocean
- Managed Databases
- Volume Storage Plugins
- Traefik

</details>

### Continuous Integration And Continuous Delivery

The Ultimate Go and React Production Setup with Docker (Part 5)

<details>
  <summary>See content</summary>

  <br/>

- Drone CI
- Portainer

</details>
