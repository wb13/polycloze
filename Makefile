languages = $(shell python -m scripts.languages)
latest_sentences = $(shell find build/tatoeba/sentences.*.csv | sort -r | head -n 1)
latest_links = $(shell find build/tatoeba/links.*.csv | sort -r | head -n 1)

define add_language
.PHONY:	$(1)
$(1):	build/sqlite/$(1).db

build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv	&:	build/sentences/$(1).tsv
	python -m scripts.tokenizer $(1) -o build/languages/$(1) < $$<

build/languages/$(1)/words.txt build/languages/$(1)/non-words.txt	&:	build/languages/$(1)/words.csv
	python python/blacklist/blacklist/uncsv.py $$< | PYTHONPATH=python/blacklist python -m blacklist $(1) \
		-b build/languages/$(1)/non-words.txt \
		-w build/languages/$(1)/words.txt

build/sqlite/$(1).db:	build/languages/$(1)/non-words.txt build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv
	mkdir -p build/sqlite
	rm -f $$@
	./scripts/check-migrations.sh migrations/
	./scripts/migrate.sh $$@ migrations/
	./scripts/importer.py $$@ -i $$< \
		-s build/languages/$(1)/sentences.csv \
		-w build/languages/$(1)/words.csv
endef

.PHONY:	all
all:	build/translations.db

build/ids.txt:	$(latest_sentences)
	languages=$$(printf "${languages}" | tr '[:space:]' '|'); \
	./scripts/sentences.sh $${languages} | ./scripts/format.sh id > $@

build/translations.csv:	build/ids.txt build/subset build/symmetric
	./build/subset $< < $(latest_links) | ./build/symmetric > $@

build/translations.db:	build/translations.csv
	./scripts/make-translation-db.sh $< $@

$(foreach lang,$(languages),$(eval $(call add_language,$(lang))))

.PHONY:	download
download:
	./scripts/download.sh
	python ./scripts/partition.py ./build/sentences < $(latest_sentences)


### nim stuff

build/subset build/symmetric:	build/%:	src/%.nim
	nim c -o:$@ --stackTrace:off --checks:off --opt:speed $<
