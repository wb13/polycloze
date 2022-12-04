# pylint: disable=too-many-locals
"""Partition Tatoeba links by language pair."""

from argparse import ArgumentParser, Namespace
from collections import deque
from contextlib import contextmanager
from os import utime
from pathlib import Path
from shutil import copytree
import sys
from tempfile import TemporaryDirectory
import typing as t


def sentence_languages(sentences: Path) -> dict[int, str]:
    """Return dict: sentence ID -> language code.

    The dict should probably be < 500MB for 10_000_000 sentences.
    """
    with open(sentences, encoding="utf-8") as file:
        return {
            int((row := line.split("\t", maxsplit=2))[0]): row[1]
            for line in file
        }


class LinkFiles:
    """Keeps track of open files for language pair links.

    Limits number of open files to `limit`.
    """
    def __init__(self, basedir: Path, limit: int = 100) -> None:
        assert limit > 0

        self.basedir = basedir
        self.limit = limit

        Key: t.TypeAlias = tuple[str, str]
        self.files: dict[Key, t.TextIO] = {}
        self.recent: deque[Key] = deque()

    def get(self, source_language: str, target_language: str) -> t.TextIO:
        """Return text IO stream to write link to."""
        assert source_language < target_language

        key = (source_language, target_language)

        try:
            return self.files[key]
        except KeyError:
            # Evict oldest if there are too many open files.
            if len(self.recent) > self.limit:
                oldest = self.recent.popleft()
                file = self.files.pop(oldest)
                file.close()

            self.recent.append(key)
            return self.files.setdefault(
                key,
                open(
                    self.basedir/f"{source_language}-{target_language}.csv",
                    "a",    # Append to keep previously written links.
                    encoding="utf-8",
                ),
            )

    def close(self) -> None:
        """Close all open files."""
        self.recent.clear()
        for _, file in self.files.items():
            file.close()


@contextmanager
def link_files(basedir: Path) -> t.Iterator[LinkFiles]:
    """Context manager for LinkFiles."""
    files = LinkFiles(basedir)
    try:
        yield files
    finally:
        files.close()


def partition_links(links: Path, sentences: Path, outdir: Path) -> None:
    """Partition Tatoeba links by language pair.

    Writes one CSV file of links for each language pair.
    Since translations are symmetric, only one of L1-L2 and L2-L1 will be
    written (L1 < L2).

    - `links`: path to Tatoeba links TSV file
    - `sentences`: path to Tatoeba sentences TSV file
    """
    language = sentence_languages(sentences)

    with (
        TemporaryDirectory() as tmpname,
        open(links, encoding="utf-8") as infile,
        link_files(tempdir := Path(tmpname)) as files,
    ):
        for line in infile:
            row = line.split("\t")
            source = int(row[0])
            target = int(row[1])

            try:
                source_language = language[source]
                target_language = language[target]
            except KeyError:
                # Some tatoeba links refer to deleted sentences.
                continue

            # Make sure `source_language < target_language`.
            if source_language >= target_language:
                continue
            outfile = files.get(source_language, target_language)
            print(f"{source},{target}", file=outfile)

        copytree(tempdir, outdir, dirs_exist_ok=True)

        for path in outdir.iterdir():
            utime(path)
        utime(outdir)


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "-l",
        dest="links",
        type=Path,
        help="Tatoeba links.csv file (actually a TSV file)",
        required=True,
    )
    parser.add_argument(
        "-s",
        dest="sentences",
        type=Path,
        help="Tatoeba sentences.csv file (actually a TSV file)",
        required=True,
    )
    parser.add_argument(
        "-o",
        dest="outdir",
        type=Path,
        help="output directory for partitioned links",
        required=True,
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    if args.outdir.is_file():
        sys.exit(f"{args.outdir!s} is a file")

    partition_links(args.links, args.sentences, args.outdir)


if __name__ == "__main__":
    main(parse_args())
