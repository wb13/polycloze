#!/usr/bin/env python

"""Tokenizes sentences from standard input and outputs CSV files."""

from argparse import ArgumentParser, Namespace
from collections import Counter
import csv
from dataclasses import dataclass
from functools import reduce
import json
from pathlib import Path
import sys
import tarfile
from tempfile import TemporaryDirectory
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


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "language",
        help="ISO 639-3 language code",
    )
    parser.add_argument(
        "-o",
        dest="output",
        help="output file",
        required=True,
    )
    return parser.parse_args()


def write_csv(path: Path | str, rows: t.Iterable[t.Sequence[t.Any]]) -> None:
    with open(path, "w", newline="") as csvfile:
        writer = csv.writer(csvfile)
        for row in rows:
            writer.writerow(row)


def main() -> None:
    args = parse_args()
    tokenizer = Tokenizer(load_language(args.language))
    sentences = {}

    try:
        while line := input():
            sentences[line] = tokenizer.tokenize(line)
    except EOFError:
        pass

    with TemporaryDirectory() as tmpdirname:
        tmpdir = Path(tmpdirname)
        words_csv = tmpdir/"words.csv"
        sentences_csv = tmpdir/"sentences.csv"

        write_csv(words_csv, count_words(sentences).most_common())
        write_csv(
            sentences_csv,
            ([sentence, json.dumps(tokens)] for sentence, tokens in sentences.items()),
        )

        with tarfile.open(args.output, "w:gz") as tar:
            tar.add(tmpdir, arcname=args.language)


if __name__ == "__main__":
    main()
