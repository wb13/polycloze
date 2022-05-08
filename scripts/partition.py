#!/usr/bin/env python

"""Partitions tatoeba sentences.tsv into multiple files (one per language)."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
from sys import exit


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("out", help="output directory")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    out = Path(args.out)
    if not out.is_dir():
        exit(f"{args.out} is not a directory")

    files = {}
    try:
        while line := input():
            [id_, language, sentence] = line.split("\t")
            if language not in files:
                files[language] = open(out/f"{language}.tsv", "a")
            print(f"{id_}\t{sentence}", file=files[language])
    except EOFError:
        pass
    finally:
        for file in files.values():
            file.close()


if __name__ == "__main__":
    main()
