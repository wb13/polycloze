"""Computes word and sentence difficulty."""

from argparse import ArgumentParser, Namespace
import csv
import json
from pathlib import Path
from shutil import move
from sqlite3 import Connection, connect
from tempfile import TemporaryDirectory
import typing as t


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


def compute_difficulty(
    sentence: list[str],
    words: dict[str, WordDifficulty],
) -> int:
    """Compute sentence difficulty.

    Returns -1 if the sentence contains an out-of-vocabulary word.
    Also updates word difficulty for each word in the sentence.
    """
    # Compute sentence difficulty.
    difficulty = 0
    keys = [token.casefold() for token in sentence]
    for word in keys:
        # Heuristic rule for excluding non-words, but not punctuation symbols.
        if word not in words:
            if len(word) > 1:
                return -1
            continue
        value = words[word]

        # NOTE Uses `value.frequency_class` instead of `value.difficulty`,
        # because `value.difficulty` is not stable yet.
        # It becomes stable after all the sentences have been examined.
        difficulty = max(difficulty, value.frequency_class)

    # Record sentence difficulty.
    for word in keys:
        try:
            value = words[word]
            value.add_example(difficulty)
        except KeyError:
            pass
    return difficulty


def load_sentences(
    con: Connection,
    language: Path,
    words: dict[str, WordDifficulty],
) -> None:
    """Load sentence into database.

    Excludes sentences that contains invalid words (not in `words`).
    """
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
            open(temp/"sentences-oov.txt", "a", encoding="utf-8") as outfile,
        ):
            reader = csv.reader(infile)
            next(reader)    # Skip header.

            for row in reader:
                text = row[1]
                tatoeba_id = int(row[0])
                tokens = row[2]
                difficulty = compute_difficulty(json.loads(tokens), words)

                # Only include sentence that don't contain OOV words.
                if difficulty >= 0:
                    con.execute(query, (text, tatoeba_id, tokens, difficulty))
                else:
                    # Log sentence with OOV word
                    print(text, file=outfile)

        move(temp/"sentences-oov.txt", language/"sentences-oov.txt")


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


def compute_difficulty_values(language: Path) -> None:
    """Compute difficulty values for all words and sentences in a language."""
    words = get_words(language)
    with TemporaryDirectory() as tmpname:
        tempdir = Path(tmpname)

        with connect(tempdir/"sentences.db") as con:
            load_sentences(con, language, words)

        # NOTE Does not recompute sentence difficulty.
        # But maybe it should?

        with connect(tempdir/"words.db") as con:
            write_words(con, words)

        move(tempdir/"sentences.db", language/"sentences.db")
        move(tempdir/"words.db", language/"words.db")


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "language",
        type=Path,
        help="path to language directory",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    compute_difficulty_values(args.language)


if __name__ == "__main__":
    main(parse_args())
