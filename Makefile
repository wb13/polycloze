packages = database flashcards review_scheduler sentence_picker translator word_queue word_scheduler \
					 cmd/flashcards cmd/sentence_picker cmd/srs cmd/translator

.PHONY:	build
build:
	@for package in $(packages); do \
		(cd $$package; go build -tags sqlite_math_functions); \
	done

.PHONY:	test
test:
	cd database; go test -tags sqlite_math_functions
	cd review_scheduler; go test -tags sqlite_math_functions
