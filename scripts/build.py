"""Course builder."""

from argparse import ArgumentParser, Namespace, RawDescriptionHelpFormatter
from concurrent.futures import ProcessPoolExecutor
from pathlib import Path
from shutil import move
import sys
from tempfile import TemporaryDirectory

from . import download, migrate, populate, task
from .dependency import is_outdated
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


def build_course(lang1: str, lang2: str) -> None:
    """Build lang1 -> lang2 course.

    Assumes the following requirements have been built:
    - build/languages/{lang1}
    - build/languages/{lang2}
    - build/translations/{lang1}-{lang2}.csv (or {lang2}-{lang1}.csv)
    """
    assert lang1 != lang2
    build = Path("build")
    courses = build/"courses"
    courses.mkdir(parents=True, exist_ok=True)

    course = build/"courses"/f"{lang1}-{lang2}.db"
    translations = (
        build/"translations"/f"{lang1}-{lang2}.csv"
        if lang1 < lang2
        else build/"translations"/f"{lang2}-{lang1}.csv"
    )

    if not is_outdated(
        [course],
        [build/"languages"/lang1, build/"languages"/lang2, translations],
    ):
        return

    with TemporaryDirectory() as tmpname:
        tmp = Path(tmpname)
        database = tmp/"scratch.db"

        # Apply migrations in empty database file.
        migrate.main(
            Namespace(
                database=database,
                migrations=Path(__file__).parent.parent/"migrations",
            ),
        )

        # Populate database
        populate.main(
            Namespace(
                reversed=lang1 < lang2,
                database=database,
                l1=build/"languages"/lang1,
                l2=build/"languages"/lang2,
                translations=translations,
            ),
        )

        # Replace existing course with new one.
        # shutil.move is used instead of Path.replace, because Path.replace
        # might raise OSError: Invalid cross-device link
        move(database, course)


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
    print("Tokenizing words...")
    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(task.language_tokenizer(lang))
            for lang in sorted(set(l1s + l2s))
        ]
        for future in futures:
            future.result()

    # Build translations.
    print("Processing translations...")
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
    print("Building courses...")
    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(build_course, lang1, lang2)
            for lang1 in l1s
            for lang2 in l2s
            if lang1 != lang2
        ]
        for future in futures:
            future.result()


if __name__ == "__main__":
    main(parse_args())
