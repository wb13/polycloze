"""Language definitions with heuristic word classifier."""

from dataclasses import dataclass, field
from importlib import import_module
import typing as t

if t.TYPE_CHECKING:
    from spacy.tokenizer import Tokenizer   # type: ignore


def import_tokenizer(module: str, name: str) -> "Tokenizer":
    """Import tokenizer from spacy.

    Usage example:
    ```
    tokenize = import_tokenizer("spacy.lang.en", "English")
    tokenize(text)
    ```
    """
    mod = import_module(module)
    nlp = getattr(mod, name)()
    return nlp.tokenizer


@dataclass
class Language:
    code: str
    name: str
    bcp47: str  # goes in html lang attribute

    tokenizer: t.Callable[[], "Tokenizer"]
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

languages["dan"] = Language(
    code="dan",
    name="Danish",
    bcp47="da",
    tokenizer=lambda: import_tokenizer("spacy.lang.da", "Danish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzæøå"),
)

languages["deu"] = Language(
    code="deu",
    name="German",
    bcp47="de",
    tokenizer=lambda: import_tokenizer("spacy.lang.de", "German"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzäéöüß"),
    symbols=set("-.'"),
)

languages["eng"] = Language(
    code="eng",
    name="English",
    bcp47="en",
    tokenizer=lambda: import_tokenizer("spacy.lang.en", "English"),
    alphabet=set("abcdefghijklmnopqrstuvwxyz"),
    symbols=set("-.'"),
)

languages["fin"] = Language(
    code="fin",
    name="Finnish",
    bcp47="fi",
    tokenizer=lambda: import_tokenizer("spacy.lang.fi", "Finnish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzåäöšž"),
)

languages["fra"] = Language(
    code="fra",
    name="French",
    bcp47="fr",
    tokenizer=lambda: import_tokenizer("spacy.lang.fr", "French"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzàâæçéèêëîïôœùûüÿ"),
)

languages["hrv"] = Language(
    code="hrv",
    name="Croatian",
    bcp47="hr",
    tokenizer=lambda: import_tokenizer("spacy.lang.hr", "Croatian"),
    alphabet=set("abcčćdđefghijklmnoprsštuvzž"),
)

languages["ita"] = Language(
    code="ita",
    name="Italian",
    bcp47="it",
    tokenizer=lambda: import_tokenizer("spacy.lang.it", "Italian"),
    alphabet=set("abcdefghilmnopqrstuvzàèéìíîòóùú"),
)

languages["lit"] = Language(
    code="lit",
    name="Lithuanian",
    bcp47="lt",
    tokenizer=lambda: import_tokenizer("spacy.lang.lt", "Lithuanian"),
    alphabet=set("aąbcčdeęėfghiįyjklmnoprsštuųūvzž"),
)

languages["nld"] = Language(
    code="nld",
    name="Dutch",
    bcp47="nl",
    tokenizer=lambda: import_tokenizer("spacy.lang.nl", "Dutch"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzĳäëïöüáéíóú"),
)

languages["nob"] = Language(
    code="nob",
    name="Norwegian Bokmål",
    bcp47="nb",
    tokenizer=lambda: import_tokenizer("spacy.lang.nb", "Norwegian"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzæøå"),
)

languages["pol"] = Language(
    code="pol",
    name="Polish",
    bcp47="pl",
    tokenizer=lambda: import_tokenizer("spacy.lang.pl", "Polish"),
    alphabet=set("aąbcćdeęfghijklłmnńoópqrsśtuvwxyzźż"),
)

languages["por"] = Language(
    code="por",
    name="Portuguese",
    bcp47="pt",
    tokenizer=lambda: import_tokenizer("spacy.lang.pt", "Portuguese"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzáâãàçéêíóôõú"),
)

languages["ron"] = Language(
    code="ron",
    name="Romanian",
    bcp47="ro",
    tokenizer=lambda: import_tokenizer("spacy.lang.ro", "Romanian"),
    alphabet=set("aăâbcdefghiîjklmnopqrsştţuvwxyz"),
)

languages["spa"] = Language(
    code="spa",
    name="Spanish",
    bcp47="es",
    tokenizer=lambda: import_tokenizer("spacy.lang.es", "Spanish"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáéíóúü"),
    symbols=set("-."),
)

languages["swe"] = Language(
    code="swe",
    name="Swedish",
    bcp47="sv",
    tokenizer=lambda: import_tokenizer("spacy.lang.sv", "Swedish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzåäöáüè"),
)

languages["tgl"] = Language(
    code="tgl",
    name="Tagalog",
    bcp47="tl",
    tokenizer=lambda: import_tokenizer("spacy.lang.tl", "Tagalog"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáàâéèêëíìîóòôúùû"),
    symbols=set("-.'"),
)
