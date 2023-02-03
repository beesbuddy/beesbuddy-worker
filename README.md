# Beesbuddy worker

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

If in configuration `HotReload` is set to true, app will try to run ui with hot reload via `npm run hot` command. I found it usefull at specific cases for me.

You can alway run ui with hot reload by using make or just executing `npm run hot`.
