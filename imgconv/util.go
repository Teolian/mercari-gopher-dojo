package imgconv

import (
	"image"
	"image/color"
	"image/draw"
)

// toPaletted converts any image.Image to a simple paletted image for GIF.
// (For superior visual quality you'd add quantization; here we keep stdlib-only.)
func toPaletted(src image.Image) *image.Paletted {
	b := src.Bounds()
	palette := []color.Color{color.White, color.Black}
	dst := image.NewPaletted(b, palette)
	draw.FloydSteinberg.Draw(dst, b, src, b.Min)
	return dst
}
