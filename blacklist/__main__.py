from argparse import ArgumentParser, Namespace
from sys import exit

from . import alphabet


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("language", help="ISO 639-3 language code")
    parser.add_argument(
        "-w",
        "--whitelist",
        help="show whitelisted words",
        action="store_true",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    try:
        alphabet_ = alphabet.load(args.language)
    except alphabet.UnsupportedLanguage:
        exit(f"unsupported language: {args.language}")

    try:
        if args.whitelist:
            while line := input():
                if alphabet.is_word(alphabet_, line):
                    print(line)
        else:
            while line := input():
                if not alphabet.is_word(alphabet_, line):
                    print(line)
    except EOFError:
        pass


if __name__ == "__main__":
    main()
