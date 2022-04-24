#!/usr/bin/env python

from argparse import ArgumentParser, Namespace

class LanguageLoader:
    @staticmethod
    def deu():
        from spacy.lang.de import German
        return German()

    @staticmethod
    def eng():
        from spacy.lang.en import English
        return English()

    @staticmethod
    def fra():
        from spacy.lang.fr import French
        return French()

    @staticmethod
    def ita():
        from spacy.lang.it import Italian
        return Italian()

    @staticmethod
    def spa():
        from spacy.lang.es import Spanish
        return Spanish()


languages = [name for name in dir(LanguageLoader) if not name.startswith("_")]


def load_language(code: str):
    language = getattr(LanguageLoader, code, None)
    if not language:
        sys.exit("unknown language code")
    return language()


def main() -> None:
    for lang in languages:
        print(lang)


if __name__ == "__main__":
    main()
