"""Shrink course files to target size.

For demo purposes only.
"""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from sqlite3 import Connection, connect
from time import time


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "-t",
        dest="target_count",
        type=int,
        default=10000,
        help="target number of sentences",
    )
    parser.add_argument(
        "courses",
        type=Path,
        nargs="+",
        help="course files to shrink",
    )
    return parser.parse_args()


def cap_sentences(con: Connection, target_count: int) -> None:
    """Cap number of sentences to target count."""
    query = """
delete from sentence where id not in (
    select id from sentence order by random() limit ?)
"""
    con.execute(query, (target_count,))


def prune(con: Connection) -> None:
    """Delete unused data."""

    # 0. delete sentences
    # 1. delete words orphaned by deleted sentences
    con.execute("""
delete from contains where sentence not in (select id from sentence)
""")
    con.execute("delete from word where id not in (select word from contains)")

    # 2. delete translations orphaned by deleted sentences
    con.execute("""
delete from translates where source not in (select tatoeba_id from sentence)
""")
    con.execute("""
delete from translation where tatoeba_id not in (select target from translates)
""")


def main() -> None:
    args = parse_args()
    for path in args.courses:
        print(f"shrinking {path!s}")
        start = time()
        with connect(path) as con:
            cap_sentences(con, args.target_count)
            prune(con)
            con.commit()
            con.executescript("vacuum")

            print("took", time() - start)


if __name__ == "__main__":
    main()
