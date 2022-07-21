#!/usr/bin/env python

import sys


class LanguageLoader:
    @staticmethod
    def cyo():
        from spacy.lang.tl import Tagalog
        return Tagalog()

    @staticmethod
    def deu():
        from spacy.lang.de import German
        return German()

    @staticmethod
    def eng():
        from spacy.lang.en import English
        return English()

    @staticmethod
    def spa():
        from spacy.lang.es import Spanish
        return Spanish()

    @staticmethod
    def tgl():
        from spacy.lang.tl import Tagalog
        return Tagalog()


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
