#!/usr/bin/env bash

# Copyright (c) 2022 Levi Gruspe
# License: MIT, or AGPLv3 or later

# For manual database testing

if [[ "$1" == "" ]]; then
	echo "missing L1 ISO 639-3 language code"
	exit 1
fi

if [[ "$2" == "" ]]; then
	echo "missing L2 ISO 639-3 language code"
	exit 1
fi

if [[ "$1" == "$2" ]]; then
	echo "invalid language pair"
	exit 1
fi

pair() {
	if [[ "$1" < "$2" ]]; then
		echo "$1-$2"
	else
		echo "$2-$1"
	fi
}

sql=".read database/migrations/1_init_schema.up.sql
.read database/migrations/2_add_interval_table.up.sql
.read database/migrations/3_add_student_table.up.sql

attach database '$HOME/.local/share/polycloze/$(pair "$1" "$2").db' as course;

.timer on
"

tmpfile="$(mktemp)"
trap 'rm -f "$tmpfile"' EXIT

echo "$sql" > "$tmpfile"
sqlite3 -init "$tmpfile"
