"""Tokenizes sentences from standard input and outputs CSV files."""

from argparse import ArgumentParser, Namespace
from collections import Counter
import csv
from dataclasses import dataclass
import fileinput
import json
from pathlib import Path
import sys
import typing as t

from .language import languages, Language

if t.TYPE_CHECKING:
    from spacy.tokenizer import Tokenizer as SpacyTokenizer     # type: ignore


@dataclass
class Tokenizer:
    tokenizer: "SpacyTokenizer"

    def tokenize(self, sentence: str) -> list[str]:
        tokens = []
        for token in self.tokenizer(sentence):
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


def write_sentences(
    outfile: Path,
    infile: Path | None,
    tokenizer: Tokenizer,
    word_counter: WordCounter,
) -> None:
    """Write tokenized sentences to output file.

    infile: file containing list of sentences.
    Pass None to get sentences from stdin.
    """
    with (
        open(outfile, "w", encoding="utf-8") as csvfile,
        fileinput.input(files=infile or "-") as file,
    ):
        writer = csv.writer(csvfile)
        writer.writerow(["tatoeba_id", "text", "tokens"])
        for line in file:
            id_, line = line.split("\t", maxsplit=1)
            line = line.strip()
            sentence = Sentence(
                id=int(id_),
                text=line,
                tokens=tokenizer.tokenize(line),
            )
            word_counter.add(sentence.tokens)
            writer.writerow(sentence.row())


def write_words(
    output: Path,
    word_counter: WordCounter,
    language: Language,
    log: Path,
) -> None:
    """log: where to write nonwords."""
    with (
        open(output, "w", newline="", encoding="utf-8") as file,
        open(log, "w", encoding="utf-8") as logfile,
    ):
        writer = csv.writer(file)
        writer.writerow(["word", "frequency"])
        for row in word_counter.count():
            if language.is_word(row[0]):
                writer.writerow(row)
            else:
                print(row[0], file=logfile)


def process_language(
    language_code: str,
    output: Path,
    file: Path | None = None,
) -> None:
    """Tokenize sentences in file and write all necessary outputs.

    output: where to write files
    file: input file of sentences, or stdin if value is None
    """
    output.mkdir(parents=True, exist_ok=True)

    language = languages[language_code]
    tokenizer = Tokenizer(language.tokenizer())
    word_counter = WordCounter()
    write_sentences(output/"sentences.csv", file, tokenizer, word_counter)
    write_words(
        output/"words.csv",
        word_counter,
        language,
        output/"nonwords.txt",
    )


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
        type=Path,
        required=True,
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    if args.language not in languages:
        sys.exit(f"unsupported language: {args.language}")
    if args.output.is_file():
        sys.exit(f"{args.output} is a file")

    process_language(args.language, args.output, args.file)


if __name__ == "__main__":
    main(parse_args())
