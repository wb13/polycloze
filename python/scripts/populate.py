"""Generate course sqlite database."""

from argparse import ArgumentParser, Namespace
from collections import Counter
import csv
import json
from pathlib import Path
from sqlite3 import Connection, connect
import sys
import typing as t

from .language import languages


def targets(translations: Path, reverse: bool = False) -> set[int]:
    result = set()
    with open(translations, encoding="utf-8") as file:
        reader = csv.reader(file)
        for row in reader:
            result.add(int(row[0] if reverse else row[1]))
    return result


def populate_translates(
    con: Connection,
    path: Path | str,
    reverse: bool = False,
) -> None:
    query = "insert into translates (source, target) values (?, ?)"
    with open(path, encoding="utf-8") as file:
        reader = csv.reader(file)
        for row in reader:
            source = int(row[0])
            target = int(row[1])
            if reverse:
                source, target = target, source

            con.execute(query, (source, target))
        con.commit()


def populate_sentence(con: Connection, course: Path) -> None:
    con.execute(
        "ATTACH DATABASE ? AS ts",
        (str((course/"sentences.db").resolve()),),
    )
    con.execute("""
        INSERT INTO sentence (tatoeba_id, text, tokens, frequency_class)
        SELECT tatoeba_id, text, tokens, difficulty
        FROM ts.sentence
    """)
    con.commit()


def populate_word(con: Connection, course: Path) -> None:
    """Insert words into database.

    May include words that don't belong to any sentence.
    """
    con.execute(
        "ATTACH DATABASE ? AS tw",
        (str((course/"words.db").resolve()),),
    )
    query = """
        INSERT INTO word (word, frequency_class)
        SELECT word, difficulty
        FROM tw.word
    """
    con.execute(query)
    con.commit()


def populate_translation(
    con: Connection,
    language: Path,
    translations: Path,
    reverse: bool = False,
) -> None:
    _targets = targets(translations, reverse)
    query = "insert into translation (tatoeba_id, text) values (?, ?)"

    with open(language/"sentences.csv", encoding="utf-8") as file:
        reader = csv.reader(file)
        next(reader)
        for row in reader:
            tatoeba_id = int(row[0])
            if tatoeba_id in _targets:
                con.execute(query, (tatoeba_id, row[1]))
        con.commit()


def escape(value: str) -> str:
    """Escape sqlite string."""
    replaced = value.replace("'", "''")
    return f"'{replaced}'"


def query_words(con: Connection, words: t.Sequence[str]) -> t.Iterable[int]:
    value = ", ".join(escape(word.casefold()) for word in words)
    query = f"select id from word where word in ({value})"
    return (id_ for id_, in con.execute(query))


def populate_contains(con: Connection, max_number_examples: int) -> None:
    """Link words to sentence they belong to.

    Caps number of linked sentence per word to `max_number_examples`.
    """
    counter: Counter[int] = Counter()   # counts example sentences

    query = "SELECT id, tokens FROM sentence ORDER BY frequency_class ASC"
    for id_, tokens in con.execute(query):
        query = "insert into contains (sentence, word) values (?, ?)"
        word_ids = list(query_words(con, json.loads(tokens)))
        values = (
            (id_, word_id)
            for word_id in word_ids
            if counter[word_id] < max_number_examples
        )
        con.executemany(query, values)

        # Count should only be increased once, even if the word appears
        # multiple times in the sentence.
        for word_id in word_ids:
            counter.update([word_id])
    con.commit()


def infer_language(code: str) -> tuple[str, str, str]:
    try:
        language = languages[code]
        return (code, language.name, language.bcp47)
    except KeyError:
        sys.exit(f"unknown language code: {code}")


def populate_language(con: Connection, course: Path) -> None:
    lang1, lang2 = course.name.split("-")
    query = "insert into language (id, code, name, bcp47) values (?, ?, ?, ?)"
    con.execute(query, ("l1", *infer_language(lang1)))
    con.execute(query, ("l2", *infer_language(lang2)))
    con.commit()


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


def populate(
    database: Path,
    course: Path,
    l1_dir: Path,
    translations: Path,
    reversed_: bool,
) -> None:
    """Populate course database.

    reversed: whether or not translation table columns are swapped
    """
    with connect(database) as con:
        populate_language(con, course)
        populate_translates(con, translations, reversed_)
        populate_sentence(con, course)
        populate_word(con, course)
        populate_translation(con, l1_dir, translations, reversed_)
        populate_contains(con, max_number_examples=30)


def main(args: Namespace) -> None:
    populate(
        args.database,
        args.l1,
        args.l2,
        args.translations,
        args.reversed,
    )


if __name__ == "__main__":
    main(parse_args())
