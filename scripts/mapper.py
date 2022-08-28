"""Maps translations."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
import sys


def get_ids(path: Path) -> set[str]:
    with open(path, "r", encoding="utf-8") as file:
        return {line.split()[0] for line in file}


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("l1", type=Path, help="L1 sentences TSV file")
    parser.add_argument("l2", type=Path, help="L2 sentences TSV file")
    parser.add_argument("links", type=Path, help="Tatoeba links file")
    parser.add_argument(
        "-o",
        dest="output",
        type=Path,
        help="output file (default: stdout)",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    if not args.l1.is_file():
        sys.exit(f"{args.l1!s} does not exist")
    if not args.l2.is_file():
        sys.exit(f"{args.l2!s} does not exist")
    if not args.links.is_file():
        sys.exit(f"{args.links!s} does not exist")

    l1_ids = get_ids(args.l1)
    l2_ids = get_ids(args.l2)

    with open(args.links, "r", encoding="utf-8") as file:
        if not args.output:
            for line in file:
                source, target = line.split()
                if source in l1_ids and target in l2_ids:
                    print(f"{source},{target}")
        else:
            with open(args.output, "w", encoding="utf-8") as outfile:
                for line in file:
                    source, target = line.split()
                    if source in l1_ids and target in l2_ids:
                        print(f"{source},{target}", file=outfile)


if __name__ == "__main__":
    main(parse_args())
