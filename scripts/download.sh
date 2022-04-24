#!/usr/bin/env bash

# Download sentences from tatoeba.

wget https://downloads.tatoeba.org/exports/sentences.tar.bz2
tar -xvf sentences.tar.bz2
rm sentences.tar.bz2

mkdir -p build/tatoeba
mv sentences.csv "build/tatoeba/sentences.$(date -I).csv"
