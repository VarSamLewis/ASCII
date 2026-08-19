package cmd

import (
	"fmt"
	"image"
	"log/slog"
	"os"

	ascii "ascii/ascii/src/internal/ascii"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var createCmd = &cobra.Command{
	Use:   "create [image]",
	Short: "Convert an image to ASCII art and save as an image file",
	Long:  "Takes an input image (jpg/png/gif) and outputs an ASCII art version as an image file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method, _ := cmd.Flags().GetString("method")
		xPixel, _ := cmd.Flags().GetInt("x")
		yPixel, _ := cmd.Flags().GetInt("y")
		output, _ := cmd.Flags().GetString("output")
		res, _ := cmd.Flags().GetString("res")

		slog.Debug("create command", "input", args[0], "output", output, "method", method)

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

		if err := ascii.RenderImage(grid, output, res); err != nil {
			return fmt.Errorf("rendering image: %w", err)
		}

		imgW, imgH := ascii.GridDimensions(grid)
		fmt.Printf("Saved to %s\n", output)
		slog.Debug("output info", "grid_cols", imgW, "grid_rows", imgH)

		return nil
	},
}

func init() {
	createCmd.Flags().StringP("method", "m", "luminosity", "brightness method: average, lightness, luminosity")
	createCmd.Flags().Int("x", 0, "horizontal pixel step (0 = auto)")
	createCmd.Flags().Int("y", 0, "vertical pixel step (0 = auto)")
	createCmd.Flags().StringP("output", "o", "output.jpg", "output image file path")
	createCmd.Flags().String("res", "", "target output resolution WxH (e.g. 1920x1080)")

	rootCmd.AddCommand(createCmd)
}
