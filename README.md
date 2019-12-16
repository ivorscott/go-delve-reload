# go-delve-reload

### Setup a Go development environment in Docker

#### Commands

```makefile
make api # develop api with live reload

make debug-api # use delve on the same api in a separate container (no live reload)
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

TODO: Live reload while debugging
