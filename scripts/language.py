"""Language definitions with heuristic word classifier."""

from dataclasses import dataclass, field


@dataclass
class Language:
    code: str
    name: str
    bcp47: str  # goes in html lang attribute

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

# source for bcp47 codes:
# https://www.iana.org/assignments/language-subtag-registry/

# source for cyo:
# https://web.archive.org/web/20120403120048/http://www.cuyonon.org/clcp8.html
languages["cyo"] = Language(
    code="cyo",
    name="Cuyonon",
    bcp47="cyo",
    spacy_path=("tl", "Tagalog"),
    alphabet=set("abdeghiklmnoprstwy'"),
)

languages["dan"] = Language(
    code="dan",
    name="Danish",
    bcp47="da",
    spacy_path=("da", "Danish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzæøå"),
)

languages["deu"] = Language(
    code="deu",
    name="German",
    bcp47="de",
    spacy_path=("de", "German"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzäéöüß"),
    symbols=set("-.'"),
)

languages["eng"] = Language(
    code="eng",
    name="English",
    bcp47="en",
    spacy_path=("en", "English"),
    alphabet=set("abcdefghijklmnopqrstuvwxyz"),
    symbols=set("-.'"),
)

languages["fin"] = Language(
    code="fin",
    name="Finnish",
    bcp47="fi",
    spacy_path=("fi", "Finnish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzåäöšž"),
)

languages["fra"] = Language(
    code="fra",
    name="French",
    bcp47="fr",
    spacy_path=("fr", "French"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzàâæçéèêëîïôœùûüÿ"),
)

languages["hrv"] = Language(
    code="hrv",
    name="Croatian",
    bcp47="hr",
    spacy_path=("hr", "Croatian"),
    alphabet=set("abcčćdđefghijklmnoprsštuvzž"),
)

languages["ita"] = Language(
    code="ita",
    name="Italian",
    bcp47="it",
    spacy_path=("it", "Italian"),
    alphabet=set("abcdefghilmnopqrstuvzàèéìíîòóùú"),
)

languages["lit"] = Language(
    code="lit",
    name="Lithuanian",
    bcp47="lt",
    spacy_path=("lt", "Lithuanian"),
    alphabet=set("aąbcčdeęėfghiįyjklmnoprsštuųūvzž"),
)

languages["nld"] = Language(
    code="nld",
    name="Dutch",
    bcp47="nl",
    spacy_path=("nl", "Dutch"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzĳäëïöüáéíóú"),
)

languages["nob"] = Language(
    code="nob",
    name="Norwegian Bokmål",
    bcp47="nb",
    spacy_path=("nb", "Norwegian"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzæøå"),
)

languages["pol"] = Language(
    code="pol",
    name="Polish",
    bcp47="pl",
    spacy_path=("pl", "Polish"),
    alphabet=set("aąbcćdeęfghijklłmnńoópqrsśtuvwxyzźż"),
)

languages["por"] = Language(
    code="por",
    name="Portuguese",
    bcp47="pt",
    spacy_path=("pt", "Portuguese"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzáâãàçéêíóôõú"),
)

languages["ron"] = Language(
    code="ron",
    name="Romanian",
    bcp47="ro",
    spacy_path=("ro", "Romanian"),
    alphabet=set("aăâbcdefghiîjklmnopqrsştţuvwxyz"),
)

languages["spa"] = Language(
    code="spa",
    name="Spanish",
    bcp47="es",
    spacy_path=("es", "Spanish"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáéíóúü"),
    symbols=set("-."),
)

languages["swe"] = Language(
    code="swe",
    name="Swedish",
    bcp47="sv",
    spacy_path=("sv", "Swedish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzåäöáüè"),
)

languages["tgl"] = Language(
    code="tgl",
    name="Tagalog",
    bcp47="tl",
    spacy_path=("tl", "Tagalog"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáàâéèêëíìîóòôúùû"),
    symbols=set("-.'"),
)


if __name__ == "__main__":
    for code in languages:
        print(code)
