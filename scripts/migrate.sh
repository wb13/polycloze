#!/usr/bin/env bash

# Upgrades database to latest version.

database="$1"
migrations="$2"

user_version() {
	echo "pragma user_version" | sqlite3 "$1"
}

file_version() {
	basename "$1" | sed 's/\([0-9]\+\).*.sql/\1/g'
}

migrate() {
	local scripts

	if [[ "$database" == "" ]]; then
		echo "missing database file"
		exit 1
	fi

	if [[ "$migrations" == "" ]]; then
		echo "missing migrations directory"
		exit 1
	fi

	echo "$database is on version" "$(user_version "$database")"

	scripts=$(find "$migrations"/*.sql | sort -V)
	for script in $scripts
	do
		if [[ $(file_version "$script") -le $(user_version "$database") ]]; then
			echo "Skipping $script"
			continue
		fi

		echo "Applying $script"
		sqlite3 -bail "$database" < "$script"
		local rc=$?
		if [[ $rc -ne 0 ]]; then
			exit $rc
		fi
	done
}

migrate
