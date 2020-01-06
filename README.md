# go-delve-reload

### The Ultimate Go and React Development Setup with Docker

This repository is paired with a [blog post](https://blog.ivorscott.com/ultimate-go-development-with-docker). If you follow along, the project starter is available under the `starter` branch.

## Contents

- VSCode Setup
- Multi-stage Builds
- Docker Compose
- Using Makefiles
- Postgres
- Traefik
- Live Reloading
- Debugging
- Testing

### Folder structure

```
├── .vscode
|  ├── launch.json
├── README.md
├── api
|  ├── Dockerfile
|  ├── cmd
|  |  └── api
|  |     └── main.go
|  ├── go.mod
|  ├── go.sum
|  ├── internal
|  |  ├── api
|  |  |  ├── client.go
|  |  |  ├── handlers.go
|  |  |  ├── handlers_test.go
|  |  |  ├── helpers.go
|  |  |  ├── middleware.go
|  |  |  └── routes.go
|  |  └── models
|  |     ├── models.go
|  |     └── postgres
|  |        └── products.go
|  ├── pkg
|  |  └── secrets
|  |     └── secrets.go
|  └── scripts
|     └── create-db.sh
├── client
|  ├── Dockerfile
|  ├── README.md
|  ├── package-lock.json
|  ├── package.json
|  ├── public
|  └── src
|     ├── App.css
|     ├── App.js
|     ├── App.test.js
|     ├── index.css
|     ├── index.js
|     ├── logo.svg
|     ├── serviceWorker.js
|     └── setupTests.js
├── docker-compose.yml
├── makefile
└── secrets
   ├── postgres_db
   ├── postgres_host
   ├── postgres_passwd
   └── postgres_user
```

### Usage

1 - Create a secrets folder in the project root.
Add the following secret files:

```
└── secrets
   ├── postgres_db
   ├── postgres_host
   ├── postgres_passwd
   └── postgres_user
```

In each file add your secret value.

2 - Add the following domains to your machine's /etc/hosts file

```
127.0.0.1       client.local api.local debug.api.local traefik.api.local pgadmin.local
```

3 - In a terminal, and under the project root, execute `make`.

#### Commands

```makefile
make # launch fullstack app (frontend/backend)

make api # develop api with live reload

make debug-api # use delve on the same api in a separate container (no live reload)

make debug-db # use pgcli to inspect postgres db

make dump # create a db backup

make exec user="..." service="..." cmd="..." # execute command in running container

make tidy # clean up unused api dependencies

make test-api # run api tests

make test-client # run client tests

make down # tear down all containers

```

#### Using the debugger in VSCode

Run the debuggable api. Set a break point on a route handler. Click 'Launch remote' then visit the route in the browser.

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

**Note:**

_Always replicate the production environment as much as possible. We do this locally by using self-signed certificates automated by Traefik._

_In your browser, you may see the "Your connection is not private" message. This is common when using self-signed certificates._

_Simply click "Advanced", and then "Proceed to ... (unsafe)"._
