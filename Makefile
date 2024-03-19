COMMIT := $(shell git rev-parse --short HEAD)

build:
	docker build -t yada:$(COMMIT) .

generate:
	sqlc generate

run: build stop
	docker run \
		--restart=always --name=yada \
		-v $(PWD)/data:/yada/data \
		-d yada:$(COMMIT)

fmt/code:
	gofumpt -l -w .

fmt/imports:
	gci write --skip-generated -s standard -s 'prefix(github.com/vitkarpenko)' -s default --custom-order .

stop:
	docker kill yada || true && docker rm yada || true

pull:
	git pull

update: pull build stop run

deploy:
	ssh vitkarpenko@pi 'cd yada; make update'
