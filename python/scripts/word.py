"""Defines how words are stored in the database."""


class Word(str):
    """Canonical representation of a word (case-folded and soft hyphens
    stripped out."""
    def __new__(cls, content: str) -> "Word":
        content = content.replace("\xAD", "")
        return super().__new__(cls, content.casefold())
