"""Make sqlite dump."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from sqlite3 import Connection, connect


def dump_db(con: Connection) -> None:
    """Make database dump."""
    for line in con.iterdump():
        print(line)


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("database", type=Path, help="database file")
    return parser.parse_args()


def main(args: Namespace) -> None:
    with connect(args.database) as con:
        print("-- Generated using python/scripts/dump.py")
        dump_db(con)


if __name__ == "__main__":
    main(parse_args())
