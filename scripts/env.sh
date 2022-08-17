#!/usr/bin/env bash

# Activate dev environment

sudo docker build -t dev-polycloze-data .
container="$(sudo docker run -dit dev-polycloze-data)"

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

run-it() {
	sudo docker exec "$container" sh -c 'rm -rf "/src/*"'
	sudo docker cp . "$container:/src"
	sudo docker exec -it "$container" "$@"
}
