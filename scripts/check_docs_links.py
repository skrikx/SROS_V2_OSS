#!/usr/bin/env python3
from __future__ import annotations

import re
import sys
from pathlib import Path
from urllib.parse import unquote


LINK_PATTERN = re.compile(r"!?\[[^\]]*\]\(([^)\s]+)(?:\s+\"[^\"]*\")?\)")


def should_skip(link: str) -> bool:
    return link.startswith(("http://", "https://", "mailto:", "tel:", "data:", "#"))


def resolve_target(repo_root: Path, source_file: Path, link: str) -> Path | None:
    clean = link.strip().strip("<>")
    if should_skip(clean):
        return None

    clean = clean.split("#", 1)[0].split("?", 1)[0]
    clean = unquote(clean)
    if not clean:
        return None

    if clean.startswith("/"):
        return repo_root / clean.lstrip("/")
    return (source_file.parent / clean).resolve()


def markdown_files(repo_root: Path) -> list[Path]:
    files = [repo_root / "README.md", repo_root / "CONTRIBUTING.md"]
    docs_dir = repo_root / "docs"
    files.extend(sorted(docs_dir.rglob("*.md")))
    return [path for path in files if path.exists()]


def main() -> int:
    repo_root = Path(__file__).resolve().parents[1]
    failures: list[str] = []

    for source in markdown_files(repo_root):
        rel_source = source.relative_to(repo_root).as_posix()
        for line_no, line in enumerate(source.read_text(encoding="utf-8").splitlines(), start=1):
            for match in LINK_PATTERN.finditer(line):
                raw_link = match.group(1)
                target = resolve_target(repo_root, source, raw_link)
                if target is None:
                    continue
                if not target.exists():
                    failures.append(
                        f"{rel_source}:{line_no}: missing link target '{raw_link}' -> '{target.relative_to(repo_root).as_posix()}'"
                    )

    if failures:
        print("docs_link_check_failed")
        for failure in failures:
            print(failure)
        return 1

    print("docs_link_check_passed")
    return 0


if __name__ == "__main__":
    sys.exit(main())
