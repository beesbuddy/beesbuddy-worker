BINARY=beesbuddy-worker
STATIC_DIR=static

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	cd ${STATIC_DIR} && find . -type f ! -name 'embed.go' -delete && rm -rf css js

static: clean
	npm install; npm run prod

hot: npm install
	npm run hot

build: static
	go build -o ${BINARY} .;

.PHONY: clean static