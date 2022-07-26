# polycloze-data

Scripts for generating course files used by [polycloze](https://github.com/lggruspe/polycloze).

## Usage

```bash
# install requirements
python -m venv env
. env/bin/activate
pip install -r requirements/install.requirements.txt

# build course files
make

# install files into data directory
make install
```

## Supported languages

- deu (German)
- eng (English)
- spa (Spanish)
- tgl (Tagalog)

You can add more by making the appropriate changes in `scripts/language.py`.

## Licenses

All scripts in this repository are available under the [MIT license](./LICENSE).

Sentences and translation data are taken from [Tatoeba](https://tatoeba.org),
which are released under [CC BY 2.0 FR](https://creativecommons.org/licenses/by/2.0/fr).
