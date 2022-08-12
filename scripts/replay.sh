#!/usr/bin/env bash

# Copyright (c) 2022 Levi Gruspe
# License: MIT, or AGPLv3 or later

lang="spa"

if [[ "$1" != "" ]]; then
	lang="$1"
fi

./build/replay -v "$HOME/.local/state/polycloze/logs/user/$lang.log" test.db
mv test.db "$HOME/.local/state/polycloze/reviews/user/$lang.db"
