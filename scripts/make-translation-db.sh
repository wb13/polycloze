#!/usr/bin/env bash

translations_csv="$1"
output_db="$2"

if [[ "$translations_csv" == "" ]]; then
	echo "missing translations.csv file"
	exit 1
fi

if [[ "$output_db" == "" ]]; then
	echo "missing path to output sqlite file"
	exit 1
fi

sql="create table if not exists translation (source, target);

.mode csv
.import '$translations_csv' translation

update translation
set source = cast(source as int),
		target = cast(target as int);
";

echo "$sql" | sqlite3 "$output_db"
