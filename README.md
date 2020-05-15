# The Go and React Series

## Deploying with Swarm and Traefik pt.7

## Contents

- Digital Ocean
- Docker Hub
- Docker Machine
- Docker Swarm
- Healthchecks
- Traefik
- Deployment

### Requirements

- VSCode
- Docker
- DockerHub Account
- Digital Ocean Account
- Auth0 Account
- A Domain Name (Try Namecheap)
- Docker Machine
- Managed Database (ElephantSQL is Free)

## Getting Started

```
git clone https://github.com/ivorscott/go-delve-reload
cd go-delve-reload
git checkout part7
```

[Setup VSCode](https://blog.ivorscott.com/ultimate-go-react-development-setup-with-docker#setting-up-vscode)

### Usage

1 - Copy .env.sample and rename it to .env

The contents of .env should look like this:

```bash
# ENVIRONMENT VARIABLES

API_PORT=4000
CLIENT_PORT=3000

AUTH0_DOMAIN=
AUTH0_AUDIENCE=
AUTH0_CLIENT_ID=

POSTGRES_DB=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=db
POSTGRES_NET=postgres-net

REACT_APP_BACKEND=https://localhost:4000/v1
API_WEB_FRONTEND_ADDRESS=https://localhost:3000
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

## The Go and React Series

### Go and React Development with Docker pt.1

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

### Transitioning to Go pt.2

<details>

  <summary>See content</summary>

  <br/>

- Why Go?
- Challenges

</details>

### Building an API with Go pt.3

<details>
  <summary>See content</summary>

  <br/>

- Package Oriented Design
- Configuration
- Database Connection
- Docker Secrets
- Graceful Shutdown
- Middleware
- Handling Requests
- Error Handling
- Seeding & Migrations (With Go-Migrate)
- Integration Testing (With TestContainers-Go)

</details>

### My API Workflow with Go pt.4

<details>
  <summary>See content</summary>

  <br/>

- A Demo
- Profiling

</details>

### OAuth 2 with Auth0 in Go pt.5

<details>
  <summary>See content</summary>

  <br/>

- OAuth
- Auth0
- Authentication
- Authorization

</details>

### Observability Metrics in Go pt.6

<details>
  <summary>See content</summary>

  <br/>

- Open Telemetry
- Prometheus
- Grafana

</details>

### Deploying with Swarm and Traefik pt.7

<details>
  <summary>See content</summary>

  <br/>
  
- Digital Ocean
- Docker Hub
- Docker Machine
- Docker Swarm
- Healthchecks
- Traefik
- Deployment

</details>

### CICD with Portainer and Drone pt.8

<details>
  <summary>See content</summary>

  <br/>

- Drone CI
- Portainer

</details>
