"""Partitions tatoeba sentences.tsv into separate files per language."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from shlex import quote
from shutil import copytree
from subprocess import CalledProcessError, run
import sys
from tempfile import TemporaryDirectory


class SortError(Exception):
    """Raised when `sort` program fails."""


def sort_file_lines(infile: Path, outfile: Path) -> None:
    """Write sorted lines in `infile` to `outfile`.

    May raise `SortError`.
    """
    with TemporaryDirectory() as tmpname:
        sorted_file = Path(tmpname)/"sorted.csv"
        args = [
            "sort",
            "-k2,2",
            quote(str(infile.resolve())),
            "-o",
            quote(str(sorted_file.resolve())),
        ]
        try:
            run(args, check=True)
        except CalledProcessError as exc:
            raise SortError from exc

        sorted_file.replace(outfile)


def partition(infile: Path, outdir: Path) -> None:
    """Extract sentences from `infile`.

    Sentences are partitioned by language.
    These are first written in a temp directory before being copied to the
    output directory, so that output files are never half-finished.

    May raise `SortError`.
    """
    with TemporaryDirectory() as tmpname:
        # Sort by language, so all sentences in the same language are together.
        sorted_file = Path(tmpname)/"sorted.csv"
        staging = Path(tmpname)/"staging"
        staging.mkdir()

        sort_file_lines(infile, sorted_file)

        with open(sorted_file, encoding="utf-8") as lines:
            prev = None
            file = None
            for line in lines:
                [id_, language, sentence] = line.strip().split("\t")

                if language != prev:
                    print(language)
                    if file is not None:
                        file.close()
                    prev = language
                    file = open(    # pylint: disable=consider-using-with
                        staging/f"{language}.tsv",
                        "a",
                        encoding="utf-8",
                    )

                print(f"{id_}\t{sentence}", file=file)

            if file is not None:
                file.close()

        copytree(staging, outdir, dirs_exist_ok=True)


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("out", help="output directory")
    parser.add_argument(
        "-f",
        dest="file",
        type=Path,
        help="sentences.csv file",
        required=True,
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    out = Path(args.out)
    if out.is_file():
        sys.exit("destination already exists as a file")
    partition(args.file, out)


if __name__ == "__main__":
    main(parse_args())
