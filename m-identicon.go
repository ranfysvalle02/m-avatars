package main

import (
	"fmt"
	"reflect"
	"strings"
)

type HSL [3]int

type Sprite struct {
	Size    int
	Palette [3]HSL
	Grid    []uint8
}

// =====================================================================
// 1. PURE MATH CORE (Zero external dependencies)
// =====================================================================

func GenerateCore(seed string, size int, densityPct uint32) Sprite {
	normalized := strings.TrimSpace(strings.ToLower(seed))

	// FNV-1a 32-bit hash over UTF-8 bytes
	var state uint32 = 0x811C9DC5
	for i := 0; i < len(normalized); i++ {
		state ^= uint32(normalized[i])
		state *= 0x01000193
	}
	if state == 0 {
		state = 0x6D2B79F5
	}

	// Xorshift32 PRNG
	nextU32 := func() uint32 {
		state ^= state << 13
		state ^= state >> 17
		state ^= state << 5
		return state
	}

	// Deterministic HSL palette
	hue := int(nextU32() % 360)
	palette := [3]HSL{
		{0, 0, 94},
		{hue, 65, 56},
		{(hue + 145) % 360, 75, 60},
	}

	// Left-right symmetric orbit projection
	grid := make([]uint8, size*size)
	half := (size + 1) / 2
	for y := 0; y < size; y++ {
		for x := 0; x < half; x++ {
			if nextU32()%100 < densityPct {
				var colorIdx uint8 = 1
				if nextU32()%100 < 28 {
					colorIdx = 2
				}
				grid[y*size+x] = colorIdx
				grid[y*size+(size-1-x)] = colorIdx
			}
		}
	}

	return Sprite{Size: size, Palette: palette, Grid: grid}
}

// =====================================================================
// 2. HOST SHELL (SVG Renderer + Golden Test Verification)
// =====================================================================

func RenderSVG(s Sprite) string {
	var b strings.Builder
	canvas := float64(s.Size) + 1.0

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.1f %.1f" width="420" height="420">`+"\n", canvas, canvas)
	b.WriteString(`  <defs><rect id="px" x="0.05" y="0.05" width="0.9" height="0.9" rx="0.25"/></defs>` + "\n")
	fmt.Fprintf(&b, `  <rect width="100%%" height="100%%" fill="hsl(%d, %d%%, %d%%)"/>`+"\n",
		s.Palette[0][0], s.Palette[0][1], s.Palette[0][2])

	for colorIdx := uint8(1); colorIdx <= 2; colorIdx++ {
		c := s.Palette[colorIdx]
		fmt.Fprintf(&b, `  <g fill="hsl(%d, %d%%, %d%%)">`+"\n", c[0], c[1], c[2])
		for idx, cell := range s.Grid {
			if cell == colorIdx {
				x := float64(idx%s.Size) + 0.5
				y := float64(idx/s.Size) + 0.5
				fmt.Fprintf(&b, `    <use href="#px" x="%.1f" y="%.1f"/>`+"\n", x, y)
			}
		}
		b.WriteString("  </g>\n")
	}

	b.WriteString("</svg>")
	return b.String()
}

func main() {
	sprite := GenerateCore("captain-pixel", 5, 52)
	expectedGrid := []uint8{
		0, 0, 1, 0, 0,
		0, 1, 1, 1, 0,
		1, 2, 1, 2, 1,
		1, 2, 0, 2, 1,
		0, 1, 2, 1, 0,
	}
	if !reflect.DeepEqual(sprite.Grid, expectedGrid) {
		panic("golden vector mismatch")
	}
	fmt.Println("Go parity verified! Palette:", sprite.Palette)
}
