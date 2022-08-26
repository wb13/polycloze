#!/usr/bin/env python

"""Partitions tatoeba sentences.tsv into multiple files (one per language)."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from shutil import copytree
import sys
from tempfile import TemporaryDirectory


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


def output_parsed_line(line: str, basedir: Path) -> None:
    [id_, language, sentence] = line.strip().split("\t")
    with open(basedir/f"{language}.tsv", "a", encoding="utf-8") as outfile:
        print(f"{id_}\t{sentence}", file=outfile)


def partition(inputfile: Path | None, basedir: Path) -> None:
    if inputfile is None:
        try:
            while line := input():
                output_parsed_line(line, basedir)
        except EOFError:
            pass
    else:
        with open(inputfile, encoding="utf-8") as infile:
            for line in infile:
                output_parsed_line(line, basedir)


def main() -> None:
    args = parse_args()
    out = Path(args.out)
    if out.is_file():
        sys.exit("output file already exists and is not a directory")

    with TemporaryDirectory() as tmpname:
        tmp = Path(tmpname)
        partition(args.file, tmp)
        copytree(tmp, out, dirs_exist_ok=True)


if __name__ == "__main__":
    main()
