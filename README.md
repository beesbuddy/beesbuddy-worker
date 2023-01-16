# beesbuddy-worker

Application responsible for dispatching messages from MQTT to time series database and WS

## Hot to generate swagger documentation

First of all you need to install swaggo with:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

Then you need to run:

```sh
swag init
```
