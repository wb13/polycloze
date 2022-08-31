"""Maps translations."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
import sys


def get_ids(path: Path) -> set[str]:
    with open(path, "r", encoding="utf-8") as file:
        return {line.split()[0] for line in file}


def map_translations(
    lang1: Path,
    lang2: Path,
    links: Path,
    output: Path | None = None,
) -> None:
    """Map translations between L1 and L2 languages.

    lang1: path to L1 sentences
    lang2: path to L2 sentences
    links: Tatoeba links file
    output: output file or None (stdout)
    """
    assert lang1.is_file()
    assert lang2.is_file()
    assert links.is_file()

    l1_ids = get_ids(lang1)
    l2_ids = get_ids(lang2)

    with open(links, "r", encoding="utf-8") as file:
        if not output:
            for line in file:
                source, target = line.split()
                if source in l1_ids and target in l2_ids:
                    print(f"{source},{target}")
        else:
            output.parent.mkdir(parents=True, exist_ok=True)
            with open(output, "w", encoding="utf-8") as outfile:
                for line in file:
                    source, target = line.split()
                    if source in l1_ids and target in l2_ids:
                        print(f"{source},{target}", file=outfile)


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
    map_translations(args.l1, args.l2, args.links, args.output)


if __name__ == "__main__":
    main(parse_args())
