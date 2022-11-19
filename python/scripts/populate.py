"""Generate course sqlite database."""

from argparse import ArgumentParser, Namespace
import csv
import json
from pathlib import Path
from sqlite3 import Connection, connect
import sys
import typing as t

from .language import languages


def sources(translations: Path, reverse: bool = False) -> set[int]:
    result = set()
    with open(translations, encoding="utf-8") as file:
        reader = csv.reader(file)
        for row in reader:
            result.add(int(row[0] if not reverse else row[1]))
    return result


def targets(translations: Path, reverse: bool = False) -> set[int]:
    return sources(translations, not reverse)


def get_words(language: Path) -> dict[str, int]:
    """Get words in language, mapped to frequency_class."""
    path = language/"words.csv"
    with open(path, encoding="utf-8") as file:
        reader = csv.reader(file)
        next(reader)    # Skip header
        return {row[0]: int(row[2]) for row in reader}


def compute_frequency_class(
    sentence: t.Iterable[str],
    words: dict[str, int],
) -> int:
    """Sentence frequency class = max(frequency class of word in sentence).

    Assumes the tokens in the sentence have been casefolded.
    """
    return max(words.get(token, 0) for token in sentence)


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


def populate_sentence(  # pylint: disable=too-many-locals
    con: Connection,
    words: dict[str, int],
    language: Path,
    translations: Path,
    reverse: bool = False,
) -> None:
    _sources = sources(translations, reverse)
    query = """
insert into sentence (tatoeba_id, text, tokens, frequency_class)
values (?, ?, ?, ?)
"""
    with (
        open(language/"sentences.csv", encoding="utf-8") as file,
        open(language/"skipped-sentences.txt", "a", encoding="utf-8") as log,
    ):
        reader = csv.reader(file)
        next(reader)
        for row in reader:
            tatoeba_id = int(row[0])
            text = row[1]
            tokens = row[2]
            if tatoeba_id in _sources:
                sentence = [token.casefold() for token in json.loads(tokens)]
                for token in sentence:
                    if len(token) > 1 and token not in words:
                        # NOTE This is a heuristic for excluding non-words,
                        # but not punctuation symbols.
                        break
                else:
                    # Insert sentence only if all tokens are words or
                    # punctuation.
                    frequency_class = compute_frequency_class(sentence, words)
                    con.execute(
                        query,
                        (tatoeba_id, text, tokens, frequency_class),
                    )
                    continue

                # Log skipped sentence.
                print(text, file=log)
        con.commit()


def populate_word(con: Connection, words: dict[str, int]) -> None:
    """Insert words into database.

    May include words that don't belong to any DB.
    """
    query = "INSERT INTO word (word, frequency_class) VALUES (?, ?)"
    con.executemany(query, words.items())
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
            text = row[1]
            if tatoeba_id in _targets:
                con.execute(query, (tatoeba_id, text))
        con.commit()


def escape(value: str) -> str:
    """Escape sqlite string."""
    replaced = value.replace("'", "''")
    return f"'{replaced}'"


def query_words(con: Connection, words: t.Sequence[str]) -> t.Iterable[int]:
    value = ", ".join(escape(word.casefold()) for word in words)
    query = f"select id from word where word in ({value})"
    return (id_ for id_, in con.execute(query))


# Word ID -> number of example sentences
count_examples: dict[int, int] = {}


def populate_contains(con: Connection, max_number_examples: int) -> None:
    """Link words to sentence they belong to.

    Caps number of linked sentence per word to `max_number_examples`.
    """
    query = "SELECT id, tokens FROM sentence ORDER BY frequency_class ASC"
    for id_, tokens in con.execute(query):
        query = "insert into contains (sentence, word) values (?, ?)"
        word_ids = list(query_words(con, json.loads(tokens)))
        values = (
            (id_, word_id)
            for word_id in word_ids
            if count_examples.get(word_id, 0) < max_number_examples
        )
        con.executemany(query, values)

        # Count should only be increased once, even if the word appears
        # multiple times in the sentence.
        for word_id in word_ids:
            count_examples[word_id] = count_examples.get(word_id, 0) + 1
    con.commit()


def infer_language(path: Path) -> tuple[str, str, str]:
    try:
        code = path.name
        language = languages[code]
        return (code, language.name, language.bcp47)
    except KeyError:
        sys.exit(f"unknown language code: {path.name}")


def populate_language(con: Connection, lang1: Path, lang2: Path) -> None:
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
    l1_dir: Path,
    l2_dir: Path,
    translations: Path,
    reversed_: bool,
) -> None:
    """Populate course database.

    reversed: whether or not translation table columns are swapped
    """
    with connect(database) as con:
        populate_language(con, l1_dir, l2_dir)
        populate_translates(con, translations, reversed_)
        words = get_words(l2_dir)
        populate_sentence(
            con,
            words,
            l2_dir,
            translations,
            reversed_,
        )
        populate_word(con, words)
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
