# ascii

Convert images to ASCII art.

## Build

```sh
go build -o ascii ./src/
```

## Usage

```
ascii [command] [flags]
```

### Global Flags

```
--debug    enable debug logging
-h, --help help
```

---

### show

Display an image as colored ASCII art in the terminal.

```sh
ascii show <image> [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-m, --method` | `luminosity` | Brightness method: `average`, `lightness`, `luminosity` |
| `--x` | `0` (auto) | Horizontal pixel step |
| `--y` | `0` (auto) | Vertical pixel step |

**Examples:**

```sh
ascii show photo.jpg
ascii show photo.jpg -m average
ascii show photo.jpg --x 4 --y 4
```

---

### create

Convert an image to ASCII art and save as an image file.

```sh
ascii create <image> [flags]
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | `output.jpg` | Output image file path |
| `-m, --method` | `luminosity` | Brightness method: `average`, `lightness`, `luminosity` |
| `--x` | `0` (auto) | Horizontal pixel step |
| `--y` | `0` (auto) | Vertical pixel step |
| `--res` | (default cell size) | Target resolution `WxH` (e.g. `1920x1080`) |

**Examples:**

```sh
ascii create photo.jpg -o ascii_art.jpg
ascii create photo.jpg -o art.png -m luminosity --res 1920x1080
```

---

### test

```sh
go test ./src/tests/... -v
```
