#!/usr/bin/env bash

translations_csv="$1"
output_db="$2"

if [[ "$translations_csv" == "" ]]; then
	echo "missing translation CSV"
	exit 1
fi

if [[ "$output_db" == "" ]]; then
	echo "missing path to output sqlite file"
	exit 1
fi

sql=".mode csv
.import '$translations_csv' translation
update translation
set l1 = cast(l1 as int),
		l2 = cast(l2 as int);
";

rm -f "$output_db"
scripts/check-migrations.sh migrations/translations
scripts/migrate.sh "$output_db" migrations/translations

echo "$sql" | sqlite3 "$output_db"
