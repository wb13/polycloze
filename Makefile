packages = api buffer database flashcards review_scheduler sentence_picker translator word_queue word_scheduler \
					 cmd/api cmd/buffer cmd/flashcards cmd/sentence_picker cmd/srs cmd/translator

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
