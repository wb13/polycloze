#!/usr/bin/env bash

# Format result of sentences.sh.

case "$1" in
	id)
		sed "s/\([0-9]\+\)\t.\+/\1/g"
		;;
	sentence)
		sed "s/.\+\t.\+\t\(.\+\)/\1/g"
		;;
	id-sentence)
		sed "s/\(.\+\)\t.\+\t\(.\+\)/\1\t\2/g"
		;;

	*)
		echo "invalid format: id, sentence or id-sentence only"
		exit 1
esac
