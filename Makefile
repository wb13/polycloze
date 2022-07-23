languages = $(shell python -m scripts.language)
pairs = $(foreach l1,$(languages), $(foreach l2,$(languages), $(l1)-$(l2)))
latest_sentences = $(shell find build/tatoeba/sentences.*.csv | sort -r | head -n 1)
latest_links = $(shell find build/tatoeba/links.*.csv | sort -r | head -n 1)

define add_language
.PHONY:	$(1)
$(1):	build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv

build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv	&:	build/sentences/$(1).tsv
	python -m scripts.tokenizer $(1) -o build/languages/$(1) < $$<
endef

define add_pair
.PHONY:	$(1)-$(2)
$(1)-$(2):	build/courses/$(1)-$(2).db

build/translations/$(1)-$(2).csv:	build/sentences/$(1).tsv build/sentences/$(2).tsv $$(latest_links)
	if [[ "$(1)" < "$(2)" ]]; then \
		mkdir -p build/translations; \
		python -m scripts.mapper $$^ > $$@; \
	fi

build/courses/$(1)-$(2).db:	build/translations/$(1)-$(2).csv $(1) $(2)
	mkdir -p build/courses
	rm -f $$@
	if [[ "$(1)" != "$(2)" ]]; then \
		./scripts/check-migrations.sh migrations/; \
		./scripts/migrate.sh $$@ migrations/; \
	fi
	if [[ "$(1)" < "$(2)" ]]; then \
		python -m scripts.populate -r $$@ build/languages/$(1) build/languages/$(2) $$<; \
	fi
	if [[ "$(2)" < "$(1)" ]]; then \
		python -m scripts.populate $$@ build/languages/$(1) build/languages/$(2) build/translations/$(2)-$(1).csv; \
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
	mkdir -p "$$HOME/.local/share/polycloze"
	cp build/courses/*.db "$$HOME/.local/share/polycloze"

# For applying migrations to existing build files
.PHONY:	migrate
migrate:
	./scripts/check-migrations.sh migrations/
	for course in ./build/courses/*.db; do \
		./scripts/migrate.sh "$$course" migrations/; \
	done

.PHONY:	check
check:
	pylint scripts -d C0115,C0116
	flake8 --max-complexity 10 scripts
	mypy --strict scripts
