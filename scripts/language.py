from dataclasses import dataclass, field


@dataclass
class Language:
    code: str
    name: str

    # e.g. ("de", "German") for spacy.lang.de.German
    spacy_path: tuple[str, str]

    alphabet: set[str]
    symbols: set[str] = field(default_factory=set)


languages = {}

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
