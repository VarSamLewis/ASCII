package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/term"
	"image"
	_ "image/jpeg"
	"os"
)

func main() {
	method := flag.String("method", "luminosity", "brightness method: average, lightness, luminosity")
	xPixel := flag.Int("x", 0, "horizontal pixel step (0 = auto)")
	yPixel := flag.Int("y", 0, "vertical pixel step (0 = auto)")

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

	for y := 0; y < height; y += yStep {
		// row := make([]rune, 0, width/xStep)

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
			//row = append(row, rune(ramp[index]))

			char := string(ramp[index])

			color.RGB(
				int(r8),
				int(g8),
				int(b8),
			).Print(char)
		}

		//fmt.Println(string(row))
		fmt.Println()
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
