#!/usr/bin/env bash

# Copyright (c) 2022 Levi Gruspe
# License: MIT, or APGLv3 or later

# Activate dev environment

sudo docker build -t dev-polycloze .
container="$(sudo docker run -dit dev-polycloze)"

# Change prompt.
prompt="$PS1"
PS1="(${container:0:4}) $prompt"

deactivate() {
	echo "Deactivating environment..."
	unset -f deactivate
	unset -f run
	PS1="$prompt"
	sudo docker stop "$container"
}

run() {
	sudo docker exec "$container" sh -c 'rm -rf "/src/*"'
	sudo docker cp . "$container:/src"
	sudo docker exec "$container" "$@"
}
