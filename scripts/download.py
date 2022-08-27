"""Tatoeba data download manager."""

from argparse import ArgumentParser, Namespace
from concurrent.futures import ThreadPoolExecutor
from datetime import date, datetime
from os import remove
from pathlib import Path
import re
import sys
import typing as t

from requests import get, head
from requests.exceptions import HTTPError
from tqdm import tqdm   # type: ignore


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
    kind: Kind

    @property
    def filename(self) -> str:
        return "{}_{}.tar.bz2".format(  # noqa; pylint: disable=consider-using-f-string
            self.kind,
            self.last_modified.isoformat(),
        )

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
    def from_path(path: Path) -> t.Optional["DownloadRecord"]:
        pattern = "(links|sentences)_([0-9]+-[0-9]+-[0-9]+).tar.bz2"

        result = re.match(pattern, path.name)
        if not result:
            return None

        groups = result.groups()
        if len(groups) < 2:
            return None
        kind, date_ = groups
        return DownloadRecord(
            content_length=path.stat().st_size,
            last_modified=date.fromisoformat(date_),
            kind=t.cast(Kind, kind),
        )

    @staticmethod
    def from_headers(kind: Kind, headers: dict[str, str]) -> "DownloadRecord":
        return DownloadRecord(
            content_length=int(headers["Content-Length"]),
            last_modified=parse_date(headers["Last-Modified"]),
            kind=kind,
        )


def list_downloads(downloads: Path) -> list[DownloadRecord]:
    """List recorded files in downloads directory."""
    records = []
    for path in downloads.glob("*.tar.bz2"):
        record = DownloadRecord.from_path(path)
        if record is not None:
            records.append(record)
    return records


def fetch_headers(url: str) -> dict[str, str]:
    """Get HTTP response headers.

    May raise requests.exceptions.HTTPError.
    """
    response = head(url)
    response.raise_for_status()
    return t.cast(dict[str, str], response.headers)


def stream(url: str, chunk_size: int = 1024) -> t.Iterable[bytes]:
    """Download content in stream.

    May raise requests.exceptions.HTTPError.
    """
    response = get(url, stream=True)
    response.raise_for_status()
    for chunk in response.iter_content(chunk_size):
        yield chunk


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
        progress_bar = tqdm(total=record.content_length)
        content = bytearray()
        for chunk in stream(url):
            progress_bar.update(len(chunk))
            content.extend(chunk)
        dest.write_bytes(content)


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
