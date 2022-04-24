#!/usr/bin/env bash

./scripts/download-sentences.sh &
./scripts/download-translations.sh &

wait
