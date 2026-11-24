# m-avatars spec v1: the sheet music

This document is the program. Every file named `identicon.*` is just one
performance of it.

A conforming implementation turns `(seed, size, density)` into a **sheet**:
a palette of three HSL colors and a `size × size` grid of the digits `0`, `1`,
and `2`. Two conforming implementations must produce identical sheets for
every input. [`golden.json`](golden.json) is the test that decides whether
they do.

Everything after the sheet (pixels, SVG, terminal blocks, LEGO bricks) is
out of scope. Renderers can do whatever they like.

---

## Inputs

| Name      | Type          | Domain                            |
|-----------|---------------|-----------------------------------|
| `seed`    | text          | any Unicode string                |
| `size`    | integer       | `1 ≤ size ≤ 64`                   |
| `density` | integer       | `0 ≤ density ≤ 100` (a percentage)|

The core can assume its inputs are in range. Checking them is the host's job.

All arithmetic below is on **unsigned 32-bit integers**. Every `+`, `*`, `<<`
wraps modulo 2³². `>>` is a logical (zero-fill) shift.

## Step 1: Normalize (bytes, not characters)

1. Encode `seed` as UTF-8 bytes.
2. Strip leading and trailing bytes in the set
   `{0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x20}` (ASCII whitespace).
3. Map each byte in `0x41–0x5A` (`A–Z`) to itself `+ 0x20` (`a–z`).
   Leave every other byte untouched.

Do **not** use your language's `lower()`, `toLowerCase()`, `strip()`,
`trim()`, or Unicode normalization. They disagree with each other on inputs
like `İstanbul`, `ΣΑΣ`, and `"\uFEFFbom"`, and that disagreement would leak
into the avatar. If a host wants NFC or locale-aware case folding, it does
that *before* calling the core.

## Step 2: Hash (FNV-1a, 32-bit)

```
h = 0x811C9DC5
for each byte b:
    h = h XOR b
    h = h * 0x01000193
if h == 0:
    h = 0x6D2B79F5
```

## Step 3: Random stream (xorshift32)

`h` becomes the generator state `s`. Each call to `next()` does:

```
s = s XOR (s << 13)
s = s XOR (s >> 17)
s = s XOR (s << 5)
return s
```

Calls to `next()` happen in **exactly** the order below. One extra or missing
draw shifts every later decision.

## Step 4: Palette

```
hue     = next() mod 360
palette = [ (0,   0,  94),                 # 0: background
            (hue, 65, 56),                 # 1: primary
            ((hue + 145) mod 360, 75, 60) ] # 2: accent
```

Each entry is `(hue°, saturation%, lightness%)`.

## Step 5: Grid (mirrored)

```
grid = [0] * (size * size)                 # row-major
half = (size + 1) / 2                      # integer division
for y in 0 .. size-1:
    for x in 0 .. half-1:
        if next() mod 100 < density:
            color = 2 if next() mod 100 < 28 else 1   # 2nd draw ONLY on a hit
            grid[y*size + x]              = color
            grid[y*size + (size - 1 - x)] = color
```

## Output

```
sheet = { size, palette, grid }
```

The canonical text form used by `golden.json` writes each grid row as a
string of digits, for example `["00100", "01110", "12121", "12021", "01210"]`.

## Conformance

An implementation conforms to v1 if it reproduces every vector in
[`golden.json`](golden.json). Any change that alters a vector's output is a
new spec version. Bump `spec`, keep the old vectors, and don't call it a
refactor.
