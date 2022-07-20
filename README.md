# polycloze-data

All the sentences are from [Tatoeba](https://tatoeba.org), which are released
under [CC BY 2.0][cc].

## Changes

- `sentences.csv`: list of tokenized sentences
- `words.csv`: list of "words" sorted by frequency

These are also released under [CC BY 2.0][cc].

## Usage

### blacklist

```bash
python -m scripts.blacklist spa < words.txt > blacklist.txt
```

## Supported languages

- deu (German)
- eng (English)
- spa (Spanish)
- tgl (Tagalog)

## How to add languages

Modify `scripts/languages.py`, `scripts/populate.py`, and `scripts/alphabet.py`.


[cc]: https://creativecommons.org/licenses/by/2.0
