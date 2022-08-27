languages = $(shell python -m scripts.language)
pairs = $(foreach l1,$(languages), $(foreach l2,$(languages), $(l1)-$(l2)))

define add_language
.PHONY:	$(1)
$(1):	build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv

build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv	&:	build/sentences/$(1).tsv
	mkdir -p build/logs/nonwords
	python -m scripts.tokenizer $(1) -o build/languages/$(1) -l build/logs/nonwords/$(1).txt < $$<
endef

define add_pair
.PHONY:	$(1)-$(2)
$(1)-$(2):	build/courses/$(1)-$(2).db

build/translations/$(1)-$(2).csv:	build/sentences/$(1).tsv build/sentences/$(2).tsv build/tatoeba/links.csv
	mkdir -p build/translations
	if [[ "$(1)" -ge "$(2)" ]]; then \
		touch $$@; \
	fi
	if [[ "$(1)" < "$(2)" ]]; then \
		mkdir -p build/translations; \
		python -m scripts.mapper $$^ > $$@; \
	fi

build/courses/$(1)-$(2).db:	build/translations/$(1)-$(2).csv build/languages/$(1)/sentences.csv build/languages/$(1)/words.csv build/languages/$(2)/sentences.csv build/languages/$(2)/words.csv
	mkdir -p build/courses
	rm -f $$@
	./scripts/check-migrations.sh migrations/; \
	./scripts/migrate.sh $$@ migrations/; \
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

build/tatoeba/links.csv build/tatoeba/sentences.csv:
	python -m scripts.download
	python -m scripts.untar

build/sentences/:	build/tatoeba/sentences.csv build/tatoeba/links.csv
	python -m scripts.partition $@ < $<

.PHONY:	install
install:
	mkdir -p "$$HOME/.local/share/polycloze"
	for course in ./build/courses/*.db; do \
		if $$(python -m scripts.course "$$course"); then \
			echo "Installing $$course"; \
			cp "$$course" "$$HOME/.local/share/polycloze"; \
		fi \
	done

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
	flake8 --max-complexity 12 scripts
	mypy --strict scripts
