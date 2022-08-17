#!/usr/bin/env bash

# Export docker-built courses out of container.

container="$1"
output="$2"

if [[ $container = "" ]]; then
	echo "missing container argument"
	exit 1
fi

if [[ $output = "" ]]; then
	echo "missing output directory"
	exit 1
fi

sudo docker cp "$container:/src/build/courses" "$output"
