FROM golang:1.19 as build
LABEL MAINTAINER="Viktor Nareiko"

WORKDIR /go/src/beesbuddy-worker

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo .

FROM alpine:latest as release

WORKDIR /app

RUN mkdir /data
COPY dev.default.json dev.default.json
COPY staging.default.json /data/staging.default.json
COPY --from=build /go/src/beesbuddy-worker/beesbuddy-worker .

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/beesbuddy-worker

EXPOSE 4000
EXPOSE 1883
EXPOSE 443
EXPOSE 80

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["./beesbuddy-worker", "worker", "serve"]