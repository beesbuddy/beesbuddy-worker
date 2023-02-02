# beesbuddy-worker

Application responsible for dispatching messages from MQTT to time series database and WS

## How to generate swagger documentation

First of all you need to install swaggo with:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

Then you need to run:

```sh
swag init
```

## How to run app with air

```sh
air -c dev.air.default.toml web serve
```

If in configuration `isProd` is specified as false, hot reload is initialized and hot reload is initialize via `npm run hot` command.