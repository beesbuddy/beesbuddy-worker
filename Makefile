APP_NAME = beesbuddy-worker

clean:
	if [ -f ${APP_NAME} ] ; then rm ${APP_NAME} ; fi

swag:
	swag init;

build:
	swag init;
	go build -o ${APP_NAME} .;

build-image:
	docker build -t ${APP_NAME} .

dev:
	air -c dev.air.default.toml worker serve

token:
	go run main.go make token e641c5f30441812f79130ff0518fbff2

health:
	curl -X GET -H "Accept: application/json" http://localhost:4000/metrics

.PHONY: clean
