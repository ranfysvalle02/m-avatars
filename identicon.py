#!/usr/bin/env python3
"""m-avatars: one melody, many instruments. Python performance of SPEC.md.

    python3 identicon.py captain-pixel
    python3 identicon.py --format svg captain-pixel > avatar.svg
    python3 identicon.py --check golden.json
"""

from __future__ import annotations

import argparse
import json
import sys
from dataclasses import dataclass

# =====================================================================
# LANE 1: THE MELODY
# Pure: data in, data out. No I/O, no clock, no platform opinions.
# =====================================================================

HSL = tuple[int, int, int]
MASK32 = 0xFFFFFFFF


@dataclass(frozen=True)
class Sheet:
    size: int
    palette: tuple[HSL, HSL, HSL]
    grid: tuple[int, ...]


def compose(seed: str, size: int, density: int) -> Sheet:
    # bytes.strip()/.lower() are ASCII-only by definition, which is what SPEC.md requires.
    state = 0x811C9DC5
    for byte in seed.encode("utf-8").strip().lower():
        state = ((state ^ byte) * 0x01000193) & MASK32
    if state == 0:
        state = 0x6D2B79F5

    def next_u32() -> int:
        nonlocal state
        state ^= (state << 13) & MASK32
        state ^= state >> 17
        state ^= (state << 5) & MASK32
        return state

    hue = next_u32() % 360
    palette = ((0, 0, 94), (hue, 65, 56), ((hue + 145) % 360, 75, 60))

    grid = [0] * (size * size)
    half = (size + 1) // 2
    for y in range(size):
        for x in range(half):
            if next_u32() % 100 < density:
                color = 2 if next_u32() % 100 < 28 else 1
                grid[y * size + x] = color
                grid[y * size + (size - 1 - x)] = color

    return Sheet(size, palette, tuple(grid))


# =====================================================================
# LANE 2: THE INSTRUMENTS
# Each one reads a Sheet and makes no musical decisions of its own.
# =====================================================================


def rows(sheet: Sheet) -> list[str]:
    n = sheet.size
    return ["".join(map(str, sheet.grid[i : i + n])) for i in range(0, n * n, n)]


def render_json(sheet: Sheet) -> str:
    return json.dumps({"size": sheet.size, "palette": sheet.palette, "grid": rows(sheet)}, separators=(",", ":"))


def render_text(sheet: Sheet) -> str:
    glyphs = {"0": "  ", "1": "██", "2": "▒▒"}
    return "\n".join("".join(glyphs[c] for c in row) for row in rows(sheet))


STAMPS = {
    "square": '<rect id="px" width="1" height="1"/>',
    "rounded": '<rect id="px" x="0.05" y="0.05" width="0.9" height="0.9" rx="0.25"/>',
    "circle": '<circle id="px" cx="0.5" cy="0.5" r="0.45"/>',
}


def render_svg(sheet: Sheet, shape: str = "rounded") -> str:
    def hsl(c: HSL) -> str:
        return f"hsl({c[0]}, {c[1]}%, {c[2]}%)"

    view = sheet.size + 1
    out = [
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {view} {view}" width="420" height="420">',
        f"  <defs>{STAMPS[shape]}</defs>",
        f'  <rect width="100%" height="100%" fill="{hsl(sheet.palette[0])}"/>',
    ]
    for color in (1, 2):
        out.append(f'  <g fill="{hsl(sheet.palette[color])}">')
        for i, cell in enumerate(sheet.grid):
            if cell == color:
                x, y = i % sheet.size, i // sheet.size
                out.append(f'    <use href="#px" x="{x + 0.5}" y="{y + 0.5}"/>')
        out.append("  </g>")
    out.append("</svg>")
    return "\n".join(out)


# =====================================================================
# LANE 3: THE STAGE
# The only lane that touches the outside world: argv, files, stdout.
# =====================================================================

HANDSHAKE = ("captain-pixel", 5, 52, ["00100", "01110", "12121", "12021", "01210"])


def check(path: str) -> int:
    with open(path, encoding="utf-8") as f:
        golden = json.load(f)
    failures = 0
    for v in golden["vectors"]:
        sheet = compose(v["seed"], v["size"], v["density"])
        ok = rows(sheet) == v["grid"] and [list(c) for c in sheet.palette] == v["palette"]
        failures += not ok
        print(f"{'PASS' if ok else 'FAIL'}  {v['seed']!r} size={v['size']} density={v['density']}")
    print(f"python: {len(golden['vectors']) - failures}/{len(golden['vectors'])} vectors match {golden['spec']}")
    return 1 if failures else 0


def main(argv: list[str]) -> int:
    p = argparse.ArgumentParser(description="Deterministic identicons. See SPEC.md.")
    p.add_argument("seed", nargs="?", help="any name, e.g. captain-pixel")
    p.add_argument("--size", type=int, default=5, help="grid size, 1-64 (default 5)")
    p.add_argument("--density", type=int, default=52, help="fill percentage, 0-100 (default 52)")
    p.add_argument("--format", choices=["text", "svg", "json"], default="text")
    p.add_argument("--shape", choices=sorted(STAMPS), default="rounded", help="SVG pixel shape")
    p.add_argument("--check", metavar="GOLDEN", help="verify against golden vectors and exit")
    args = p.parse_args(argv)

    seed, size, density, expected = HANDSHAKE
    if rows(compose(seed, size, density)) != expected:
        print("handshake failed: this build does not play the m-avatars melody", file=sys.stderr)
        return 2

    if args.check:
        return check(args.check)
    if args.seed is None:
        p.error("a seed is required (or use --check)")
    if not 1 <= args.size <= 64:
        p.error("--size must be between 1 and 64")
    if not 0 <= args.density <= 100:
        p.error("--density must be between 0 and 100")

    sheet = compose(args.seed, args.size, args.density)
    if args.format == "svg":
        print(render_svg(sheet, args.shape))
    elif args.format == "json":
        print(render_json(sheet))
    else:
        print(render_text(sheet))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
