#!/usr/bin/env bash

# Print ID of sentences in specified language.

if [[ "$1" == "" ]]; then
	echo "missing ISO 639-3 language code"
	exit 1
fi

latest=$(find build/tatoeba/sentences.*.csv | sort -r | head -n 1)

grep -h -P "\t$1\t" "$latest" | sed "s/\([0-9]\+\)\t.\+/\1/g"
