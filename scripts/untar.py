"""Unarchive tatoeba data."""

from argparse import ArgumentParser, Namespace
from concurrent.futures import ProcessPoolExecutor
from os import utime
from pathlib import Path
from shutil import copytree
import sys
import tarfile
from tempfile import TemporaryDirectory

from .dependency import is_outdated
from .download import latest_data


def untar(destination: Path, infile: Path) -> None:
    """Unarchive single file into destination."""
    with tarfile.open(infile, "r:bz2") as tar:
        tar.extractall(destination)


def parse_args() -> Namespace:
    parser = ArgumentParser(
        description="Unarchive Tatoeba data.",
    )
    parser.add_argument(
        "-l",
        dest="links",
        type=Path,
        help="Tatoeba links.tar.bz2 file",
    )
    parser.add_argument(
        "-s",
        dest="sentences",
        type=Path,
        help="Tatoeba sentences.tar.bz2 file",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    downloads = Path("build")/"tatoeba"
    if not args.links or not args.sentences:
        try:
            links, sentences = latest_data(downloads)
            args.links = links.destination(downloads)
            args.sentences = sentences.destination(downloads)

            assert args.links.is_file()
            assert args.sentences.is_file()
        except AssertionError:
            sys.exit("no data found")

    if not is_outdated(
        [downloads/"links.csv", downloads/"sentences.csv"],
        [args.links, args.sentences],
    ):
        return

    print("Extracting data...")
    with (
        ProcessPoolExecutor() as executor,
        TemporaryDirectory() as tmpname,
    ):
        tmp = Path(tmpname)
        futures = [
            executor.submit(untar, tmp, args.links),
            executor.submit(untar, tmp, args.sentences),
        ]
        for future in futures:
            future.result()
        copytree(tmp, downloads, dirs_exist_ok=True)

    # Update mtime, because original mtime < mtime of tar archive.
    utime(downloads/"links.csv")
    utime(downloads/"sentences.csv")
    print("Done extracting data")


if __name__ == "__main__":
    main(parse_args())
