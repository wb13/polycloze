#!/usr/bin/env bash

# Print sentences in the specified language (in ISO 639-3).

if [[ "$1" == "" ]]; then
	echo "missing ISO 639-3 language code"
	exit 1
fi

latest=$(find sentences/sentences.*.csv | sort -r | head -n 1)
grep -h -P "\t$1\t" "$latest" | sed "s/.\+\t.\+\t\(.\+\)/\1/g"
