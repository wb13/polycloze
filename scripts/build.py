"""Course builder."""

from argparse import ArgumentParser, Namespace, RawDescriptionHelpFormatter
from graphlib import TopologicalSorter  # pylint: disable=unused-import
import sys

from . import dependency, task
from .dependency import execute, Task  # pylint: disable=unused-import
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
        "-B",
        dest="build_always",
        action="store_true",
        help="build all targets unconditionally",
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

    if args.build_always:
        dependency.BUILD_ALWAYS = True

    # Build dependency graph
    deps: "TopologicalSorter[Task]" = TopologicalSorter()
    deps.add(task.decompress_links, task.download_latest)
    deps.add(task.decompress_sentences, task.download_latest)
    deps.add(task.prepare_sentences, task.decompress_sentences)

    for lang in sorted(set(l1s + l2s)):
        deps.add(task.language_tokenizer(lang), task.prepare_sentences)

    for lang1 in l1s:
        for lang2 in l2s:
            if lang1 < lang2:
                deps.add(
                    task.translation_mapper(lang1, lang2),
                    task.language_tokenizer(lang1),
                    task.language_tokenizer(lang2),
                    task.decompress_links,
                )

    for lang1 in l1s:
        for lang2 in l2s:
            if lang1 != lang2:
                deps.add(
                    task.course_builder(lang1, lang2),
                    task.translation_mapper(
                        min(lang1, lang2),
                        max(lang1, lang2),
                    ),
                    task.language_tokenizer(lang1),
                    task.language_tokenizer(lang2),
                )

    execute(deps)


if __name__ == "__main__":
    main(parse_args())
