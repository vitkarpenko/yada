COMMIT := $(shell git rev-parse --short HEAD)

build:
	docker build -t yada:$(COMMIT) .

run: build
	docker run --restart=always --name=yada -d yada:$(COMMIT)

stop:
	docker kill yada && docker rm yada

pull:
	git pull

update: pull build stop run