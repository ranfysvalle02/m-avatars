# m-avatars

----

# Mary Had a Little Lamb and the Secret to Software That Never Dies

**E – D – C – D – E – E – E.**

Seven notes. Even if you have never taken a music lesson in your life, you probably already hear *Mary Had a Little Lamb* playing in your head.

You can tap those seven notes on a rainbow plastic xylophone, blow them through a squeaky recorder in music class, play them on a grand piano, or beep them out of a microwave. The plastic xylophone eventually cracks. The recorder gets lost in a backpack. Computers get replaced every few years. Yet *Mary Had a Little Lamb* never breaks—because **the melody is not glued to the toy**.

Most software makes the exact opposite mistake. We glue our ideas directly to the machine.

And in an era where AI assistants can generate code in seconds, learning how to separate **the melody** from **the instrument** is the secret superpower that turns fragile weekend experiments into creations that last forever.

---

### The Glue Trap

Right now, anyone with a curious mind and an AI helper can build a working app, website, or game in an afternoon. That is an awesome superpower—until you try to take your creation somewhere new.

Suppose you build an avatar generator: a program that takes a player's username, like `"captain-pixel"`, and turns it into a unique robot face. If you build it the usual way, everything gets mashed together into one sticky ball:

* Your program grabs random numbers using a trick that only exists in Python.
* It mixes math with art, gluing `<svg>` image tags right inside the counting loop.
* It saves files to your laptop while it is still deciding where the robot's eyes go.

It works great on your computer today. Then tomorrow, you want that exact same robot face to show up inside a phone app, on a website, and on a sticker printer. Suddenly, nothing fits. Your idea is trapped inside the Python xylophone, and you have to smash it apart and rebuild it from scratch.

---

### The "Aha!" Moment: Seven Notes on Graph Paper

What if you never built with glue at all?

Imagine that instead of painting a picture directly onto a screen, you wrote your avatar down on a piece of **graph paper** using nothing but three numbers—just like writing `E – D – C – D – E – E – E` on sheet music:

* `0` means **empty space**
* `1` means **main color**
* `2` means **sparkle color**

That grid of numbers is your *Mary Had a Little Lamb*. It doesn't care about screens, browsers, file formats, or what year it is. Once you have that grid of `0`s, `1`s, and `2`s, **any instrument on Earth can play it**—a web browser, a game console, a terminal window, or a robotic LEGO printer.

Every rock-solid system in computing—from video game engines to GitHub's default profile pictures—is built on this exact split:

```text
  "captain-pixel" (Any Name)
         │
         ▼
┌─────────────────────────────────────────────┐
│  1. THE MELODY (Pure Math Sheet Music)      │
│  • Uses only simple counting & number math  │
│  • Zero files, zero screens, zero glue      │
│  • Emits a flat grid of 0s, 1s, and 2s      │
└─────────────────────────────────────────────┘
         │
         ▼  [0, 0, 1, 0, 0, 0, 1, 1, 1, 0, ...]
┌─────────────────────────────────────────────┐
│  2. THE INSTRUMENT (The Art Printer)        │
│  • Reads the 0s, 1s, and 2s                 │
│  • Draws squares, circles, or pixels!       │
└─────────────────────────────────────────────┘

```

#### Part 1: The Melody (Three Simple Math Tricks)

To turn any name into a piece of sheet music that plays identically on every computer, we use three simple rules:

1. **The Letter Blender:** We drop the letters of `"captain-pixel"` into a number-mixer that scrambles them into one starting number. The exact same name *always* makes the exact same starting number.
2. **The Infinite Spinner:** We spin that starting number like a board-game spinner over and over to get an endless stream of lucky numbers.
3. **Folded-Paper Painting:** Have you ever painted dots on the left half of a piece of paper and folded it in half so the wet paint stamps onto the right side? Our math does that with `0`s, `1`s, and `2`s. Because both sides match, random squares instantly turn into alien faces, space cruisers, or royal shields.

#### Part 2: The Instrument (The Art Printer)

Once **The Melody** finishes writing down the 25 numbers on our graph paper, **The Instrument** steps in. It is a tiny, swappable translator that reads `0`, `1`, and `2` and stamps shapes onto a real picture file (`avatar.svg`).

If a brand-new coding language is invented ten years from now, you don't touch The Melody. You just hand the sheet music to a new Instrument.

---

### Why AI Turns This Into a Superpower

For decades, programmers followed a strict rule called **DRY** (*Don't Repeat Yourself*). Copying code into two different languages was considered risky because humans make typos when translating by hand.

AI flips that rule upside down.

Because **The Melody** is just 30 lines of pure counting—with zero external libraries and zero messy file operations—an AI can translate it from Python to Go, JavaScript, or Rust in two seconds flat. Even better, you don't have to *guess* whether the AI translated it right. You check it with a **Secret Handshake**:

> **Input Word:** `"captain-pixel"`
> **Expected Notes:** `[0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 2, 1, 2, 1, 1, 2, 0, 2, 1, 0, 1, 2, 1, 0]`

Just like hearing `E – D – C – D – E – E – E` tells you immediately that someone is playing the right song, checking those 25 numbers on startup gives you **mathematical proof** that your Python server, your Go server, and a friend's phone across the world will draw the exact same avatar down to the last pixel—without ever talking to each other.

---

### Play It Yourself: Two Instruments, One Melody

Below are two complete, single-file programs—one in **Go** and one in **Python**. Neither needs any extra downloads or libraries. Tweak the settings at the top, run either file, and open `avatar.svg` in your browser to see your creation.

#### 1. Single-File Go (`identicon.go`)

Run in your terminal with `go run identicon.go`:

```go
package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// =====================================================================
// 🎨 CHANGE THESE SETTINGS (Make It Your Own Flavor!)
// =====================================================================

const (
	Username   = "captain-pixel" // Try your name, a pet, or a secret code!
	GridSize   = 5               // Try 5 (classic), 7 (detailed), or 11 (giant shield)
	DensityPct = 52              // How full the grid gets (1 to 99)
	PixelShape = "rounded"       // Try "square", "rounded", or "circle"
	OutputFile = "avatar.svg"
)

// =====================================================================
// 🎼 1. THE MELODY (Pure Math Sheet Music — Plays Identically Everywhere)
// =====================================================================

type HSL [3]int

type SheetMusic struct {
	Size    int
	Palette [3]HSL
	Grid    []uint8
}

func ComposeSong(seed string, size int, densityPct uint32) SheetMusic {
	cleanName := strings.TrimSpace(strings.ToLower(seed))

	// Trick 1: The Letter Blender (turns letters into one starting number)
	var spinner uint32 = 0x811C9DC5
	for i := 0; i < len(cleanName); i++ {
		spinner ^= uint32(cleanName[i])
		spinner *= 0x01000193
	}
	if spinner == 0 {
		spinner = 0x6D2B79F5
	}

	// Trick 2: The Infinite Spinner (gives us endless lucky numbers)
	nextNumber := func() uint32 {
		spinner ^= spinner << 13
		spinner ^= spinner >> 17
		spinner ^= spinner << 5
		return spinner
	}

	// Spin once to pick a main color from the 360-degree color wheel
	hue := int(nextNumber() % 360)
	palette := [3]HSL{
		{0, 0, 94},                  // 0 = Soft gray background
		{hue, 65, 56},               // 1 = Main character color
		{(hue + 145) % 360, 75, 60}, // 2 = Matching highlight color
	}

	// Trick 3: Folded-Paper Painting (mirror left side onto right side)
	grid := make([]uint8, size*size)
	half := (size + 1) / 2
	for y := 0; y < size; y++ {
		for x := 0; x < half; x++ {
			if nextNumber()%100 < densityPct {
				var colorNumber uint8 = 1
				if nextNumber()%100 < 28 {
					colorNumber = 2
				}
				grid[y*size+x] = colorNumber
				grid[y*size+(size-1-x)] = colorNumber
			}
		}
	}

	return SheetMusic{Size: size, Palette: palette, Grid: grid}
}

// =====================================================================
// 🎹 2. THE INSTRUMENT (Turns Sheet Music Into a Real SVG Picture)
// =====================================================================

var stamps = map[string]string{
	"square":  `<rect id="px" width="1" height="1"/>`,
	"rounded": `<rect id="px" x="0.05" y="0.05" width="0.9" height="0.9" rx="0.25"/>`,
	"circle":  `<circle id="px" cx="0.5" cy="0.5" r="0.45"/>`,
}

func PlaySVG(song SheetMusic, shape string) string {
	stamp, ok := stamps[shape]
	if !ok {
		stamp = stamps["square"]
	}

	var picture strings.Builder
	canvas := float64(song.Size) + 1.0

	fmt.Fprintf(&picture, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.1f %.1f" width="420" height="420">`+"\n", canvas, canvas)
	fmt.Fprintf(&picture, "  <defs>%s</defs>\n", stamp)
	fmt.Fprintf(&picture, `  <rect width="100%%" height="100%%" fill="hsl(%d, %d%%, %d%%)"/>`+"\n",
		song.Palette[0][0], song.Palette[0][1], song.Palette[0][2])

	for colorNumber := uint8(1); colorNumber <= 2; colorNumber++ {
		c := song.Palette[colorNumber]
		fmt.Fprintf(&picture, `  <g fill="hsl(%d, %d%%, %d%%)">`+"\n", c[0], c[1], c[2])
		for idx, cell := range song.Grid {
			if cell == colorNumber {
				x := float64(idx%song.Size) + 0.5
				y := float64(idx/song.Size) + 0.5
				fmt.Fprintf(&picture, `    <use href="#px" x="%.1f" y="%.1f"/>`+"\n", x, y)
			}
		}
		picture.WriteString("  </g>\n")
	}

	picture.WriteString("</svg>\n")
	return picture.String()
}

func checkSecretHandshake() {
	test := ComposeSong("captain-pixel", 5, 52)
	expected := []uint8{
		0, 0, 1, 0, 0,
		0, 1, 1, 1, 0,
		1, 2, 1, 2, 1,
		1, 2, 0, 2, 1,
		0, 1, 2, 1, 0,
	}
	if !reflect.DeepEqual(test.Grid, expected) {
		panic("Uh oh! The Secret Handshake notes didn't match!")
	}
}

func main() {
	checkSecretHandshake()

	song := ComposeSong(Username, GridSize, DensityPct)
	svg := PlaySVG(song, PixelShape)

	if err := os.WriteFile(OutputFile, []byte(svg), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Secret handshake passed! Saved %s for %q (%dx%d, %s)\n",
		OutputFile, Username, GridSize, GridSize, PixelShape)
}

```

---

#### 2. Single-File Python (`identicon.py`)

Run in your terminal with `python identicon.py`:

```python
from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

# =====================================================================
# 🎨 CHANGE THESE SETTINGS (Make It Your Own Flavor!)
# =====================================================================

USERNAME = "captain-pixel"      # Try your name, a pet, or a secret code!
GRID_SIZE = 5                   # Try 5 (classic), 7 (detailed), or 11 (giant shield)
DENSITY_PCT = 52                # How full the grid gets (1 to 99)
PIXEL_SHAPE = "rounded"         # Try "square", "rounded", or "circle"
OUTPUT_FILE = Path("avatar.svg")

# =====================================================================
# 🎼 1. THE MELODY (Pure Math Sheet Music — Plays Identically Everywhere)
# =====================================================================

HSL = tuple[int, int, int]


@dataclass(frozen=True)
class SheetMusic:
    size: int
    palette: tuple[HSL, HSL, HSL]
    grid: list[int]


def compose_song(seed: str, size: int = 5, density_pct: int = 52) -> SheetMusic:
    # Trick 1: The Letter Blender (turns letters into one starting number)
    spinner = 0x811C9DC5
    for byte in seed.strip().lower().encode("utf-8"):
        spinner ^= byte
        spinner = (spinner * 0x01000193) & 0xFFFFFFFF
    if spinner == 0:
        spinner = 0x6D2B79F5

    # Trick 2: The Infinite Spinner (gives us endless lucky numbers)
    def next_number() -> int:
        nonlocal spinner
        spinner ^= (spinner << 13) & 0xFFFFFFFF
        spinner ^= (spinner >> 17) & 0xFFFFFFFF
        spinner ^= (spinner << 5) & 0xFFFFFFFF
        return spinner

    # Spin once to pick a main color from the 360-degree color wheel
    hue = next_number() % 360
    palette: tuple[HSL, HSL, HSL] = (
        (0, 0, 94),                   # 0 = Soft gray background
        (hue, 65, 56),                # 1 = Main character color
        ((hue + 145) % 360, 75, 60),  # 2 = Matching highlight color
    )

    # Trick 3: Folded-Paper Painting (mirror left side onto right side)
    grid = [0] * (size * size)
    half = (size + 1) // 2
    for y in range(size):
        for x in range(half):
            if (next_number() % 100) < density_pct:
                color_number = 2 if (next_number() % 100) < 28 else 1
                grid[y * size + x] = color_number
                grid[y * size + (size - 1 - x)] = color_number

    return SheetMusic(size=size, palette=palette, grid=grid)


# =====================================================================
# 🎹 2. THE INSTRUMENT (Turns Sheet Music Into a Real SVG Picture)
# =====================================================================

STAMPS = {
    "square":  '<rect id="px" width="1" height="1"/>',
    "rounded": '<rect id="px" x="0.05" y="0.05" width="0.9" height="0.9" rx="0.25"/>',
    "circle":  '<circle id="px" cx="0.5" cy="0.5" r="0.45"/>',
}


def play_svg(song: SheetMusic, shape: str = "rounded") -> str:
    canvas = song.size + 1.0
    bg = song.palette[0]
    stamp = STAMPS.get(shape, STAMPS["square"])

    lines = [
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {canvas:.1f} {canvas:.1f}" width="420" height="420">',
        f"  <defs>{stamp}</defs>",
        f'  <rect width="100%" height="100%" fill="hsl({bg[0]}, {bg[1]}%, {bg[2]}%)"/>',
    ]

    for color_number in (1, 2):
        h, s, l = song.palette[color_number]
        lines.append(f'  <g fill="hsl({h}, {s}%, {l}%)">')
        for idx, cell in enumerate(song.grid):
            if cell == color_number:
                x = (idx % song.size) + 0.5
                y = (idx // song.size) + 0.5
                lines.append(f'    <use href="#px" x="{x:.1f}" y="{y:.1f}"/>')
        lines.append("  </g>")

    lines.append("</svg>")
    return "\n".join(lines) + "\n"


def check_secret_handshake() -> None:
    test = compose_song("captain-pixel", 5, 52)
    assert test.grid == [
        0, 0, 1, 0, 0,
        0, 1, 1, 1, 0,
        1, 2, 1, 2, 1,
        1, 2, 0, 2, 1,
        0, 1, 2, 1, 0,
    ], "Uh oh! The Secret Handshake notes didn't match!"


if __name__ == "__main__":
    check_secret_handshake()

    song = compose_song(USERNAME, GRID_SIZE, DENSITY_PCT)
    svg = play_svg(song, PIXEL_SHAPE)
    OUTPUT_FILE.write_text(svg, encoding="utf-8")

    print(f"Secret handshake passed! Saved {OUTPUT_FILE} for {USERNAME!r} ({GRID_SIZE}x{GRID_SIZE}, {PIXEL_SHAPE})")

```

---

Tools change every year. Coding languages go in and out of style. AI models get faster every month.

If you spend your life gluing ideas to today's gadgets, you will always be starting over. **Don't build a xylophone. Write down `E – D – C – D – E – E – E`.**
