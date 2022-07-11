#!/usr/bin/env python

"""Sets language metadata in database (info table)."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from sqlite3 import Connection, connect
import sys

from .languages import languages


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "db",
        help="language sqlite3 database",
        type=Path,
    )
    return parser.parse_args()


metadata = {
    "cyo": ("Cuyonon", "Cuyonon"),
    "deu": ("German", "Deutsch"),
    "eng": ("English", "English"),
    "spa": ("Spanish", "Español"),
    "tgl": ("Tagalog", "Tagalog"),
}


def check(language: str) -> None:
    if language not in languages or language not in metadata:
        sys.exit(f"unsupported language: {language}")


def set_metadata(con: Connection, code: str) -> None:
    query = """
insert into info (code, english, native) values (?, ?, ?)
"""
    english, native = metadata[code]
    con.execute(query, (code, english, native))
    con.commit()


def main() -> None:
    args = parse_args()

    language = args.db.stem
    check(language)
    with connect(args.db) as con:
        set_metadata(con, language)


if __name__ == "__main__":
    main()
