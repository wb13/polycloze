.PHONY:	all
all:
	@echo 'Try: python -m scripts.build -h'

.PHONY:	install
install:
	mkdir -p "$$HOME/.local/share/polycloze"
	cp ./build/courses/*.db "$$HOME/.local/share/polycloze"

.PHONY:	check
check:
	pylint scripts -d C0115,C0116
	flake8 --max-complexity 12 scripts
	mypy --strict scripts
