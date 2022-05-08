from argparse import ArgumentParser, Namespace
from sys import exit

from . import alphabet


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("language", help="ISO 639-3 language code")
    parser.add_argument(
        "-b",
        help="non-words output file",
        required=True,
    )
    parser.add_argument(
        "-w",
        help="words output file",
        required=True,
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    try:
        alphabet_ = alphabet.load(args.language)
    except alphabet.UnsupportedLanguage:
        exit(f"unsupported language: {args.language}")

    try:
        with open(args.b, "w") as file_b:
            with open(args.w, "w") as file_w:
                while line := input():
                    if alphabet.is_word(alphabet_, line):
                        print(line, file=file_w)
                    else:
                        print(line, file=file_b)
    except EOFError:
        pass


if __name__ == "__main__":
    main()
