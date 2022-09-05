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

from spacy.language import Language as SpacyLanguage

from .language import languages, Language


def load_spacy_language(code: str) -> SpacyLanguage:
    assert code in languages

    parent, name = languages[code].spacy_path
    mod = import_module(f"spacy.lang.{parent}")
    return t.cast(SpacyLanguage, getattr(mod, name)())


@dataclass
class Tokenizer:
    nlp: SpacyLanguage

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
    log: Path | None,
) -> None:
    if log:
        log.parent.mkdir(parents=True, exist_ok=True)
        log.write_text("", encoding="utf-8")

    with open(output, "w", newline="", encoding="utf-8") as file:
        writer = csv.writer(file)
        writer.writerow(["word", "frequency"])
        for row in word_counter.count():
            if language.is_word(row[0]):
                writer.writerow(row)
            elif log:
                with open(log, "a", encoding="utf-8") as logfile:
                    print(row[0], file=logfile)


def process_language(
    language_code: str,
    output: Path,
    file: Path | None = None,
    log: Path | None = None,
) -> None:
    """Tokenize sentences in file and write all necessary outputs.

    output: where to write files
    file: input file of sentences, or stdin if value is None
    log: optional log file for non-words
    """
    output.mkdir(parents=True, exist_ok=True)

    tokenizer = Tokenizer(load_spacy_language(language_code))
    word_counter = WordCounter()
    write_sentences(output/"sentences.csv", file, tokenizer, word_counter)
    write_words(
        output/"words.csv",
        word_counter,
        languages[language_code],
        log,
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
    parser.add_argument(
        "-l",
        dest="log",
        type=Path,
        help="log non-words",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    if args.language not in languages:
        sys.exit(f"unsupported language: {args.language}")
    if args.output.is_file():
        sys.exit(f"{args.output} is a file")

    process_language(args.language, args.output, args.file, args.log)


if __name__ == "__main__":
    main(parse_args())
