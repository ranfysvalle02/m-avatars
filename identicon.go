// m-avatars: one melody, many instruments. Go performance of SPEC.md.
//
//	go run identicon.go captain-pixel
//	go run identicon.go -format svg captain-pixel > avatar.svg
//	go run identicon.go -check golden.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// =====================================================================
// LANE 1: THE MELODY
// Pure: data in, data out. No I/O, no clock, no platform opinions.
// =====================================================================

type HSL [3]int

type Sheet struct {
	Size    int
	Palette [3]HSL
	Grid    []uint8
}

func Compose(seed string, size int, density uint32) Sheet {
	// ASCII-only trim and case fold, per SPEC.md. strings.ToLower/TrimSpace are Unicode-aware and would diverge.
	normalized := strings.Trim(seed, " \t\n\v\f\r")

	var state uint32 = 0x811C9DC5
	for i := 0; i < len(normalized); i++ {
		b := normalized[i]
		if 'A' <= b && b <= 'Z' {
			b += 'a' - 'A'
		}
		state ^= uint32(b)
		state *= 0x01000193
	}
	if state == 0 {
		state = 0x6D2B79F5
	}

	next := func() uint32 {
		state ^= state << 13
		state ^= state >> 17
		state ^= state << 5
		return state
	}

	hue := int(next() % 360)
	palette := [3]HSL{{0, 0, 94}, {hue, 65, 56}, {(hue + 145) % 360, 75, 60}}

	grid := make([]uint8, size*size)
	half := (size + 1) / 2
	for y := 0; y < size; y++ {
		for x := 0; x < half; x++ {
			if next()%100 < density {
				var color uint8 = 1
				if next()%100 < 28 {
					color = 2
				}
				grid[y*size+x] = color
				grid[y*size+(size-1-x)] = color
			}
		}
	}

	return Sheet{Size: size, Palette: palette, Grid: grid}
}

// =====================================================================
// LANE 2: THE INSTRUMENTS
// Each one reads a Sheet and makes no musical decisions of its own.
// =====================================================================

func Rows(s Sheet) []string {
	rows := make([]string, s.Size)
	for y := range rows {
		var b strings.Builder
		for _, cell := range s.Grid[y*s.Size : (y+1)*s.Size] {
			b.WriteByte('0' + cell)
		}
		rows[y] = b.String()
	}
	return rows
}

func RenderJSON(s Sheet) string {
	out, _ := json.Marshal(struct {
		Size    int      `json:"size"`
		Palette [3]HSL   `json:"palette"`
		Grid    []string `json:"grid"`
	}{s.Size, s.Palette, Rows(s)})
	return string(out)
}

func RenderText(s Sheet) string {
	glyphs := strings.NewReplacer("0", "  ", "1", "██", "2", "▒▒")
	return glyphs.Replace(strings.Join(Rows(s), "\n"))
}

var stamps = map[string]string{
	"square":  `<rect id="px" width="1" height="1"/>`,
	"rounded": `<rect id="px" x="0.05" y="0.05" width="0.9" height="0.9" rx="0.25"/>`,
	"circle":  `<circle id="px" cx="0.5" cy="0.5" r="0.45"/>`,
}

func RenderSVG(s Sheet, shape string) string {
	hsl := func(c HSL) string { return fmt.Sprintf("hsl(%d, %d%%, %d%%)", c[0], c[1], c[2]) }
	view := s.Size + 1

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="420" height="420">`+"\n", view, view)
	fmt.Fprintf(&b, "  <defs>%s</defs>\n", stamps[shape])
	fmt.Fprintf(&b, `  <rect width="100%%" height="100%%" fill="%s"/>`+"\n", hsl(s.Palette[0]))
	for color := uint8(1); color <= 2; color++ {
		fmt.Fprintf(&b, `  <g fill="%s">`+"\n", hsl(s.Palette[color]))
		for i, cell := range s.Grid {
			if cell == color {
				x, y := i%s.Size, i/s.Size
				fmt.Fprintf(&b, `    <use href="#px" x="%.1f" y="%.1f"/>`+"\n", float64(x)+0.5, float64(y)+0.5)
			}
		}
		b.WriteString("  </g>\n")
	}
	b.WriteString("</svg>")
	return b.String()
}

// =====================================================================
// LANE 3: THE STAGE
// The only lane that touches the outside world: argv, files, stdout.
// =====================================================================

var handshake = []string{"00100", "01110", "12121", "12021", "01210"}

type golden struct {
	Spec    string `json:"spec"`
	Vectors []struct {
		Seed    string   `json:"seed"`
		Size    int      `json:"size"`
		Density uint32   `json:"density"`
		Palette [3]HSL   `json:"palette"`
		Grid    []string `json:"grid"`
	} `json:"vectors"`
}

func check(path string) int {
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	var g golden
	if err := json.Unmarshal(raw, &g); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	failures := 0
	for _, v := range g.Vectors {
		s := Compose(v.Seed, v.Size, v.Density)
		status := "PASS"
		if !reflect.DeepEqual(Rows(s), v.Grid) || s.Palette != v.Palette {
			status = "FAIL"
			failures++
		}
		fmt.Printf("%s  %q size=%d density=%d\n", status, v.Seed, v.Size, v.Density)
	}
	fmt.Printf("go: %d/%d vectors match %s\n", len(g.Vectors)-failures, len(g.Vectors), g.Spec)
	if failures > 0 {
		return 1
	}
	return 0
}

func main() {
	size := flag.Int("size", 5, "grid size, 1-64")
	density := flag.Int("density", 52, "fill percentage, 0-100")
	format := flag.String("format", "text", "text, svg, or json")
	shape := flag.String("shape", "rounded", "SVG pixel shape: square, rounded, or circle")
	golden := flag.String("check", "", "verify against a golden vectors file and exit")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: identicon [flags] SEED\n       identicon -check golden.json")
		flag.PrintDefaults()
	}
	flag.Parse()

	if !reflect.DeepEqual(Rows(Compose("captain-pixel", 5, 52)), handshake) {
		fmt.Fprintln(os.Stderr, "handshake failed: this build does not play the m-avatars melody")
		os.Exit(2)
	}

	if *golden != "" {
		os.Exit(check(*golden))
	}

	fail := func(msg string) {
		fmt.Fprintln(os.Stderr, "error:", msg)
		flag.Usage()
		os.Exit(2)
	}
	switch {
	case flag.NArg() != 1:
		fail("exactly one SEED is required (put flags before it)")
	case *size < 1 || *size > 64:
		fail("-size must be between 1 and 64")
	case *density < 0 || *density > 100:
		fail("-density must be between 0 and 100")
	case *format != "text" && *format != "svg" && *format != "json":
		fail("-format must be text, svg, or json")
	case stamps[*shape] == "":
		fail("-shape must be square, rounded, or circle")
	}

	sheet := Compose(flag.Arg(0), *size, uint32(*density))
	switch *format {
	case "svg":
		fmt.Println(RenderSVG(sheet, *shape))
	case "json":
		fmt.Println(RenderJSON(sheet))
	default:
		fmt.Println(RenderText(sheet))
	}
}
