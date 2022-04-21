COMMIT := $(shell git rev-parse --short HEAD)

build:
	docker build -t yada:$(COMMIT) .

run: build
	docker run --restart=always -d yada:$(COMMIT)

stop:
	docker kill yada:$(COMMIT)

pull:
	git pull

update: pull build stop run