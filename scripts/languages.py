#!/usr/bin/env python

from argparse import ArgumentParser, Namespace

from spacy.language import Language


def deu() -> Language:
    from spacy.lang.de import German
    return German()


def fra() -> Language:
    from spacy.lang.fr import French
    return French()


def spa() -> Language:
    from spacy.lang.es import Spanish
    return Spanish()


languages = {
    "deu": deu,
    "fra": fra,
    "spa": spa,
}


def load_language(code: str) -> Language:
    language = languages.get(code)
    if not language:
        sys.exit("unknown language code")
    return language()


def main() -> None:
    for lang in languages:
        print(lang)


if __name__ == "__main__":
    main()
