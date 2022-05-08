#!/usr/bin/env bash

# Download sentence and translation data from tatoeba.

download-sentences() {
	wget https://downloads.tatoeba.org/exports/sentences.tar.bz2
	tar -xvf sentences.tar.bz2
	rm sentences.tar.bz2
	mv sentences.csv "build/tatoeba/sentences.$(date -I).csv"
}

download-translations() {
	wget https://downloads.tatoeba.org/exports/links.tar.bz2
	tar -xvf links.tar.bz2
	rm links.tar.bz2
	mv links.csv "build/tatoeba/links.$(date -I).csv"
}

mkdir -p build/tatoeba
download-sentences &
download-translations &
wait
