#!/usr/bin/env python

import sys


class LanguageLoader:
    from spacy.language import Language

    @staticmethod
    def cyo() -> Language:
        from spacy.lang.tl import Tagalog
        return Tagalog()

    @staticmethod
    def deu() -> Language:
        from spacy.lang.de import German
        return German()

    @staticmethod
    def eng() -> Language:
        from spacy.lang.en import English
        return English()

    @staticmethod
    def spa() -> Language:
        from spacy.lang.es import Spanish
        return Spanish()

    @staticmethod
    def tgl() -> Language:
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
