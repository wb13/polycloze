"""Shrink course files by reducing number of example sentences per word."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from sqlite3 import Connection, connect
from time import time


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "courses",
        type=Path,
        nargs="+",
        help="course files to shrink",
    )
    return parser.parse_args()


def cap_sentences(con: Connection) -> None:
    """Exclude sentence examples that are too difficult."""
    # Drop index.
    con.execute("DROP INDEX index_contains_word")

    # Create new contains table.
    con.execute("""
        CREATE TABLE new_contains (
            sentence INTEGER NOT NULL REFERENCES sentence,
            word INTEGER NOT NULL REFERENCES word
        )
    """)

    # Populate new table.
    con.execute("""
        INSERT INTO new_contains (sentence, word)
        SELECT sentence.id, word.id
        FROM sentence
        JOIN contains ON (sentence.id = contains.sentence)
        JOIN word ON (word.id = contains.word)
        WHERE sentence.frequency_class <= word.frequency_class
    """)

    # Replace old table with new table.
    con.execute("DROP TABLE contains")
    con.execute("ALTER TABLE new_contains RENAME TO contains")

    # Recreate index.
    con.execute("""
        CREATE INDEX index_contains_word ON contains (word)
    """)


def shrink(con: Connection) -> None:
    """Shrink course file by reducing number of example sentences per word.

    The caller doesn't have to call `.commit()` afterwards.
    """
    cap_sentences(con)
    delete_orphans(con)

    con.commit()
    con.executescript("vacuum")


def delete_orphaned_sentences(con: Connection) -> None:
    """Delete orphaned sentences."""
    query = """
        DELETE FROM sentence
        WHERE id NOT IN (
            SELECT sentence FROM contains
        )
    """
    con.execute(query)


def delete_orphaned_translations(con: Connection) -> None:
    """Delete orphaned translations."""
    con.execute("""
        DELETE FROM translates
        WHERE source NOT IN (
            SELECT tatoeba_id FROM sentence
        )
    """)
    con.execute("""
        DELETE FROM translation
        WHERE tatoeba_id NOT IN (
            SELECT target FROM translates
        )
    """)


def delete_orphans(con: Connection) -> None:
    """Delete orphaned data."""
    delete_orphaned_sentences(con)
    delete_orphaned_translations(con)

    # Some sentences might have a tatoeba translation, but the translation is
    # too long.
    con.execute("""
        DELETE FROM translates
        WHERE target NOT IN (
            SELECT tatoeba_id from translation
        )
    """)
    con.execute("""
        DELETE FROM sentence
        WHERE tatoeba_id NOT IN (
            SELECT source FROM translates
        )
    """)
    con.execute("""
        DELETE FROM contains
        WHERE sentence NOT IN (
            SELECT id FROM sentence
        )
    """)
    # Delete words that appear in untranslated sentences (e.g. including
    # sentences that do have a translation, but the translation is too long).
    query = """
        DELETE FROM word
        WHERE id NOT IN (
            SELECT word FROM contains
        )
    """
    con.execute(query)


def main() -> None:
    args = parse_args()
    for path in args.courses:
        print(f"shrinking {path!s}")
        start = time()
        with connect(path) as con:
            shrink(con)
            print("took", time() - start)


if __name__ == "__main__":
    main()
