"""Tasks executed by course builder.

These are written to be executed in parallel, and so that targets are built
only when sources are modified.
"""

from dataclasses import dataclass
from functools import cache
from pathlib import Path
from shutil import copyfile, move
from sqlite3 import connect
from tempfile import TemporaryDirectory
import typing as t

from .dependency import is_outdated, Task
from .download import download, has_been_a_week, latest_data
from .link import partition_links
from .migrate import check_scripts, migrate
from .partition import partition
from .populate import populate
from .shrink import shrink
from .tokenizer import process_language
from .untar import untar


build = Path("build")


def download_latest() -> None:
    """Download latest tatoeba data.

    Link and sentence download is one task instead of two, because subsequent
    tasks depend on both of them.
    Hence they are considered one item.
    """
    downloads = build/"tatoeba"
    if has_been_a_week(downloads):
        print("Downloading latest data from Tatoeba")
        download(downloads)


def decompress_links() -> None:
    """Decompress links.tar.bz2."""
    downloads = build/"tatoeba"
    source = latest_data(downloads)[0].destination(downloads)
    target = downloads/"links.csv"

    assert source.is_file()

    if is_outdated([target], [source]):
        print("Decompressing Tatoeba links")
        untar(downloads, source)


def decompress_sentences() -> None:
    """Decompress sentences.tar.bz2."""
    downloads = build/"tatoeba"
    source = latest_data(downloads)[1].destination(downloads)
    target = downloads/"sentences.csv"

    assert source.is_file()

    if is_outdated([target], [source]):
        print("Decompressing Tatoeba sentences")
        untar(downloads, source)


def prepare_sentences() -> None:
    """Prepare sentences for tokenization."""
    source = build/"tatoeba"/"sentences.csv"
    target = build/"sentences"

    assert source.is_file()

    if is_outdated([target], [source]):
        print("Preparing sentences")
        partition(source, target)


def prepare_links() -> None:
    """Partition links by language pair."""
    links = build/"tatoeba"/"links.csv"
    sentences = build/"tatoeba"/"sentences.csv"
    sources = [links, sentences]
    target = build/"links"

    for source in sources:
        assert source.is_file()

    if is_outdated([target], sources):
        print("Preparing links")
        partition_links(links, sentences, target)


@dataclass(unsafe_hash=True)
class LanguageTokenizerTask:
    lang: str

    @property
    def __name__(self) -> str:
        return f"LanguageTokenizerTask({self.lang})"

    def __call__(self) -> None:
        """Run this task."""
        lang = self.lang
        source = build/"sentences"/f"{lang}.tsv"

        outdir = build/"languages"/lang
        targets = [
            outdir/"sentences.csv",
            outdir/"words.csv",
            outdir/"nonwords.txt",
        ]
        assert source.is_file()

        if is_outdated(targets, [source]):
            print(f"Tokenizing words in {lang}")
            process_language(lang, output=outdir, file=source)


@cache
def language_tokenizer(lang: str) -> Task:
    """Create tokenization task for language.

    Cached so that language_tokenizer can be called repeatedly and still refer
    to the same task.
    lang should be a valid language code.

    Returns a class instance, because closures are unpicklable (needed by
    multiprocessing).
    """
    return t.cast(Task, LanguageTokenizerTask(lang))


@dataclass(unsafe_hash=True)
class CourseBuilderTask:
    lang1: str
    lang2: str

    @property
    def __name__(self) -> str:
        return f"CourseBuilderTask({self.lang1}, {self.lang2})"

    def __call__(self) -> None:
        lang1 = self.lang1
        lang2 = self.lang2

        l1_dir = build/"languages"/lang1
        l2_dir = build/"languages"/lang2
        translations = build/"links"/(
            f"{lang1}-{lang2}.csv"
            if lang1 < lang2
            else f"{lang2}-{lang1}.csv"
        )

        sources = [
            build/"test.db",
            l1_dir/"sentences.csv",
            l1_dir/"words.csv",
            l2_dir/"sentences.csv",
            l2_dir/"words.csv",
            translations,
        ]
        target = build/"polycloze"/"courses"/f"{lang1}-{lang2}.db"

        for source in sources:
            assert source.is_file()

        if is_outdated([target], sources):
            print(f"Building {lang1}->{lang2} course")
            with TemporaryDirectory() as tmpname:
                tmp = Path(tmpname)
                database = tmp/"scratch.db"

                copyfile(build/"test.db", database)
                populate(
                    database=database,
                    l1_dir=l1_dir,
                    l2_dir=l2_dir,
                    translations=translations,
                    reversed_=lang1 < lang2,
                )

                with connect(database) as con:
                    shrink(con)
                move(database, target)


@cache
def course_builder(lang1: str, lang2: str) -> Task:
    """Create task for building lang1 -> lang2 course."""
    assert lang1 != lang2
    return t.cast(Task, CourseBuilderTask(lang1, lang2))


def create_empty_course() -> None:
    """Create empty course database file for testing purposes."""
    migrations = Path(__file__).with_name("migrations")

    sources = list(migrations.glob("*.sql"))
    target = build/"test.db"

    for source in sources:
        assert source.is_file()

    if is_outdated([target], sources):
        print("Creating test.db")
        with TemporaryDirectory() as tmpname:
            tmp = Path(tmpname)
            database = tmp/"test.db"

            # Apply migrations to empty database
            with connect(database) as con:
                migrate(con, check_scripts(migrations))

            # shutil.move is used instead of Path.replace, because
            # Path.replace might raise OSError: Invalid cross-device link
            target.parent.mkdir(parents=True, exist_ok=True)
            move(database, target)


def create_course_directory() -> None:
    """Create course directory."""
    (build/"polycloze"/"courses").mkdir(parents=True, exist_ok=True)
