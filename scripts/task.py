"""Tasks executed by course builder.

These are written to be executed in parallel, and so that targets are built
only when sources are modified.
"""

from pathlib import Path

from .dependency import is_outdated
from .partition import partition


build = Path("build")


def prepare_sentences() -> None:
    """Prepare sentences for tokenization."""
    source = build/"tatoeba"/"sentences.csv"
    target = build/"sentences"

    assert source.is_file()
    print("Preparing sentences")

    if is_outdated([target], [source]):
        partition(source, target)
