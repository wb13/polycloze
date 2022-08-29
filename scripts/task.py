"""Tasks executed by course builder.

These are written to be executed in parallel, and so that targets are built
only when sources are modified.
"""

from functools import cache
from pathlib import Path
import typing as t

from .dependency import is_outdated
from .download import latest_data
from .mapper import map_translations
from .partition import partition
from .tokenizer import process_language
from .untar import untar


Task = t.Callable[[], t.Any]

build = Path("build")
downloads = build/"tatoeba"


def decompress_links() -> None:
    """Decompress links.tar.bz2."""
    source = latest_data(downloads)[0].destination(downloads)
    target = downloads/"links.csv"

    assert source.is_file()
    print("Decompressing Tatoeba links")

    if is_outdated([target], [source]):
        untar(downloads, source)


def decompress_sentences() -> None:
    """Decompress sentences.tar.bz2."""
    source = latest_data(downloads)[1].destination(downloads)
    target = downloads/"sentences.csv"

    assert source.is_file()
    print("Decompressing Tatoeba sentences")

    if is_outdated([target], [source]):
        untar(downloads, source)


def prepare_sentences() -> None:
    """Prepare sentences for tokenization."""
    source = build/"tatoeba"/"sentences.csv"
    target = build/"sentences"

    assert source.is_file()
    print("Preparing sentences")

    if is_outdated([target], [source]):
        partition(source, target)


@cache
def language_tokenizer(lang: str) -> Task:
    """Create tokenization task for language.

    Cached so that language_tokenizer can be called repeatedly and still refer
    to the same task.
    lang should be a valid language code.
    """
    def tokenize_language() -> None:
        source = build/"sentences"/f"{lang}.tsv"

        log = build/"logs"/"nonwords"/f"{lang}.txt"
        sentences = build/"languages"/lang/"sentences.csv"
        words = build/"languages"/lang/"words.csv"
        targets = [log, sentences, words]

        assert source.is_file()
        print(f"Tokenizing words in {lang}")

        if is_outdated(targets, [source]):
            process_language(
                lang,
                output=build/"languages"/lang,
                file=source,
                log=log,
            )
    return tokenize_language


@cache
def translation_mapper(lang1: str, lang2: str) -> Task:
    """Creating task for mapping translations between lang1 and lang2.

    @cache'd for the same reason as language_tokenizer.
    Asserts lang1 < lang2, because lang1->lang2 and lang2->lang1 use the same
    translation file.
    """
    assert lang1 < lang2

    def task() -> None:
        l1_sentences = build/"sentences"/f"{lang1}.tsv"
        l2_sentences = build/"sentences"/f"{lang2}.tsv"
        links = build/"tatoeba"/"links.csv"

        sources = [l1_sentences, l2_sentences, links]
        target = build/"translations"/f"{lang1}-{lang2}.csv"

        for source in sources:
            assert source.is_file()
        print(f"Mapping translations between {lang1} and {lang2}")

        if is_outdated([target], sources):
            map_translations(l1_sentences, l2_sentences, links, output=target)
    return task
