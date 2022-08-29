"""Applies migration scripts."""

from argparse import ArgumentParser, Namespace
from pathlib import Path
import re
from sqlite3 import Connection, connect
import sys
import typing as t


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


class MigrationScriptError(Exception):
    """E.g. bad or missing script version number."""
    def __init__(
        self,
        path: Path,
        file_ver: int | None,
        script_ver: int | None,
    ) -> None:
        super().__init__(path, file_ver, script_ver)
        self.path = path
        self.file_version = file_ver
        self.script_version = script_ver


class MigrationScript(t.NamedTuple):
    version: int
    text: str

    @staticmethod
    def from_path(path: Path) -> "MigrationScript":
        file_ver = file_version(path)
        script_ver = script_version(path)

        if file_ver is None or script_ver is None or file_ver != script_ver:
            raise MigrationScriptError(path, file_ver, script_ver)
        return MigrationScript(version=file_ver, text=path.read_text())


def check_scripts(migrations: Path) -> list[MigrationScript]:
    """Check scripts in migrations directory.

    Returns list of scripts (possibly not sorted).
    Raises an exception if there's an invalid script.
    """
    return [MigrationScript.from_path(p) for p in migrations.glob("*.sql")]


def migrate(con: Connection, scripts: t.Iterable[MigrationScript]) -> None:
    """Apply migration scripts."""
    for script in sorted(scripts, key=lambda script: script.version):
        if script.version <= user_version(con):
            print("Skipping:", str(script))
        else:
            print("Applying:", str(script))
            con.executescript(script.text)


def parse_args() -> Namespace:
    parser = ArgumentParser()
    parser.add_argument("database", type=Path, help="database file")
    parser.add_argument("migrations", type=Path, help="migrations directory")
    return parser.parse_args()


def main(args: Namespace) -> None:
    try:
        scripts = check_scripts(args.migrations)
    except MigrationScriptError as exc:
        path = exc.path
        file_ver = exc.file_version
        script_ver = exc.script_version
        message = (
            "Invalid migration script version number: "
            f"{path!s} file_version={file_ver} script_version={script_ver}"
        )
        sys.exit(message)

    with connect(args.database) as con:
        migrate(con, scripts)


if __name__ == "__main__":
    main(parse_args())
