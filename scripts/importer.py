#!/usr/bin/env python
from argparse import ArgumentParser, Namespace
import csv
import json
from pathlib import Path
from sqlite3 import Connection, connect
import sqlite3
import typing as t

RowImporter = t.Callable[[Connection, csv.reader], None]

def import_csv(
    con: Connection,
    path: Path | str,
    import_row: RowImporter,
) -> None:
    with open(path) as file:
        reader = csv.reader(file)
        next(reader)
        for row in reader:
            import_row(con, row)
        con.commit()


def import_accepted_word_row(ignored_words: set[str]) -> RowImporter:
    def import_word_row(con: Connection, row: tuple[str, ...]) -> None:
        word = row[0]
        if word in ignored_words:
            return
        frequency = int(row[1])
        query = "insert or ignore into word (word, frequency) values (?, ?)"
        con.execute(query, (word, frequency))

    return import_word_row


def import_sentence_row(con: Connection, row: tuple[str, ...]) -> None:
    tatoeba_id = int(row[0])
    text = row[1]
    tokens = row[2]

    query = "insert or ignore into sentence (tatoeba_id, text, tokens) values (?, ?, ?)"
    con.execute(query, (tatoeba_id, text, tokens))


def escape(value: str) -> str:
    """Escape sqlite string."""
    return "'{}'".format(value.replace("'", "''"))


def query_words(con: Connection, words: t.Sequence[str]) -> t.Iterable[int]:
    value = ", ".join(escape(word.casefold()) for word in words)
    query = "select id from word where word in ({})".format(value)
    return (id_ for id_, in con.execute(query))


def insert_contains(con: Connection) -> None:
    query = "select id, tokens from sentence"
    for id_, tokens in con.execute(query):
        query = "insert or ignore into contains (sentence, word) values (?, ?)"
        values = ((id_, word_id) for word_id in query_words(con, json.loads(tokens)))
        con.executemany(query, values)
    con.commit()


def get_max_frequency(con: Connection) -> t.Optional[int]:
    cur = con.cursor()
    cur.execute("select max(frequency) from word")
    row = cur.fetchone()
    if not row:
        return None
    return row[0]


def update_frequency_classes(con: Connection) -> None:
    max_freq = get_max_frequency(con)
    if not max_freq:
        return None

    query = """
update word set frequency_class = cast(floor(0.5 - log2(frequency / ?)) as int)
"""
    con.execute(query, (float(max_freq),))
    con.commit()


def parse_args() -> Namespace:
    """Parse command-line args."""
    parser = ArgumentParser()
    parser.add_argument(
        "database",
        help="sqlite database",
    )
    parser.add_argument(
        "-s",
        dest="sentences_csv",
        help="sentences.csv file",
        required=True,
    )
    parser.add_argument(
        "-w",
        dest="words_csv",
        help="words.csv file",
        required=True,
    )
    parser.add_argument(
        "-i",
        "--ignore",
        dest="ignore",
        help="new-line separated list of words to ignore",
    )
    return parser.parse_args()


def main():
    args = parse_args()
    ignored_words = set()
    if args.ignore:
        try:
            ignored_words = set(Path(args.ignore).read_text().splitlines())
        except FileNotFoundError:
            pass

    with connect(args.database) as con:
        import_csv(con, args.words_csv, import_accepted_word_row(ignored_words))
        import_csv(con, args.sentences_csv, import_sentence_row)
        insert_contains(con)
        update_frequency_classes(con)


if __name__ == "__main__":
    main()
