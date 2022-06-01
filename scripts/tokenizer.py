"""Tokenizes sentences from standard input and outputs CSV files."""

from argparse import ArgumentParser, Namespace
from collections import Counter
import csv
from dataclasses import dataclass
from functools import reduce
import json
from pathlib import Path
import sys
import typing as t

from spacy.language import Language

from .languages import load_language


@dataclass
class Tokenizer:
    nlp: Language

    def tokenize(self, sentence: str) -> list[str]:
        tokens = (
            [token.text for token in self.nlp.tokenizer(word)]
            for word in sentence.split()
        )
        return reduce(lambda x, y: x + [" "] + y, tokens)


class Sentence(t.NamedTuple):
    id: int | None
    text: str
    tokens: list[str]

    def __hash__(self) -> int:
        return hash(self.text)

    def row(self) -> tuple[str, str]:
        if self.id is None:
            return (self.text, json.dumps(self.tokens))
        return (self.id, self.text, json.dumps(self.tokens))


class WordCounter:
    def __init__(self):
        self.counter = Counter()

    def add(self, tokens: t.Iterable[str]) -> None:
        self.counter.update(token.casefold() for token in tokens)

    def count(self) -> list[tuple[str, int]]:
        return self.counter.most_common()


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "language",
        help="ISO 639-3 language code",
    )
    parser.add_argument(
        "-o",
        dest="output",
        help="output directory",
        required=True,
    )
    parser.add_argument(
        "--no-ids",
        dest="has_ids",
        help="input has no IDs",
        action="store_false",
    )
    return parser.parse_args()


def write_csv(
    path: Path | str,
    rows: t.Iterable[t.Sequence[t.Any]],
    *,
    header: t.Optional[t.Sequence[str]] = None,
) -> None:
    with open(path, "w", newline="") as csvfile:
        writer = csv.writer(csvfile)
        if header:
            writer.writerow(header)
        for row in rows:
            writer.writerow(row)


def main() -> None:
    args = parse_args()

    output = Path(args.output)
    if output.is_file():
        sys.exit(f"{args.output} is a file")
    output.mkdir(parents=True, exist_ok=True)

    tokenizer = Tokenizer(load_language(args.language))
    word_counter = WordCounter()

    with open(output/"sentences.csv", "w") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["tatoeba_id", "text", "tokens"])

        try:
            while line := input():
                id_ = None
                if args.has_ids:
                    # TODO handle exception
                    id_, line = line.split("\t", maxsplit=1)
                sentence = Sentence(
                    id=id_,
                    text=line,
                    tokens=tokenizer.tokenize(line),
                )
                word_counter.add(sentence.tokens)
                writer.writerow(sentence.row())
        except EOFError:
            pass

    write_csv(
        output/"words.csv",
        word_counter.count(),
        header=["word", "frequency"],
    )


if __name__ == "__main__":
    main()
