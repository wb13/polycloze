.PHONY:	all
all:

.PHONY:	build
build:
	cd api/js; npm run build
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
	golangci-lint run

.PHONY:	run
run:	build
	./polycloze

.PHONY:	init
init:
	cd api/js; npm ci
