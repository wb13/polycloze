"""Applies migration scripts."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
import re
from sqlite3 import Connection, connect
import sys


def user_version(con: Connection) -> int:
    """Get database user_version."""
    res = con.execute("pragma user_version")
    return int(res.fetchone()[0])


def file_version(script: Path) -> int | None:
    """Get file version of script from filename."""
    pattern = "([0-9]+).*.sql"
    result = re.match(pattern, script.name)
    if not result:
        return None

    groups = result.groups()
    if len(groups) < 1:
        return None
    return int(groups[0])


def script_version(script: Path) -> int | None:
    """Get script version from text."""
    pattern = r"pragma\s*?user_version\s*?=\s*?([0-9]+)"
    text = script.read_text()

    result = re.search(pattern, text, flags=re.IGNORECASE)
    if not result:
        return None

    groups = result.groups()
    if len(groups) < 1:
        return None
    return int(groups[0])


def check_script(script: Path) -> bool:
    """Check if file and script versions match."""
    return file_version(script) == script_version(script)


class InvalidScriptVersion(Exception):
    """E.g. version in filename and body don't match."""


def check_scripts(migrations: Path) -> list[tuple[int, Path]]:
    """Check scripts in migrations directory.

    Returns list of scripts sorted by version number.
    Raises InvalidScriptVersion if there's an invalid script.
    """
    scripts = []
    for script in migrations.glob("*.sql"):
        file_ver = file_version(script)
        script_ver = script_version(script)
        if file_ver is None or script_ver is None:
            print("Not a migration script:", str(script))
            continue

        if file_ver != script_ver:
            raise InvalidScriptVersion(script)

        scripts.append((file_ver, script))
    return sorted(scripts)


def migrate(con: Connection, migrations: Path) -> None:
    """Apply migration scripts."""
    for version, script in check_scripts(migrations):
        if version <= user_version(con):
            print("Skipping:", str(script))
        else:
            print("Applying:", str(script))
            con.executescript(script.read_text())


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("database", type=Path, help="database file")
    parser.add_argument("migrations", type=Path, help="migrations directory")
    return parser.parse_args()


def main(args: Namespace) -> None:
    try:
        with connect(args.database) as con:
            migrate(con, args.migrations)
    except InvalidScriptVersion as exc:
        name = str(exc.args[0])
        message = f"Filename and script version numbers don't match: {name}"
        sys.exit(message)


if __name__ == "__main__":
    main(parse_args())
