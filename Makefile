.PHONY:	build
build:
	go build -tags sqlite_math_functions

.PHONY:	run
run:
	go run -tags sqlite_math_functions .

.PHONY:	test
test:
	cd database; go test -tags sqlite_math_functions
	cd review_scheduler; go test -tags sqlite_math_functions
