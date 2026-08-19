package main

import (
	"flag"
	"fmt"
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
		row := make([]rune, 0, width/xStep)

		for x := 0; x < width; x += xStep {
			r, g, b, _ := img.At(
				x+bounds.Min.X,
				y+bounds.Min.Y,
			).RGBA()

			brightness := calc_brightness(r, g, b, *method)
			index := brightness * (len(ramp) - 1) / 255
			row = append(row, rune(ramp[index]))
		}

		fmt.Println(string(row))
	}

}

func calc_brightness(r uint32, g uint32, b uint32, method string) int {
	r8, g8, b8 := r>>8, g>>8, b>>8
	switch method {
	case "average":
		return int((r8 + g8 + b8) / 3)
	case "lightness":
		return int((max(r8, g8, b8) + min(r8, g8, b8)) / 2)
	case "luminosity":
		return int(0.21*float64(r8) + 0.72*float64(g8) + 0.07*float64(b8))
	default:
		fmt.Fprintf(os.Stderr, "unknown method: %s\n", method)
		os.Exit(1)
		return 0
	}
}
