# go-delve-reload

### Setup a Go development environment in Docker

This repository is paired with a [blog post](https://blog.ivorscott.com/ultimate-go-development-with-docker). If you follow along, the project starter is available under the `starter` branch.

## Contents

- VSCode Setup
- Multi-stage Builds
- Docker Compose
- Postgres
- Traefik
- React
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

#### Commands

```makefile
make # launch fullstack app (frontend/backend)

make api # develop api with live reload

make test-api # run api tests

make debug-api # use delve on the same api in a separate container (no live reload)

make debug-db # use pgcli to inspect postgres db

make dump # create a db backup

make exec user="..." service="..." cmd="..." # execute command in running container (user defaults to root)

make down # tear down all containers

make tidy # clean up unused api dependencies

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
