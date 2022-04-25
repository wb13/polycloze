languages = $(shell python -m scripts.languages)
latest = $(shell find build/tatoeba/sentences.*.csv | sort -r | head -n 1)

define add_language
build/sentences/$(1).txt:	$(latest)
	./scripts/sentences.sh $(1) | ./scripts/format.sh id-sentence > $$@

dist/$(1).tar.gz:	build/sentences/$(1).txt
	python -m scripts.tokenizer $(1) -o $$@ < $$<
endef

.PHONY:	all
all:

.PHONY:	dist
dist:	$(patsubst %,dist/%.tar.gz,$(languages))

build/ids.txt:	$(latest)
	@for lang in $(languages); \
	do \
		./scripts/sentences.sh $${lang} | ./scripts/format.sh id > $@; \
	done

$(foreach lang,$(languages),$(eval $(call add_language,$(lang))))
