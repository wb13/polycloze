# polycloze-data

Scripts for generating course files used by [polycloze](https://github.com/lggruspe/polycloze).

## Usage

Note: requires a Rust compiler for Japanese.

```bash
# Install requirements in a virtual environment.
python -m venv env
. env/bin/activate
pip install -r requirements/install.requirements.txt

# Build all course files.
python -m scripts.build

# Install into data directory.
cp -r ./build/polycloze ~/.local/share
```

You can also specify a course to build.
For example:

```bash
# Build English -> Spanish course
python -m scripts.build eng spa
```

See `python -m scripts.build -h` for details.

## Supported languages

- dan (Danish)
- deu (German)
- eng (English)
- epo (Esperanto)
- fin (Finnish)
- fra (French)
- hrv (Croatian)
- ita (Italian)
- jpn (Japanese)
- lit (Lithuanian)
- nld (Dutch)
- nob (Norwegian Bokmål)
- pol (Polish)
- por (Portuguese)
- ron (Romanian)
- spa (Spanish)
- swe (Swedish)
- tgl (Tagalog)
- tok (toki pona)

You can add more by making the appropriate changes in `scripts/language.py`.

## Licenses

These scripts are available under the [MIT license](./LICENSE).

Sentences and translation data are taken from [Tatoeba](https://tatoeba.org),
which are released under [CC BY 2.0 FR](https://creativecommons.org/licenses/by/2.0/fr).
