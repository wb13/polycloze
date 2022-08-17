#!/usr/bin/env bash

echo "Activating build environment..."
. scripts/env.sh

echo "Processing tatoeba sentences... (this can take a while)"
run sh -c "python -m scripts.partition build/sentences -f $(find build/tatoeba/sentences.*.csv | sort -r | head -n 1)"

echo "Building courses..."
run sh -c "make -j $(nproc) SHELL=/bin/bash"

echo "Exporting courses..."
run ./scripts/export.sh "$container" build/courses
deactivate

echo "Done :)"
echo "See build/courses"
