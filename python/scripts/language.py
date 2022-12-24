"""Language definitions with heuristic word classifier."""

from dataclasses import dataclass, field
from importlib import import_module
import typing as t

from .word import Word

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


class CharacterRange(t.NamedTuple):
    """Unicode codepoints at the boundaries of a character range.

    `start` is included in the range, `end` is not.
    Assumes `start` <= `end`.
    """
    start: str
    end: str

    def contains(self, char: str) -> bool:
        """Check if the character range contains the given character."""
        return ord(self.start) <= ord(char) < ord(self.end)


@dataclass
class Language:
    code: str
    name: str
    bcp47: str  # goes in html lang attribute

    tokenizer: t.Callable[[], "Tokenizer"]
    alphabet: set[str]
    symbols: set[str] = field(default_factory=set)

    # If this field is not empty, overrides `alphabet` and `symbols`.
    character_ranges: set[CharacterRange] = field(default_factory=set)

    def is_word(self, word: Word) -> bool:
        if not word:
            return False

        if self.character_ranges:
            for char_range in self.character_ranges:
                pass

            for char in word:
                for char_range in self.character_ranges:
                    if char_range.contains(char):
                        return True

                if any(not r.contains(char) for r in self.character_ranges):
                    return False
            return True

        if word[0] not in self.alphabet:
            return False
        return all(a in self.alphabet or a in self.symbols for a in word)


languages = {}

# source for bcp47 codes:
# https://www.iana.org/assignments/language-subtag-registry/

languages["cat"] = Language(
    code="cat",
    name="Catalan",
    bcp47="ca",
    tokenizer=lambda: import_tokenizer("spacy.lang.ca", "Catalan"),
    alphabet=set("abcdefghijlmnopqrstuvxyzàéèíïóòúüçkw"),
    symbols=set("-'0123456789"),
)

languages["dan"] = Language(
    code="dan",
    name="Danish",
    bcp47="da",
    tokenizer=lambda: import_tokenizer("spacy.lang.da", "Danish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzæøåé"),
    # é for café, diarré, idé
)

languages["deu"] = Language(
    code="deu",
    name="German",
    bcp47="de",
    tokenizer=lambda: import_tokenizer("spacy.lang.de", "German"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzàäéöüß"),
    symbols=set("-.'0123456789"),
    # Included
    # à for voilà, à la carte, à jour, Déjà-vu
    # é for Café

    # Excluded (because names):
    # á for Bogotá, Guzmán
    # ã for São
    # ç for Curaçao
    # í for Medellín, Brasília
    # ñ for Piña Colada
    # ó for Córdoba
    # ô for Rhône

    # Excluded (too few examples):
    # è for Crème, Gruyère
    # ê for Crêpe
    # û for Croûton
)

languages["ell"] = Language(
    code="ell",
    name="Greek",
    bcp47="el",
    tokenizer=lambda: import_tokenizer("spacy.lang.el", "Greek"),
    alphabet=set("αβγδεζηθικλμνξοπρσςτυφχψω"),
    symbols=set(","),
)

languages["eng"] = Language(
    code="eng",
    name="English",
    bcp47="en",
    tokenizer=lambda: import_tokenizer("spacy.lang.en", "English"),
    alphabet=set("abcdefghijklmnopqrstuvwxyz"),
    symbols=set("-.'0123456789"),
)

languages["epo"] = Language(
    code="epo",
    name="Esperanto",
    bcp47="eo",
    tokenizer=lambda: import_tokenizer("spacy.language", "Language"),
    alphabet=set("abcĉdefgĝhĥijĵklmnoprsŝtuŭvz"),
    symbols=set("-0123456789"),
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
    alphabet=set("abcdefghijklmnopqrstuvwxyzàâæçèéêëîïôùûüÿœ"),
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
    # j, k, w, x and y excluded because they only appear in loanwords?
)

# https://www.localizingjapan.com/blog/2012/01/20/regular-expressions-for-japanese-text/    # noqa; pylint: disable=line-too-long
languages["jpn"] = Language(
    code="jpn",
    name="Japanese",
    bcp47="ja",
    tokenizer=lambda: import_tokenizer("spacy.lang.ja", "Japanese"),
    alphabet=set(),

    character_ranges={
        # Hiragana
        CharacterRange("\u3041", "\u3096"),

        # Full-width katakana
        CharacterRange("\u30a0", "\u30ff"),

        # Kanji
        CharacterRange("\u3400", "\u4db5"),
        CharacterRange("\u4e00", "\u9fcb"),
        CharacterRange("\uf900", "\ufa6a"),

        # Kanji radicals
        CharacterRange("\u2e80", "\u2fd5"),

        # Katakana and punctuation (half-width)
        CharacterRange("\uff5f", "\uff9f"),

        # Symbols and punctuation
        CharacterRange("\u3000", "\u303f"),

        # Misc.
        CharacterRange("\u31f0", "\u31ff"),
        CharacterRange("\u3220", "\u3243"),
        CharacterRange("\u3280", "\u337f"),

        # Alphanumeric and punctuation (full-width)
        CharacterRange("\uff01", "\uff5e"),
    },
)

languages["lit"] = Language(
    code="lit",
    name="Lithuanian",
    bcp47="lt",
    tokenizer=lambda: import_tokenizer("spacy.lang.lt", "Lithuanian"),
    alphabet=set("aąbcčdeęėfghiįyjklmnoprsštuųūvzž"),
)

languages["mkd"] = Language(
    code="mkd",
    name="Macedonian",
    bcp47="mk",
    tokenizer=lambda: import_tokenizer("spacy.lang.mk", "Macedonian"),
    alphabet=set("абвгдѓежзѕијклљмнњопрстќуфхцчџшѐѝč"),
    symbols=set("'"),
)

languages["nld"] = Language(
    code="nld",
    name="Dutch",
    bcp47="nl",
    tokenizer=lambda: import_tokenizer("spacy.lang.nl", "Dutch"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzĳäëïöüáéíóú'"),
    # ' for contractions
    # Dutch uses diaeresis and acute accents?
    # Grave accents are only for french loanwords?
    # ç for curaçao, Française, façade
)

languages["nob"] = Language(
    code="nob",
    name="Norwegian Bokmål",
    bcp47="nb",
    tokenizer=lambda: import_tokenizer("spacy.lang.nb", "Norwegian"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzæøåôé"),
    # ô for fôr, dyrefôr
    # accute accents: allé, diaré, kafé, idé, entré, komité, kupé
    # moské, supé, trofé, diskré
    # ç for provençalsk (loan words only; excluded)
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
    alphabet=set("abcdefghijlmnopqrstuvwxyzáâãàçéêíóôõú"),
    # w for watt
    # y for hobby
    # è for ampère (excluded because no other words found)
    # ü excluded because no longer used?
)

languages["ron"] = Language(
    code="ron",
    name="Romanian",
    bcp47="ro",
    tokenizer=lambda: import_tokenizer("spacy.lang.ro", "Romanian"),
    alphabet=set("aăâbcdefghiîjklmnopqrsștțuvwxyzşţ"),
    # NOTE ş and ţ are used because ș and ț are hard to type
)

languages["rus"] = Language(
    code="rus",
    name="Russian",
    bcp47="ru",
    tokenizer=lambda: import_tokenizer("spacy.lang.ru", "Russian"),
    alphabet=set("бвгджзклмнпрстфхцчшщаеёиоуыэюяйьъ"),
)

languages["spa"] = Language(
    code="spa",
    name="Spanish",
    bcp47="es",
    tokenizer=lambda: import_tokenizer("spacy.lang.es", "Spanish"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáéíóúü"),

    # Space included because strings like "EE. UU." get tokenized as one word.
    symbols=set("-.'0123456789 "),
)

languages["swe"] = Language(
    code="swe",
    name="Swedish",
    bcp47="sv",
    tokenizer=lambda: import_tokenizer("spacy.lang.sv", "Swedish"),
    alphabet=set("abcdefghijklmnopqrstuvwxyzåäöáüèé"),
    # é for words like idé and armé (integrated; included)
    # à for french loanwords (non-integrated; excluded)
)

languages["tgl"] = Language(
    code="tgl",
    name="Tagalog",
    bcp47="tl",
    tokenizer=lambda: import_tokenizer("spacy.lang.tl", "Tagalog"),
    alphabet=set("abcdefghijklmnñopqrstuvwxyzáàâéèêëíìîóòôúùû'"),
    symbols=set("-.0123456789"),
)

languages["tok"] = Language(
    code="tok",
    name="toki pona",
    bcp47="tok",
    tokenizer=lambda: import_tokenizer("spacy.language", "Language"),
    alphabet=set("aeijklmnopstuw"),
)

languages["ukr"] = Language(
    code="ukr",
    name="Ukrainian",
    bcp47="uk",
    tokenizer=lambda: import_tokenizer("spacy.lang.uk", "Ukrainian"),
    alphabet=set("абвгґдеєжзиіїйклмнопрстуфхцчшщьюя'"),
)
