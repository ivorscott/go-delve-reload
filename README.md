# go-delve-reload

### Setup a Go development environment in Docker

This repository is paired with a [blog post](https://blog.ivorscott.com/ultimate-go-development-with-docker). If you follow along, the project starter is available under the `starter` branch.

## Contents

- VSCode Setup
- Multi-stage Builds
- Docker Compose
- Traefik
- Live Reloading
- Debugging

#### Commands

```makefile
make api # develop api with live reload

make debug-api # use delve on the same api in a separate container (no live reload)

make api-d # tear down all containers

make test # run unit tests

make run cmd="..." # run command in new api container

make exec cmd="..." # execute command in existing api container

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
