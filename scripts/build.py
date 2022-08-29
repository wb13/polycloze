"""Course builder."""

from argparse import ArgumentParser, Namespace, RawDescriptionHelpFormatter
from concurrent.futures import ProcessPoolExecutor
from pathlib import Path
import sys

from . import download, task
from .language import languages as supported_languages


class UnknownLanguage(Exception):
    """E.g. non ISO 639-3 language code or not supported."""


def parse_languages(languages: str) -> list[str]:
    """Parse languages str from command-line.

    Raises UnknownLanguage on error.
    """
    if languages == "_":
        return list(supported_languages)
    result = []
    for lang in languages.split(","):
        if lang in supported_languages:
            result.append(lang)
        else:
            raise UnknownLanguage(lang)
    return result


def parse_args() -> Namespace:
    description = "Build course files."
    epilog = """
examples:
  build.py eng
    Build all English -> * courses.

  build.py _ eng
    Build all * -> English courses.

  build.py eng fra,spa
    Build English -> French and English -> Spanish courses.
"""
    parser = ArgumentParser(
        description=description,
        epilog=epilog,
        formatter_class=RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "l1",
        default="_",
        nargs="?",
        help="default: '_' (all languages)",
    )
    parser.add_argument(
        "l2",
        default="_",
        nargs="?",
        help="default: '_' (all languages)",
    )
    return parser.parse_args()


def main(args: Namespace) -> None:
    try:
        l1s = parse_languages(args.l1)
        l2s = parse_languages(args.l2)
    except UnknownLanguage as exc:
        sys.exit(f"unknown language: {exc.args[0]}")

    build = Path("build")

    # Download latest data.
    download.main(
        Namespace(
            ls=False,
            downloads=build/"tatoeba",
        ),
    )

    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(task.decompress_links),
            executor.submit(task.decompress_sentences),
        ]
        for future in futures:
            future.result()

    task.prepare_sentences()

    # Build languages, sentences, etc.
    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(task.language_tokenizer(lang))
            for lang in sorted(set(l1s + l2s))
        ]
        for future in futures:
            future.result()

    # Build translations.
    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(task.translation_mapper(lang1, lang2))
            for lang1 in l1s
            for lang2 in l2s
            if lang1 < lang2
        ]
        for future in futures:
            future.result()

    # Build courses
    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(task.course_builder(lang1, lang2))
            for lang1 in l1s
            for lang2 in l2s
            if lang1 != lang2
        ]
        for future in futures:
            future.result()


if __name__ == "__main__":
    main(parse_args())
