"""Tatoeba data download manager."""

from argparse import ArgumentParser, Namespace
from concurrent.futures import ThreadPoolExecutor
from datetime import date, datetime
from os import remove
from pathlib import Path
import sys
import typing as t

from requests import get, head
from requests.exceptions import HTTPError


LINKS = "https://downloads.tatoeba.org/exports/links.tar.bz2"
SENTENCES = "https://downloads.tatoeba.org/exports/sentences.tar.bz2"

Kind = t.Literal["links", "sentences"]


def parse_date(time: str) -> date:
    """Parse Tatoeba date."""
    layout = "%a, %d %b %Y %H:%M:%S %Z"
    return datetime.strptime(time, layout).date()


class DownloadRecord(t.NamedTuple):
    content_length: int
    last_modified: date
    etag: str
    kind: Kind

    @property
    def filename(self) -> str:
        kind = "s" if self.kind == "sentences" else "l"
        date_ = self.last_modified.isoformat()
        return f"{kind}_{date_}_{self.etag}.tar.bz2"

    def destination(self, downloads: Path) -> Path:
        """Return destination in downloads directory."""
        dest = downloads/self.filename
        if not dest.exists():
            return dest

        # Delete if inconsistent
        if dest.stat().st_size != self.content_length:
            remove(dest)
        return dest

    @staticmethod
    def from_path(path: Path) -> "DownloadRecord":
        stem = path.with_name(path.stem).stem

        kind: Kind = "links" if stem.startswith("l") else "sentences"
        date_ = len("2020-01-01")
        timestamp = stem[2:2 + date_]
        etag = stem[date_+3:]
        return DownloadRecord(
            content_length=path.stat().st_size,
            last_modified=date.fromisoformat(timestamp),
            etag=etag,
            kind=kind,
        )

    @staticmethod
    def from_headers(kind: Kind, headers: dict[str, str]) -> "DownloadRecord":
        return DownloadRecord(
            content_length=int(headers["Content-Length"]),
            last_modified=parse_date(headers["Last-Modified"]),
            etag=headers["ETag"][1:-1],
            kind=kind,
        )


def list_downloads(downloads: Path) -> list[DownloadRecord]:
    """List recorded files in downloads directory."""
    return list(map(DownloadRecord.from_path, downloads.glob("*.tar.bz2")))


def fetch_headers(url: str) -> dict[str, str]:
    """Get HTTP response headers.

    May raise requests.exceptions.HTTPError.
    """
    response = head(url)
    response.raise_for_status()
    return t.cast(dict[str, str], response.headers)


def download(url: str) -> bytes:
    """Download file.

    May raise requests.exceptions.HTTPError.
    """
    response = get(url)
    response.raise_for_status()
    return bytes(response.content)


def save_missing(url: str, downloads: Path, kind: Kind) -> None:
    """Download missing file into downloads directory."""
    assert downloads.is_dir()
    headers = fetch_headers(url)
    record = DownloadRecord.from_headers(kind, headers)
    dest = record.destination(downloads)
    if dest.exists():
        print(f"Cached: {url} ({dest!s})", file=sys.stderr)
    else:
        print(f"Downloading: {url}", file=sys.stderr)
        dest.write_bytes(download(url))


def parse_args() -> Namespace:
    parser = ArgumentParser(
        description="Tatoeba data download manager.",
    )
    parser.add_argument(
        "--ls",
        action="store_true",
        help="list downloaded data",
    )
    parser.add_argument(
        "downloads",
        type=Path,
        default=Path("build")/"tatoeba",
        nargs="?",
        help="downloads directory (default: build/tatoeba)",
    )
    return parser.parse_args()


def latest_data(downloads: Path) -> tuple[DownloadRecord, DownloadRecord]:
    """Find latest matching data in downloads directory.

    May raise Exception if no matching data is found.
    """
    links = []
    sentences = []
    for record in list_downloads(downloads):
        if record.kind == "links":
            links.append(record)
        else:
            sentences.append(record)

    def key(record: DownloadRecord) -> date:
        return record.last_modified

    links.sort(key=key)
    sentences.sort(key=key)

    while links and sentences:
        link = links[-1]
        sentence = sentences[-1]
        if link.last_modified < sentence.last_modified:
            sentences.pop()
        elif sentence.last_modified < link.last_modified:
            links.pop()
        else:
            return link, sentence
    raise Exception("no matching data found")


def main() -> None:
    args = parse_args()

    if args.ls:
        for record in list_downloads(args.downloads):
            print(record.kind, str(args.downloads/record.filename))
            print("ETag:", record.etag)
            print("Last-Modified:", record.last_modified)
            print("Content-Length:", record.content_length)
            print()
        return

    args.downloads.mkdir(parents=True, exist_ok=True)
    try:
        with ThreadPoolExecutor(max_workers=2) as executor:
            futures = [
                executor.submit(save_missing, LINKS, args.downloads, "links"),
                executor.submit(
                    save_missing,
                    SENTENCES,
                    args.downloads,
                    "sentences",
                ),
            ]
            for future in futures:
                future.result()
    except HTTPError:
        print("download failed", file=sys.stderr)


if __name__ == "__main__":
    main()
