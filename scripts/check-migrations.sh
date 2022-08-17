#!/usr/bin/env bash

# Check migration scripts.

migrations="$1"

file_version() {
	basename "$1" | sed 's/\([0-9]\+\).*.sql/\1/g'
}

script_version() {
	rg -Noir "\$1" 'pragma.*user_version.*=.*(\d+)' "$1"
}

check() {
	local scripts
	local has_error

	if [[ "$migrations" == "" ]]; then
		echo "missing migrations directory"
		exit 1
	fi

	has_error=0
	scripts=$(find "$migrations"/*.sql | sort -V)
	for script in $scripts
	do
		if [[ $(script_version "$script") -ne $(file_version "$script") ]]; then
			echo "Filename and script versions don't match: $script"
			has_error=1
		fi
	done

	exit $has_error
}

check
