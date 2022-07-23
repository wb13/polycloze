"""Language definitions with heuristic word classifier."""

from dataclasses import dataclass, field


@dataclass
class Language:
    code: str
    name: str

    # e.g. ("de", "German") for spacy.lang.de.German
    spacy_path: tuple[str, str]

    alphabet: set[str]
    symbols: set[str] = field(default_factory=set)

    def is_word(self, word: str) -> bool:
        word = word.casefold()
        if word[0] not in self.alphabet:
            return False
        return all(a in self.alphabet or a in self.symbols for a in word)


languages = {}

# source for cyo:
# https://web.archive.org/web/20120403120048/http://www.cuyonon.org/clcp8.html
languages["cyo"] = Language(
    code="cyo",
    name="Cuyonon",
    spacy_path=("tl", "Tagalog"),
    alphabet=set("abdeghiklmnoprstwy'"),
)

languages["deu"] = Language(
    code="deu",
    name="German",
    spacy_path=("de", "German"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzäéöüß"),
    symbols=set("-.'"),
)

languages["eng"] = Language(
    code="eng",
    name="English",
    spacy_path=("en", "English"),
    alphabet=set("abcdefghijklmnopqrstuvwxyz"),
    symbols=set("-.'"),
)

languages["spa"] = Language(
    code="spa",
    name="Spanish",
    spacy_path=("es", "Spanish"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáéíóúü"),
    symbols=set("-."),
)

languages["tgl"] = Language(
    code="tgl",
    name="Tagalog",
    spacy_path=("tl", "Tagalog"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáàâéèêëíìîóòôúùû"),
    symbols=set("-.'"),
)
