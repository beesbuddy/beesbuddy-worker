BINARY=beesbuddy-worker

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

swag:
	swag init;

build:
	swag init;
	go build -o ${BINARY} .;

dev:
	 air -c dev.air.default.toml worker serve

token:
	go run main.go make token e641c5f30441812f79130ff0518fbff2

.PHONY: clean
