.PHONY:	all
all:
	@echo 'Try: python -m scripts.build -h'

.PHONY:	install
install:
	cp -r ./build/polycloze "$$HOME/.local/share"

.PHONY:	check
check:
	pylint scripts -d C0115,C0116
	flake8 --max-complexity 12 scripts
	mypy --strict scripts

.PHONY:	init
init:
	pip install --upgrade pip
	pip install -r requirements/dev.requirements.txt
