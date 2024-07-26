package styles

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// for use in images.
func ModifyQuoteImage(src image.Image, font *font.Face) *gg.Context {
	width, height := src.Bounds().Max.X, src.Bounds().Max.Y;

	// average_luminosity, _ := utils.GetAverageBrightnessOfRGBA(src, width, height);
	screenResolution := image.Rect(0, 0, width, height);

	return ComposeQuoteFrameImage(&src, *font, screenResolution, "You're amazing at what you do.");
}

func ComposeQuoteFrameImage(
	img *image.Image,
	font font.Face,
	resolution image.Rectangle,
	text string,
) *gg.Context {
	dc := gg.NewContext(resolution.Max.X, resolution.Max.Y);

	grad := gg.NewLinearGradient(float64(resolution.Max.X) / 2, 0, float64(resolution.Max.X), float64(resolution.Max.Y));
	grad.AddColorStop(0, color.White);
	grad.AddColorStop(0.5, color.Black);
	
	dc.SetColor(color.White);
	dc.SetFillStyle(grad);
	dc.MoveTo(0, 0);
	dc.LineTo(0, float64(dc.Height()));
	dc.LineTo(float64(dc.Width()), float64(dc.Height()));
	dc.LineTo(float64(dc.Width()), 0);
	dc.ClosePath();
	dc.Fill();

	return dc;
}
