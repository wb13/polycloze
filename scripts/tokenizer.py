"""Tokenizes sentences from standard input and outputs CSV files."""

from argparse import ArgumentParser, Namespace
from collections import Counter
import csv
from dataclasses import dataclass
import fileinput
from importlib import import_module
import json
from pathlib import Path
import sys
import typing as t

from spacy.language import Language

from .language import languages


def load_spacy_language(code: str) -> Language:
    if code not in languages:
        sys.exit("unknown language code")
    parent, name = languages[code].spacy_path
    mod = import_module(f"spacy.lang.{parent}")
    return t.cast(Language, getattr(mod, name)())


@dataclass
class Tokenizer:
    nlp: Language

    def tokenize(self, sentence: str) -> list[str]:
        tokens = []
        for token in self.nlp.tokenizer(sentence):
            tokens.append(token.text)
            if token.whitespace_:
                tokens.append(token.whitespace_)
        return tokens


class Sentence(t.NamedTuple):
    id: int | None
    text: str
    tokens: list[str]

    def __hash__(self) -> int:
        return hash(self.text)

    def row(self) -> tuple[str, str] | tuple[int, str, str]:
        if self.id is None:
            return (self.text, json.dumps(self.tokens))
        return (self.id, self.text, json.dumps(self.tokens))


class WordCounter:
    def __init__(self) -> None:
        self.counter: Counter[str] = Counter()

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
        "-f",
        dest="file",
        type=Path,
        help="input file (default: stdin)",
    )
    parser.add_argument(
        "-o",
        dest="output",
        help="output directory",
        required=True,
    )
    parser.add_argument(
        "-l",
        dest="log",
        help="log non-words",
    )
    parser.add_argument(
        "--no-ids",
        dest="has_ids",
        help="input has no IDs",
        action="store_false",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    log = Path(args.log) if args.log is not None else None
    if log:
        log.parent.mkdir(parents=True, exist_ok=True)
        log.write_text("", encoding="utf-8")

    output = Path(args.output)
    if output.is_file():
        sys.exit(f"{args.output} is a file")
    output.mkdir(parents=True, exist_ok=True)

    tokenizer = Tokenizer(load_spacy_language(args.language))
    word_counter = WordCounter()

    with open(output/"sentences.csv", "w", encoding="utf-8") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["tatoeba_id", "text", "tokens"])
        with fileinput.input(files=args.file or "-") as file:
            for line in file:
                id_ = None
                if args.has_ids:
                    id_, line = line.split("\t", maxsplit=1)

                assert id_
                sentence = Sentence(
                    id=int(id_),
                    text=line,
                    tokens=tokenizer.tokenize(line),
                )
                word_counter.add(sentence.tokens)
                writer.writerow(sentence.row())

    try:
        language = languages[args.language]
    except KeyError:
        sys.exit(f"unsupported language: {args.language}")
    with open(output/"words.csv", "w", newline="", encoding="utf-8") as file:
        writer = csv.writer(file)
        writer.writerow(["word", "frequency"])
        for row in word_counter.count():
            if language.is_word(row[0]):
                writer.writerow(row)
            elif log:
                with open(log, "a", encoding="utf-8") as logfile:
                    print(row[0], file=logfile)


if __name__ == "__main__":
    main(parse_args())
