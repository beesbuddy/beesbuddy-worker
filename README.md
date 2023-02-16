# Beesbuddy worker

Application responsible for dispatching messages from MQTT to time series database and broadcast provider. It have simple rest api to perform settings for workers (adding, removing subscribers) on the fly.

## How to generate swagger documentation

First of all you need to install swaggo and air with:

```bash
$ go install github.com/swaggo/swag/cmd/swag@latest
$ go install github.com/cosmtrek/air@latest
```

Then you need to run:

```bash
$ swag init
```

## How to run app with air

```bash
$ air -c dev.air.default.toml web serve
```
## How to run app with make

```bash
$ make dev
```

## How to build

Change client secret from `e641c5f30441812f79130ff0518fbff2` to something else. Warning `e641c5f30441812f79130ff0518fbff2` must be used only for local development.

Apply changes to your needs in dev.default.json

Run make command:

```bash
$ make build
```

## TODO

* Implement message handler logic (sending metrics to influx db);
