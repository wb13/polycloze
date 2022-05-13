.PHONY:	build
build:
	go build -tags sqlite_math_functions

.PHONY:	run
run:
	go run -tags sqlite_math_functions .
