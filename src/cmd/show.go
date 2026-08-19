package cmd

import (
	"fmt"
	"image"
	"log/slog"
	"os"

	ascii "ascii/ascii/src/internal/ascii"

	fatih "github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var showCmd = &cobra.Command{
	Use:   "show [image]",
	Short: "Display an image as ASCII art in the terminal",
	Long:  "Takes an input image (jpg/png/gif) and prints the ASCII art version to the terminal with colors.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")
		xPixel, _ := cmd.Flags().GetInt("x")
		yPixel, _ := cmd.Flags().GetInt("y")

		slog.Debug("show command", "input", args[0], "method", method)

		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("opening input file: %w", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("decoding image: %w", err)
		}

		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		slog.Debug("image decoded", "width", width, "height", height)

		fd := int(os.Stdout.Fd())
		terminalWidth, terminalHeight, err := term.GetSize(fd)
		if err != nil {
			slog.Debug("could not get terminal size, using defaults", "error", err)
			terminalWidth = 80
			terminalHeight = 24
		}

		slog.Debug("terminal size", "width", terminalWidth, "height", terminalHeight)

		xStep, yStep := xPixel, yPixel
		if xStep <= 0 || yStep <= 0 {
			autoX, autoY := ascii.AutoSteps(width, height, terminalWidth, terminalHeight)
			if xStep <= 0 {
				xStep = autoX
			}
			if yStep <= 0 {
				yStep = autoY
			}
		}

		slog.Debug("step sizes", "xStep", xStep, "yStep", yStep)

		grid := ascii.GenerateGrid(img, ascii.ConvertOptions{
			Method: method,
			XStep:  xStep,
			YStep:  yStep,
		})

		for _, row := range grid {
			for _, c := range row {
				fatih.RGB(int(c.R), int(c.G), int(c.B)).Print(c.Ch)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	showCmd.Flags().StringP("method", "m", "luminosity", "brightness method: average, lightness, luminosity")
	showCmd.Flags().Int("x", 0, "horizontal pixel step (0 = auto)")
	showCmd.Flags().Int("y", 0, "vertical pixel step (0 = auto)")

	rootCmd.AddCommand(showCmd)
}
