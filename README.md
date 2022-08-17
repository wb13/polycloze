# polycloze-data

Scripts for generating course files used by [polycloze](https://github.com/lggruspe/polycloze).

## Usage

```bash
# Install requirements in a virtual environment.
python -m venv env
. env/bin/activate
pip install -r requirements/install.requirements.txt

# Build course files.
make

# Copy files into data directory (~/.local/share/polycloze).
make install
```

If you only want to build a specific course, then you can use the course name as the `make` target.

```bash
# Instead of
make

# You can do
make eng-spa	# build English -> Spanish course only
```

## Building course files inside Docker

Alternatively, you can build course files in a Docker container.
The only requirements are bash and docker.

```bash
# Download tatoeba data.
./scripts/download.sh

# Build all courses.
./scripts/build-in-docker.sh
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
