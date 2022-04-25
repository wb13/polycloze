languages = $(shell python -m scripts.languages)
latest_sentences = $(shell find build/tatoeba/sentences.*.csv | sort -r | head -n 1)
latest_links = $(shell find build/tatoeba/links.*.csv | sort -r | head -n 1)

define add_language
build/sentences/$(1).txt:	$(latest_sentences)
	./scripts/sentences.sh $(1) | ./scripts/format.sh id-sentence > $$@

dist/$(1).tar.gz:	build/sentences/$(1).txt
	python -m scripts.tokenizer $(1) -o $$@ < $$<
endef

.PHONY:	all
all:

.PHONY:	dist
dist:	$(patsubst %,dist/%.tar.gz,$(languages)) dist/translations.tar.gz

build/ids.txt:	$(latest_sentences)
	languages=$$(printf "${languages}" | tr '[:space:]' '|'); \
	./scripts/sentences.sh $${languages} | ./scripts/format.sh id > $@

build/translations.tsv:	build/ids.txt build/subset build/symmetric
	./build/subset $< < $(latest_links) | ./build/symmetric > $@

dist/translations.tar.gz:	build/translations.tsv
	cd build; tar -czvf translations.tar.gz translations.tsv
	mv build/translations.tar.gz dist

$(foreach lang,$(languages),$(eval $(call add_language,$(lang))))


### nim stuff

build/subset build/symmetric:	build/%:	src/%.nim
	nim c -o:$@ --stackTrace:off --checks:off --opt:speed $<
