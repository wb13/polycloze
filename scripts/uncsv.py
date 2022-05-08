#!/usr/bin/env python

from argparse import ArgumentParser, Namespace
import csv
from pathlib import Path

"""Converts single-row CSV file into newline-separated list of strings."""


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "csv",
        help="single-row CSV file",
        type=Path,
    )
    return parser.parse_args()


def main():
    args = parse_args()
    with open(args.csv) as file:
        reader = csv.reader(file)
        for row in reader:
            print(row[0])


if __name__ == "__main__":
    main()
