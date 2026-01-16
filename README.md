# go-delve-reload

### Go and React Development with Docker pt.1

This repository is paired with a [blog post](https://medium.com/gitconnected/go-and-react-development-with-docker-pt-1-9f734a5991fc). If you follow along, the project starter is available under the `starter` branch.

[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G01SD6G)

## Contents

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
   ├── postgres_passwd
   └── postgres_user
```

### Requirements

- Docker
- VSCode

### Usage

1 - Create a secrets folder in the project root.
Add the following secret files:

```
└── secrets
   ├── postgres_db
   ├── postgres_passwd
   └── postgres_user
```

In each file add your secret value.

2 - Add the following domains to your machine's /etc/hosts file

```
127.0.0.1       client.local api.local debug.api.local traefik.api.local pgadmin.local
```

3 - In a terminal, and under the project root, execute `make`.

4 - Navigate to https://api.local/products and https://client.local in two separate tabs.

**Note:**

_To replicate the production environment as much as possible locally, we use self-signed certificates automated by Traefik._

_In your browser, you may see a warning and need to click a link to proceed to the requested page. This is common when using self-signed certificates._

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
