# poly<u>cloze</u>

A tool for studying foreign language vocabulary using cloze tests.

[Demo](https://polycloze-demo.herokuapp.com/)

## Features

**Adaptive flashcards scheduler**
:	Polycloze's flashcards scheduler adjusts the spacing between reviews based on the student's performance.
It also estimates the student's vocabulary level, so advanced learners don't have to see words they already know.

**[Digraphs](./docs/digraphs.md)**
: Even if your keyboard doesn't have all the needed keys, you can use digraphs to enter any symbol in your target language.
For example, `\a:` automatically gets turned into `ä`.

**Uses SQLite**
: This makes it easy to track your vocabulary.
Even course files are stored as SQLite databases.
By default, course files are saved in `~/.local/share/polycloze` and review data are saved in `~/.local/state/polycloze`.

## Usage

```bash
# Install everything needed to build front-end.
make init

# Run server.
make run

# Open in browser.
xdg-open http://localhost:3000
```

## Licenses

Copyright (C) 2022 Levi Gruspe

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

---

The scripts in [./scripts](./scripts) and
[./database/migrations](./database/migrations) are also available under the
terms of the MIT license.
