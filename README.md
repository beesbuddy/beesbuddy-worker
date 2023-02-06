# Beesbuddy worker

Application responsible for dispatching messages from MQTT to time series database and broadcast provider. It have simple rest api to perform settings for workers (adding, removing subscribers) on the fly.

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
## How to run app with make

```sh
make dev
```
