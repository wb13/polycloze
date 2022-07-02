packages = api cmd/api
l1 = eng
l2 = spa

.PHONY:	all
all:

.PHONY:	build
build:	$(packages)

.PHONY:	$(packages)
$(packages):	%:
	cd $@; go build

.PHONY:	test
test:
	go test -cover ./...

.PHONY:	format
format:
	gofmt -s -w .

.PHONY:	run
run:	build
	cd cmd/api; ./api $(l1) $(l2)
