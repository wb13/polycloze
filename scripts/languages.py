#!/usr/bin/env python

from argparse import ArgumentParser, Namespace


def deu():
    from spacy.lang.de import German
    return German()


def fra():
    from spacy.lang.fr import French
    return French()


def ita():
    from spacy.lang.it import Italian
    return Italian()


def spa():
    from spacy.lang.es import Spanish
    return Spanish()


languages = {
    "deu": deu,
    "fra": fra,
    "ita": ita,
    "spa": spa,
}


def load_language(code: str):
    language = languages.get(code)
    if not language:
        sys.exit("unknown language code")
    return language()


def main() -> None:
    for lang in languages:
        print(lang)


if __name__ == "__main__":
    main()
