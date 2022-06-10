#!/usr/bin/env bash

# For manual database testing

if [[ "$1" == "" ]]; then
	echo "missing L1 ISO 639-3 language code"
	exit 1
fi

if [[ "$2" == "" ]]; then
	echo "missing L2 ISO 639-3 language code"
	exit 1
fi

sql=".read database/migrations/1_init_schema.up.sql

attach database '$HOME/.local/share/polycloze/languages/$2.db' as l2;
attach database '$HOME/.local/share/polycloze/languages/$1.db' as l1;
attach database '$HOME/.local/share/polycloze/translations.db' as translation;

.timer on
"

tmpfile="$(mktemp)"
trap "rm -f '$tmpfile'" EXIT

echo "$sql" > "$tmpfile"
sqlite3 -init "$tmpfile"
