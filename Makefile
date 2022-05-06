.PHONY:	build
build:
	go build -tags sqlite_math_functions

.PHONY:	run
run:
	go run -tags sqlite_math_functions .

.PHONY:	test
test:
	cd srs; go test -tags sqlite_math_functions

.PHONY:	clean
clean:
	rm polycloze-srs test.db
