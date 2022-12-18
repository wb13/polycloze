DOCKER = $(shell command -v podman || command -v docker)

.PHONY:	all
all:	build lint test

.PHONY:	build-js
build-js:
	cd api/js; npm run build

.PHONY:	build
build:	build-js
	go build .
	go build -v -o build/ ./...

.PHONY:	test
test:	build-js
	go test -cover ./...

.PHONY:	format
format:
	gofmt -s -w .

.PHONY:	bench
bench:
	cd flashcards; go test -cpuprofile ../build/cpu.prof -bench .
	go tool pprof build/cpu.prof

.PHONY:	lint-js
lint-js:
	cd api/js; npm run check

.PHONY:	lint-go
lint-go:	format
	golangci-lint run

.PHONY:	lint
lint:	lint-js | lint-go

.PHONY:	run
run:	build-js
	go run .

.PHONY:	init
init:
	cd api/js; npm ci

.PHONY:	docker-build
docker-build:
	$(DOCKER) build -t polycloze-demo .

.PHONY:	docker-run
docker-run:
	$(DOCKER) run -p 3000:3000 --cpus 1 --memory 256m polycloze-demo
