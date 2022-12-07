"""Computes word and sentence difficulty."""

from argparse import ArgumentParser, Namespace
import csv
import json
from pathlib import Path
from shutil import move
from sqlite3 import Connection, connect
from tempfile import TemporaryDirectory
import typing as t

from .word import Word


def sources(translations: Path, reverse: bool) -> set[int]:
    """Return ID numbers of translated sentences."""
    result = set()
    with open(translations, encoding="utf-8") as file:
        reader = csv.reader(file)
        for row in reader:
            result.add(int(row[0] if not reverse else row[1]))
    return result


class WordDifficulty(t.NamedTuple):
    """Keeps track of easiest example sentences for a given word.
    Used to compute word difficulty.

    - `frequency_class`: initial guess for word difficulty
    - `examples`: list of difficulties of easiest example sentences
    - `low_count`: number of easiest example sentences to keep track of
        This should be a small number (e.g. 3)
    """
    frequency_class: int
    examples: list[int]
    low_count: int

    @property
    def difficulty(self) -> int:
        """Return word difficulty, while considering the difficulty of the
        sentences it appears in.
        """
        if not self.examples:
            return self.frequency_class
        return self.examples[-1]

    def add_example(self, example_difficulty: int) -> None:
        """Add example sentence for word."""
        self.examples.append(example_difficulty)
        self.examples.sort()
        while len(self.examples) > self.low_count:
            self.examples.pop()


def get_words(language: Path) -> dict[str, WordDifficulty]:
    """Return map of words -> difficulty.

    The result contains initial guess for difficulty: frequency_class.
    """
    with open(language/"words.csv", encoding="utf-8") as file:
        reader = csv.reader(file)
        next(reader)    # Skip header.
        return {row[0]: WordDifficulty(int(row[2]), [], 3) for row in reader}


def is_number(token: str) -> bool:
    """Check if token is a number.

    Also returns true for time, percentages, game scores, etc.
    """
    if not token:
        return False
    for char in token:
        if char in "-.,%:x+ºª€$₱¥£" or char.isdigit():
            continue
        return False
    return True


def compute_difficulty(
    sentence: list[str],
    words: dict[str, WordDifficulty],
) -> int:
    """Compute sentence difficulty.

    Returns -1 if the sentence contains an out-of-vocabulary word.
    Also updates word difficulty for each word in the sentence.
    """
    # Compute sentence difficulty.
    difficulty = -1
    keys = [Word(token) for token in sentence]
    for word in keys:
        # Heuristic rule for excluding non-words, but not punctuation symbols
        # or numbers. Loanwords are excluded.
        if word not in words:
            if len(word) > 1 and not is_number(word):
                return -1
            continue
        value = words[word]

        # NOTE Uses `value.frequency_class` instead of `value.difficulty`,
        # because `value.difficulty` is not stable yet.
        # It becomes stable after all the sentences have been examined.
        difficulty = max(difficulty, value.frequency_class)

    if difficulty < 0:
        return difficulty

    # Record sentence difficulty.
    for word in keys:
        try:
            value = words[word]
            value.add_example(difficulty)
        except KeyError:
            pass
    return difficulty


def load_sentences(     # pylint: disable=too-many-arguments,too-many-locals
    con: Connection,
    language: Path,
    words: dict[str, WordDifficulty],
    translations: Path,
    reversed_: bool,
    outdir: Path,
) -> None:
    """Load sentence into database, but only those with translations.

    Excludes sentences that contains invalid words (not in `words`).
    """
    _sources = sources(translations, reversed_)

    con.execute("""
        CREATE TABLE sentence (
            id INTEGER PRIMARY KEY,
            text TEXT NOT NULL,
            tatoeba_id INTEGER NOT NULL,
            tokens TEXT NOT NULL,   -- JSON array of tokens
            difficulty INTEGER NOT NULL
        )
    """)

    query = """
        INSERT INTO sentence (text, tatoeba_id, tokens, difficulty)
        VALUES (?, ?, ?, ?)
    """
    with TemporaryDirectory() as tmpname:
        temp = Path(tmpname)

        with (
            open(language/"sentences.csv", encoding="utf-8") as infile,
            open(temp/"skipped.csv", "w", encoding="utf-8") as outfile,
        ):
            writer = csv.writer(outfile)
            writer.writerow(["tatoeba_id", "text", "reason_for_exclusion"])

            reader = csv.reader(infile)
            next(reader)    # Skip header.
            for row in reader:
                text = row[1]
                tatoeba_id = int(row[0])
                tokens = row[2]

                # Exclude untranslated sentences.
                if tatoeba_id not in _sources:
                    writer.writerow([tatoeba_id, text, "not translated"])
                    continue

                difficulty = compute_difficulty(json.loads(tokens), words)

                # Exclude sentence contain OOV words (except for punctuation
                # symbols).
                if difficulty < 0:
                    writer.writerow([tatoeba_id, text, "contains OOV word"])
                    continue

                con.execute(query, (text, tatoeba_id, tokens, difficulty))

        outdir.mkdir(parents=True, exist_ok=True)
        move(temp/"skipped.csv", outdir/"skipped.csv")


def write_words(con: Connection, words: dict[str, WordDifficulty]) -> None:
    """Write words to the database."""
    con.execute("""
        CREATE TABLE word (
            word TEXT PRIMARY KEY,
            difficulty INTEGER NOT NULL
        )
    """)

    query = "INSERT INTO word (word, difficulty) VALUES (?, ?)"
    con.executemany(query, (
        (word, value.difficulty)
        for word, value in words.items()
    ))


def compute_difficulty_values(
    language: Path,
    outdir: Path,
    translations: Path,
    reversed_: bool,
) -> None:
    """Compute difficulty values for all words and sentences in a course.

    `outdir`: where course files will be saved
    `reversed_`: whether or not translation table columns are flipped.
    """
    words = get_words(language)
    with TemporaryDirectory() as tmpname:
        tempdir = Path(tmpname)

        with connect(tempdir/"sentences.db") as con:
            load_sentences(
                con,
                language,
                words,
                translations,
                reversed_,
                outdir,
            )

        # NOTE Does not recompute sentence difficulty.
        # But maybe it should?

        with connect(tempdir/"words.db") as con:
            write_words(con, words)

        outdir.mkdir(parents=True, exist_ok=True)
        move(tempdir/"sentences.db", outdir/"sentences.db")
        move(tempdir/"words.db", outdir/"words.db")


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "language",
        type=Path,
        help="path to language directory",
    )
    parser.add_argument(
        dest="translations",
        type=Path,
        help="path to L1->L2 translations file",
    )
    parser.add_argument(
        "-o",
        dest="outdir",
        type=Path,
        required=True,
        help="path to output course directory",
    )
    parser.add_argument(
        "-r",
        dest="reverse",
        action="store_true",
        help="flip columns in translation file",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    compute_difficulty_values(
        args.language,
        args.outdir,
        args.translations,
        args.reverse,
    )


if __name__ == "__main__":
    main(parse_args())
