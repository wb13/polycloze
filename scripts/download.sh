#!/usr/bin/env bash

# Download sentences from tatoeba.

wget https://downloads.tatoeba.org/exports/sentences.tar.bz2
tar -xvf sentences.tar.bz2
rm sentences.tar.bz2
mv sentences.csv "sentences.$(date -I).csv"
