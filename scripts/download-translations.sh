#!/usr/bin/env bash

# Download translation data from tatoeba.

wget https://downloads.tatoeba.org/exports/links.tar.bz2
tar -xvf links.tar.bz2
rm links.tar.bz2

mkdir -p build/tatoeba
mv links.csv "build/tatoeba/links.$(date -I).csv"
