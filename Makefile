languages = $(shell python -m scripts.languages)
pairs = $(foreach l1,$(languages), $(foreach l2,$(languages), $(l1)-$(l2)))
latest_sentences = $(shell find build/tatoeba/sentences.*.csv | sort -r | head -n 1)
latest_links = $(shell find build/tatoeba/links.*.csv | sort -r | head -n 1)

define add_language
.PHONY:	$(1)
$(1):	build/sqlite/$(1).db

build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv	&:	build/sentences/$(1).tsv
	python -m scripts.tokenizer $(1) -o build/languages/$(1) < $$<

build/sqlite/$(1).db:	build/languages/$(1)/non-words.txt build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv
	mkdir -p build/sqlite
	rm -f $$@
	./scripts/check-migrations.sh migrations/languages/
	./scripts/migrate.sh $$@ migrations/languages/
	./scripts/importer.py $$@ -i $$< \
		-s build/languages/$(1)/sentences.csv \
		-w build/languages/$(1)/words.csv
	python -m scripts.metadata $$@
endef

define add_pair
.PHONY:	$(1)-$(2)
$(1)-$(2):	build/courses/$(1)-$(2).db

build/translations/$(1)-$(2).csv:	build/sentences/$(1).tsv build/sentences/$(2).tsv $$(latest_links)
	if [[ "$(1)" < "$(2)" ]]; then \
		mkdir -p build/translations; \
		python -m scripts.mapper $$^ > $$@; \
	fi

build/courses/$(1)-$(2).db:	build/translations/$(1)-$(2).csv
	mkdir -p build/courses
	rm -f $$@
	./scripts/check-migrations.sh migrations/courses/
	./scripts/migrate.sh $$@ migrations/courses/
	if [[ "$(1)" < "$(2)" ]]; then \
		python -m scripts.populate $$@ build/languages/$(1) build/languages/$(2) $$<; \
	fi
	if [[ "$(2)" < "$(1)" ]]; then \
		python -m scripts.populate $$@ build/languages/$(1) build/languages/$(2) -r build/translations/$(2)-$(1).csv; \
	fi
endef

.PHONY:	all
all:	$(pairs) $(languages)

$(foreach lang,$(languages),$(eval $(call add_language,$(lang))))
$(foreach l1,$(languages),$(foreach l2,$(languages), $(eval $(call add_pair,$(l1),$(l2)))))

.PHONY:	download
download:
	./scripts/download.sh
	python ./scripts/partition.py ./build/sentences < $(latest_sentences)

.PHONY:	install
install:
	mkdir -p "$$HOME/.local/share/polycloze/languages"
	mkdir -p "$$HOME/.local/share/polycloze/translations"
	cp build/sqlite/*.db "$$HOME/.local/share/polycloze/languages"
	cp build/translations/*.db "$$HOME/.local/share/polycloze/translations"
