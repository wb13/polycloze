"""Tasks executed by course builder.

These are written to be executed in parallel, and so that targets are built
only when sources are modified.
"""

from pathlib import Path

from .dependency import is_outdated
from .download import latest_data
from .partition import partition
from .untar import untar


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
