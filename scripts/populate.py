"""Generate course sqlite database."""

from argparse import ArgumentParser, Namespace
import csv
import json
from math import floor, log2
from pathlib import Path
from sqlite3 import Connection, connect
from sys import exit
import typing as t


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "-r",
        dest="reversed",
        help="reverse translation",
        action="store_true",
    )
    parser.add_argument("database", help="sqlite database")
    parser.add_argument("l1", type=Path, help="path to L1 directory")
    parser.add_argument("l2", type=Path, help="path to L2 directory")
    parser.add_argument("translations", type=Path, help="translation CSV file")
    return parser.parse_args()


def sources(translations: Path, reverse: bool = False) -> set[int]:
    result = set()
    with open(translations) as file:
        reader = csv.reader(file)
        for row in reader:
            result.add(int(row[0] if not reverse else row[1]))
    return result


def targets(translations: Path, reverse: bool = False) -> set[int]:
    return sources(translations, not reverse)


def populate_translates(con: Connection, path: Path | str, reverse: bool = False) -> None:
    query = "insert into translates (source, target) values (?, ?)"
    with open(path) as file:
        reader = csv.reader(file)
        for row in reader:
            source = int(row[0])
            target = int(row[1])
            if reverse:
                source, target = target, source

            con.execute(query, (source, target))
        con.commit()


def populate_sentence(con: Connection, language: Path, translations: Path, reverse: bool = False) -> set[str]:
    _sources = sources(translations, reverse)
    query = "insert into sentence (tatoeba_id, text, tokens, frequency_class) values (?, ?, ?, 0)"
    words = set()
    with open(language/"sentences.csv") as file:
        reader = csv.reader(file)
        next(reader)
        for row in reader:
            tatoeba_id = int(row[0])
            text = row[1]
            tokens = row[2]
            if tatoeba_id in _sources:
                con.execute(query, (tatoeba_id, text, tokens))
                words.update(json.loads(tokens))
        con.commit()
    return words


def populate_word(con: Connection, language: Path, words: set[str]) -> None:
    query = "insert into word (word, frequency_class) values (?, ?)"

    with open(language/"words.csv") as file:
        reader = csv.reader(file)
        next(reader)
        row = next(reader)  # first row (highest-frequency word)
        max_frequency = int(row[1])
        con.execute(query, (row[0], 0))

        for row in reader:
            word = row[0]
            frequency = int(row[1])
            frequency_class = int(floor(0.5 - log2(frequency / max_frequency)))
            con.execute(query, (word.casefold(), frequency_class))
        con.commit()


def populate_translation(con: Connection, language: Path, translations: Path, reverse: bool = False) -> None:
    _targets = targets(translations, reverse)
    query = "insert into translation (tatoeba_id, text) values (?, ?)"

    with open(language/"sentences.csv") as file:
        reader = csv.reader(file)
        next(reader)
        for row in reader:
            tatoeba_id = int(row[0])
            text = row[1]
            if tatoeba_id in _targets:
                con.execute(query, (tatoeba_id, text))
        con.commit()


def escape(value: str) -> str:
    """Escape sqlite string."""
    return "'{}'".format(value.replace("'", "''"))


def query_words(con: Connection, words: t.Sequence[str]) -> t.Iterable[int]:
    value = ", ".join(escape(word.casefold()) for word in words)
    query = "select id from word where word in ({})".format(value)
    return (id_ for id_, in con.execute(query))


def populate_contains(con: Connection) -> None:
    query = "select id, tokens from sentence"
    for id_, tokens in con.execute(query):
        query = "insert into contains (sentence, word) values (?, ?)"
        values = ((id_, word_id) for word_id in query_words(con, json.loads(tokens)))
        con.executemany(query, values)
    con.commit()


def create_temp_trigger(con: Connection) -> None:
    query = """
create temp trigger t1 after insert on main.contains begin
    update sentence
    set frequency_class = max(
        frequency_class,
        (select frequency_class from word where word.id = new.word)
    )
    where sentence.id = new.sentence;
end;
"""
    con.execute(query)
    con.commit()


language_names = {
    "cyo": "Cuyonon",
    "deu": "German",
    "eng": "English",
    "spa": "Spanish",
    "tgl": "Tagalog",
}


def infer_language(path: Path) -> tuple[str, str]:
    try:
        code = path.name
        name = language_names[code]
        return (code, name)
    except KeyError:
        exit(f"unknown language code: {path.name}")



def populate_language(con: Connection, l1: Path, l2: Path) -> None:
    query = "insert into language (id, code, name) values (?, ?, ?)"
    con.execute(query, ("l1", *infer_language(l1)))
    con.execute(query, ("l2", *infer_language(l2)))
    con.commit()


def main() -> None:
    args = parse_args()

    with connect(args.database) as con:
        populate_language(con, args.l1, args.l2)

        create_temp_trigger(con)
        populate_translates(con, args.translations, args.reversed)
        words = populate_sentence(con, args.l2, args.translations, args.reversed)
        populate_word(con, args.l2, words)
        populate_translation(con, args.l1, args.translations, args.reversed)
        populate_contains(con)


if __name__ == "__main__":
    main()
