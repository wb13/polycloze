"""Rule-based classifiers for blacklisting "non-words"."""

from dataclasses import dataclass, field


@dataclass
class Language:
    alphabet: set[str]
    symbols: set[str] = field(default_factory=set)


# source for cyo:
# https://web.archive.org/web/20120403120048/http://www.cuyonon.org/clcp8.html

languages = {
    "cyo": Language(
        alphabet=set("abdeghiklmnoprstwy'"),
        # symbols=set("_.'"),
    ),
    "deu": Language(
        alphabet=set("abcdefghijklmnopqrstuvwxyzäéöüß"),
        symbols=set("-.'"),
    ),
    "eng": Language(
        alphabet=set("abcdefghijklmnopqrstuvwxyz"),
        symbols=set("-.'"),
    ),
    "fra": Language(alphabet=set("abcdefghijklmnopqrstuvwxyzéàèùâêîôûëïüÿçñ")),
    "ita": Language(alphabet=set("abcdefghilmnopqrstuvzéóàèìòùî")),
    "spa": Language(
        alphabet=set("abcdefghijklmnñopqrstuvwxyzáéíóúü"),
        symbols=set("-."),
    ),
    "tgl": Language(
        alphabet=set("abcdefghijklmnñopqrstuvwxyzáàâéèêëíìîóòôúùû"),
        symbols=set("-.'"),
    ),
}


class UnsupportedLanguage(Exception):
    pass


def is_word(language: Language, word: str) -> bool:
    word = word.casefold()
    if word[0] not in language.alphabet:
        return False
    return all(a in language.alphabet or a in language.symbols for a in word)


def load_alphabet(language: str) -> Language:
    try:
        return languages[language]
    except KeyError as exc:
        raise UnsupportedLanguage(language) from exc
