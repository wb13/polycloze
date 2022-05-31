packages = api cmd/api

.PHONY:	all
all:

.PHONY:	build
build:	$(packages)

.PHONY:	$(packages)
$(packages):	%:
	cd $@; go build -tags sqlite_math_functions

.PHONY:	test
test:
	cd database; go test -tags sqlite_math_functions
	cd review_scheduler; go test -tags sqlite_math_functions
	cd translator; go test -tags sqlite_math_functions

.PHONY:	format
format:
	@for package in $(packages); do \
		(cd $$package; gofmt -s -w .); \
	done

.PHONY:	run
run:	build
	cd cmd/api; ./api
