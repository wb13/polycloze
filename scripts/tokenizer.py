#!/usr/bin/env python

"""Tokenizes sentences from standard input and outputs CSV files."""

from argparse import ArgumentParser
from collections import Counter
import csv
from dataclasses import dataclass
from functools import reduce
import json
from pathlib import Path
import sys
import typing as t

from spacy.language import Language
from spacy.lang.de import German
from spacy.lang.es import Spanish


@dataclass
class Tokenizer:
    nlp: Language

    def tokenize(self, sentence: str) -> list[str]:
        tokens = (
            [token.text for token in self.nlp.tokenizer(word)]
            for word in sentence.split()
        )
        return reduce(lambda x, y: x + [" "] + y, tokens)


def load_language(code: str) -> Language:
    match code:
        case "deu":
            return German()
        case "spa":
            return Spanish()
        case _:
            sys.exit("unknown language code")


def count_words(sentences: dict[str, list[str]]) -> Counter:
    counter = Counter()
    for _, tokens in sentences.items():
        counter.update(tokens)
    return counter


def parse_args():
    parser = ArgumentParser()
    parser.add_argument("language", help="ISO 639-3 language code")
    return vars(parser.parse_args())


def main() -> None:
    args = parse_args()
    tokenizer = Tokenizer(load_language(args.get("language")))
    sentences = {}

    try:
        while line := input():
            sentences[line] = tokenizer.tokenize(line)
    except EOFError:
        pass

    with open("words.csv", "w", newline="") as csvfile:
        writer = csv.writer(csvfile)
        for word, count in count_words(sentences).most_common():
            writer.writerow([word, count])

    with open("sentences.csv", "w", newline="") as csvfile:
        writer = csv.writer(csvfile)
        for sentence, tokens in sentences.items():
            writer.writerow([sentence, json.dumps(tokens)])


if __name__ == "__main__":
    main()
