#!/usr/bin/env python

"""Partitions tatoeba sentences.tsv into multiple files (one per language)."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from shutil import copytree
from sys import exit
from tempfile import TemporaryDirectory


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("out", help="output directory")
    return parser.parse_args()


def partition(basedir: Path):
    files = {}
    try:
        while line := input():
            [id_, language, sentence] = line.split("\t")
            if language not in files:
                files[language] = open(basedir/f"{language}.tsv", "a")
            print(f"{id_}\t{sentence}", file=files[language])
    except EOFError:
        pass
    finally:
        for file in files.values():
            file.close()


def main() -> None:
    args = parse_args()
    out = Path(args.out)
    if out.is_file():
        exit("output file already exists and is not a directory")

    with TemporaryDirectory() as tmpname:
        tmp = Path(tmpname)
        partition(tmp)
        copytree(tmp, out, dirs_exist_ok=True)


if __name__ == "__main__":
    main()
