#!/usr/bin/env python

"""Partitions tatoeba sentences.tsv into multiple files (one per language)."""

from argparse import ArgumentParser, Namespace
import fileinput
from pathlib import Path
from shutil import copytree
import sys
from tempfile import TemporaryDirectory

from .dependency import is_outdated


def output_parsed_line(line: str, basedir: Path) -> None:
    [id_, language, sentence] = line.strip().split("\t")
    with open(basedir/f"{language}.tsv", "a", encoding="utf-8") as outfile:
        print(f"{id_}\t{sentence}", file=outfile)


def partition(inputfile: Path | None, out: Path) -> None:
    """Extract sentences from inputfile.

    Sentences are partitioned by language.
    These are first written in a temp directory before being copied to the
    output directory, so that output files are never half-finished.
    """
    with (
        fileinput.input(files=inputfile or "-", encoding="utf-8") as file,
        TemporaryDirectory() as tmpname,
    ):
        tmp = Path(tmpname)
        for line in file:
            output_parsed_line(line, tmp)
        copytree(tmp, out, dirs_exist_ok=True)


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("out", help="output directory")
    parser.add_argument(
            "-f",
            dest="file",
            type=Path,
            help="sentences.csv file (default: stdin)",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    out = Path(args.out)
    if out.is_file():
        sys.exit("output file already exists and is not a directory")
    if is_outdated([args.out], [args.file]):
        partition(args.file, out)


if __name__ == "__main__":
    main(parse_args())
