"""Course builder."""

from argparse import ArgumentParser, Namespace, RawDescriptionHelpFormatter
from concurrent.futures import ProcessPoolExecutor
from pathlib import Path
import sys
from tempfile import TemporaryDirectory

from . import download, mapper, migrate, partition, populate, tokenizer, untar
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


def build_language(lang: str) -> None:
    """Build files needed for language.

    - build/languages/{lang}/sentences.csv
    - build/languages/{lang}/words.csv
    """
    build = Path("build")
    lang_dir = build/"languages"/lang
    lang_dir.mkdir(parents=True, exist_ok=True)

    tokenizer.main(
        Namespace(
            language=lang,
            output=lang_dir,
            log=build/"logs"/"nonwords"/f"{lang}.txt",
            file=build/"sentences"/f"{lang}.tsv",
        ),
    )


def build_translations(lang1: str, lang2: str) -> None:
    """Build build/translations/{lang1}-{lang2}.csv.

    L1-L2 and L2-L1 use the same translation file, so only L1-L2 where L1 < L2
    is built.
    """
    assert lang1 < lang2

    build = Path("build")
    translations = build/"translations"
    translations.mkdir(parents=True, exist_ok=True)
    mapper.main(
        Namespace(
            l1=build/"sentences"/f"{lang1}.tsv",
            l2=build/"sentences"/f"{lang2}.tsv",
            links=build/"tatoeba"/"links.csv",
            output=build/"translations"/f"{lang1}-{lang2}.csv",
        ),
    )


def build_course(lang1: str, lang2: str) -> None:
    """Build lang1 -> lang2 course.

    Assumes the following requirements have been built:
    - build/languages/{lang1}
    - build/languages/{lang2}
    - build/translations/{lang1}-{lang2}.csv (or {lang2}-{lang1}.csv)
    """
    assert lang1 != lang2
    build = Path("build")

    with TemporaryDirectory() as tmpname:
        tmp = Path(tmpname)
        database = tmp/"scratch.db"

        # Apply migrations in empty database file.
        migrate.main(
            Namespace(
                database=database,
                migrations=Path(__file__).parent/"migrations",
            ),
        )

        # Populate database
        if lang1 < lang2:
            populate.main(
                Namespace(
                    reversed=True,
                    database=database,
                    l1=build/"languages"/lang1,
                    l2=build/"languages"/lang2,
                    translations=build/"translations"/f"{lang1}-{lang2}.csv",
                ),
            )
        else:
            populate.main(
                Namespace(
                    reversed=False,
                    database=database,
                    l1=build/"languages"/lang1,
                    l2=build/"languages"/lang2,
                    translations=build/"translations"/f"{lang2}-{lang1}.csv",
                ),
            )

        # Replace existing course with new one.
        database.replace(build/"courses"/f"{lang1}-{lang2}.db")


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

    # Unarchive downloaded data.
    untar.main(Namespace(links=None, sentences=None))

    # Partition build/tatoeba/sentences.csv, output in build/sentences/*
    print("Processing sentences...")
    partition.main(
        Namespace(
            out=build/"sentences",
            file=build/"tatoeba"/"sentences.csv",
        ),
    )

    # Build languages, sentences, etc.
    print("Tokenizing words...")
    with ProcessPoolExecutor() as executor:
        futures = [executor.submit(build_language, lang) for lang in l1s + l2s]
        for future in futures:
            future.result()

    # Build translations.
    print("Processing translations...")
    with ProcessPoolExecutor() as executor:
        futures = [
            executor.submit(build_translations, lang1, lang2)
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
