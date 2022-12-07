"""Defines how words are stored in the database."""


SOFT_HYPHEN = "\u00AD"
ZERO_WIDTH_SPACE = "\u200B"
NO_BREAK_SPACE = "\u00A0"


class Word(str):
    """Canonical representation of a word (case-folded and unneeded chars
    stripped out (e.g. soft-hyphens, zero-width spaces and non-breaking
    spaces).
    """
    def __new__(cls, content: str) -> "Word":
        # NOTE This operation is also performed in the `polycloze/text`
        # package, so any changes here should be reflected there as well.
        content = content.replace(SOFT_HYPHEN, "")

        while content.startswith(ZERO_WIDTH_SPACE):
            content.removeprefix(ZERO_WIDTH_SPACE)
        while content.endswith(ZERO_WIDTH_SPACE):
            content.removesuffix(ZERO_WIDTH_SPACE)

        while content.startswith(NO_BREAK_SPACE):
            content.removeprefix(NO_BREAK_SPACE)
        while content.endswith(NO_BREAK_SPACE):
            content.removesuffix(NO_BREAK_SPACE)

        return super().__new__(cls, content.casefold())
