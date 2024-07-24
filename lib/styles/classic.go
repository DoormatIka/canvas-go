package styles

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

func ComposeQuoteFrame(
	img *image.Paletted,
	font font.Face,
	resolution image.Rectangle,
	text string,
) *gg.Context {
	dc := gg.NewContext(resolution.Max.X, resolution.Max.Y);

	grad := gg.NewLinearGradient(float64(resolution.Max.X) / 2, 0, float64(resolution.Max.X), float64(resolution.Max.Y));
	grad.AddColorStop(0, color.Black);

	return dc;
}
