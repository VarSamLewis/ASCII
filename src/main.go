package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	fatih "github.com/fatih/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/term"
)

func parseRes(s string) (int, int, bool) {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return 0, 0, false
	}
	w, err1 := strconv.Atoi(parts[0])
	h, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || w <= 0 || h <= 0 {
		return 0, 0, false
	}
	return w, h, true
}

func main() {
	method := flag.String("method", "luminosity", "brightness method: average, lightness, luminosity")
	xPixel := flag.Int("x", 0, "horizontal pixel step (0 = auto)")
	yPixel := flag.Int("y", 0, "vertical pixel step (0 = auto)")
	output := flag.String("o", "", "output JPG file path")
	res := flag.String("res", "", "target output resolution WxH (e.g. 1920x1080)")

	flag.Parse()
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("Width: %d Pixels\n", width)
	fmt.Printf("Height: %d Pixels\n", height)

	fd := int(os.Stdout.Fd())
	terminalWidth, terminalHeight, err := term.GetSize(fd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("terminalWidth: %d \n", terminalWidth)
	fmt.Printf("terminalHeight: %d \n", terminalHeight)

	var xStep, yStep int
	if *xPixel > 0 {
		xStep = *xPixel
	} else {
		xStep = max(1, (width+terminalWidth-1)/terminalWidth)
	}
	if *yPixel > 0 {
		yStep = *yPixel
	} else {
		yStep = max(1, (height+terminalHeight/2-1)/(terminalHeight/2))
	}

	fmt.Printf("xStep: %d \n", xStep)
	fmt.Printf("yStep: %d \n", yStep)

	ramp := "`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

	type cell struct {
		ch string
		r  uint8
		g  uint8
		b  uint8
	}
	var grid [][]cell

	for y := 0; y < height; y += yStep {
		var row []cell
		for x := 0; x < width; x += xStep {
			r, g, b, _ := img.At(
				x+bounds.Min.X,
				y+bounds.Min.Y,
			).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			brightness := calc_brightness(r8, g8, b8, *method)
			index := brightness * (len(ramp) - 1) / 255

			char := string(ramp[index])

			fatih.RGB(int(r8), int(g8), int(b8)).Print(char)
			row = append(row, cell{ch: char, r: r8, g: g8, b: b8})
		}
		grid = append(grid, row)
		fmt.Println()
	}

	if *output != "" {
		imgW := 0
		for _, row := range grid {
			if len(row) > imgW {
				imgW = len(row)
			}
		}
		imgH := len(grid)

		cw, ch := 6, 12
		if *res != "" {
			resW, resH, ok := parseRes(*res)
			if !ok {
				fmt.Fprintf(os.Stderr, "invalid resolution: %s (use WxH)\n", *res)
				os.Exit(1)
			}
			cw = resW / imgW
			ch = resH / imgH
			if cw < 1 {
				cw = 1
			}
			if ch < 1 {
				ch = 1
			}
		}

		pg := image.NewRGBA(image.Rect(0, 0, imgW*cw, imgH*ch))
		d := font.Drawer{
			Dst:  pg,
			Src:  image.White,
			Face: basicfont.Face7x13,
			Dot:  fixed.P(0, 0),
		}
		for y, row := range grid {
			for x, c := range row {
				d.Src = image.NewUniform(color.RGBA{R: c.r, G: c.g, B: c.b, A: 255})
				d.Dot = fixed.P(x*cw, (y+1)*ch)
				d.DrawString(c.ch)
			}
		}
		f, err := os.Create(*output)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		jpeg.Encode(f, pg, &jpeg.Options{Quality: 95})
		fmt.Printf("Saved to %s (%dx%d)\n", *output, imgW*cw, imgH*ch)
	}
}

func calc_brightness(r uint8, g uint8, b uint8, method string) int {
	switch method {
	case "average":
		return int((uint32(r) + uint32(g) + uint32(b)) / 3)
	case "lightness":
		return int((max(r, g, b) + min(r, g, b)) / 2)
	case "luminosity":
		return int(0.21*float64(r) + 0.72*float64(g) + 0.07*float64(b))
	default:
		fmt.Fprintf(os.Stderr, "unknown method: %s\n", method)
		os.Exit(1)
		return 0
	}
}
