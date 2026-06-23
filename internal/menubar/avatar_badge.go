package menubar

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"

	"github.com/haritabh17/theirtime/internal/slack"
)

// avatarEntry holds display-ready avatar bytes and the original avatar size,
// used so the menu bar scales the face to the configured display size.
type avatarEntry struct {
	data        []byte
	contentSize int
}

var (
	colorSlackGreen = color.RGBA{R: 0x2B, G: 0xAC, B: 0x76, A: 0xFF}
	colorAwayBar    = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

func applyPresenceBadge(avatarPNG []byte, presence slack.Presence) ([]byte, error) {
	switch presence {
	case slack.PresenceActive, slack.PresenceAway:
	default:
		return avatarPNG, nil
	}

	src, _, err := image.Decode(bytes.NewReader(avatarPNG))
	if err != nil {
		return nil, err
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return avatarPNG, nil
	}

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(canvas, canvas.Bounds(), src, bounds.Min, draw.Src)
	drawPresenceBar(canvas, presence)

	var out bytes.Buffer
	if err := png.Encode(&out, canvas); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func drawPresenceBar(img draw.Image, presence slack.Presence) {
	bounds := img.Bounds()
	thickness := presenceBarThickness(bounds.Dx())
	if thickness == 0 {
		return
	}
	radius := thickness / 2
	if radius == 0 {
		radius = 1
	}
	barRight := bounds.Min.X + thickness
	centerX := bounds.Min.X + radius
	topCenterY := bounds.Min.Y + radius
	bottomCenterY := bounds.Max.Y - radius - 1
	barColor := presenceBarColor(presence)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < barRight; x++ {
			if insidePresenceBar(x, y, radius, centerX, topCenterY, bottomCenterY) {
				img.Set(x, y, barColor)
			}
		}
	}
}

func insidePresenceBar(x, y int, radius, centerX, topCenterY, bottomCenterY int) bool {
	if y >= topCenterY && y <= bottomCenterY {
		return true
	}
	cy := topCenterY
	if y > bottomCenterY {
		cy = bottomCenterY
	}
	dx := x - centerX
	dy := y - cy
	return dx*dx+dy*dy <= radius*radius
}

func presenceBarThickness(width int) int {
	if width <= 0 {
		return 0
	}
	thickness := width / 8
	if thickness < 3 {
		thickness = 3
	}
	if thickness > width {
		return width
	}
	return thickness
}

func presenceBarColor(presence slack.Presence) color.Color {
	if presence == slack.PresenceActive {
		return colorSlackGreen
	}
	return colorAwayBar
}

func imageSquareSize(raw []byte) int {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil || cfg.Width <= 0 {
		return 0
	}
	return cfg.Width
}

func displayAvatar(showPresence bool, raw []byte, presence slack.Presence) []byte {
	if !showPresence || len(raw) == 0 {
		return raw
	}
	badged, err := applyPresenceBadge(raw, presence)
	if err != nil {
		return raw
	}
	return badged
}
