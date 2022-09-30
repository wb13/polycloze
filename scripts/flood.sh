#!/usr/bin/env bash

# Copyright (c) 2022 Levi Gruspe
# License: MIT, or AGPLv3 or later

duration="2"

if [[ "$1" != "" ]]; then
	duration="$1"
fi

while true; do
	if ! curl http://localhost:3000/eng/spa; then
		exit 1
	fi
	sleep "$duration"
done
