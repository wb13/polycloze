"""Checks course quality."""

from argparse import ArgumentParser, Namespace
from sqlite3 import connect, Connection
import sys


def quality_check(con: Connection) -> bool:
    """Course should have words in every frequency class and there should be
    enough number of frequency classes.
    """
    query = """
select frequency_class, count(*) from word group by frequency_class
"""
    prev = -1
    for frequency_class, count in con.execute(query):
        if count <= 0 or prev + 1 != frequency_class:
            return False
        prev = frequency_class
    return prev >= 8


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("course", help="course database string")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    with connect(args.course) as con:
        if not quality_check(con):
            sys.exit(1)


if __name__ == "__main__":
    main()
