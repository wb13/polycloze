from argparse import ArgumentParser, Namespace
import csv
from pathlib import Path

"""Converts CSV file into newline-separated list of strings (entries of first
column only)."""


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument(
        "csv",
        help="CSV file",
        type=Path,
    )
    parser.add_argument(
        "--no-header",
        dest="has_header",
        help="CSV file has no header",
        action="store_false",
    )
    return parser.parse_args()


def main():
    args = parse_args()
    with open(args.csv) as file:
        reader = csv.reader(file)
        if args.has_header:
                next(reader, None)
        for row in reader:
            print(row[0])


if __name__ == "__main__":
    main()
