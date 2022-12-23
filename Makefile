BINARY=beesbuddy-worker
PUBLIC_DIR=public

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -d ${PUBLIC_DIR} ] ; then rm -rf ${PUBLIC_DIR} ; fi

static: clean
	cd ui;\
	npm install;\
	npm run build

build: static
	go build -o ${BINARY} .;

.PHONY: clean static