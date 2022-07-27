.PHONY:	all
all:

.PHONY:	build-js
build-js:
	cd api/js; npm run build

.PHONY:	build
build:	build-js
	go build .
	go build -v -o build/ ./...

.PHONY:	test
test:
	go test -cover ./...

.PHONY:	format
format:
	gofmt -s -w .

.PHONY:	bench
bench:
	cd flashcards; go test -cpuprofile ../build/cpu.prof -bench .
	go tool pprof build/cpu.prof

.PHONY:	lint
lint:
	cd api/js; npm run check
	golangci-lint run

.PHONY:	run
run:	build-js
	go build .
	./polycloze

.PHONY:	init
init:
	cd api/js; npm ci
