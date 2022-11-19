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


def bump_up_frequency_class(con: Connection) -> None:
    """Bump up frequency class of words that never appear as the most difficult
    word in a sentence.
    The new frequency class of such words will be the lowest frequency class of
    any sentence that the word appears in.
    """
    query = """
        UPDATE word
        SET frequency_class = affected.frequency_class
        FROM (
            SELECT a AS id, c AS frequency_class
            FROM (
                SELECT
                    word.id AS a,
                    word.frequency_class AS b,
                    min(sentence.frequency_class) AS c
                FROM word
                JOIN contains ON (word.id = contains.word)
                JOIN sentence ON (sentence.id = contains.sentence)
                GROUP BY (word.id)
            )
            WHERE b < c
        ) AS affected
        WHERE word.id = affected.id;
    """
    con.execute(query)


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
    bump_up_frequency_class(con)
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

    # There may be orphaned words, not because sentences were removed,
    # but because they originally didn't belong to any sentence.
    # This happens with words that appear in untranslated sentences.
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
