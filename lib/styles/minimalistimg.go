package styles

import (
	"fmt"
	"image"

	"golang.org/x/image/font"
	"github.com/fogleman/gg"

	"canvas/lib/utils"
)

func ModifyMinimalistImage(src *image.Image, font *font.Face, text string) (*image.Image, error) {
	switch im := (*src).(type) {
	case *image.RGBA:
		return modifyMinimalistRGBA(im, font, text), nil;
	default:
		return nil, fmt.Errorf("Image is not of type RGBA.");
	}
}

// for use in images.
func modifyMinimalistRGBA(src *image.RGBA, font *font.Face, text string) *image.Image {
	width, height := src.Rect.Max.X, src.Rect.Max.Y;

	average_luminosity, _ := utils.GetAverageBrightnessOfRGBA(src, width, height);
	screenResolution := image.Rect(0, 0, width, height);

	dc := composeMinimalistFrameRGBA(src, *font, text, screenResolution, average_luminosity);
	dcImg := dc.Image();

	return &dcImg;
}

func composeMinimalistFrameRGBA(
	img *image.RGBA,
	font font.Face,
	text string, 
	resolution image.Rectangle,
	average_luminosity uint32,
) *gg.Context {
	screenWidth := resolution.Max.X;
	screenHeight := resolution.Max.Y;

	dc := gg.NewContextForImage(img);

	var r, g, b int;
	if average_luminosity > 150 {
		r = 0;
		g = 0;
		b = 0;
	} else {
		r = 255;
		g = 255;
		b = 255;
	}

	dc.SetFontFace(font);

	offset := 10.0;
	dc.SetRGBA255(r, g, b, 255);
	dc.DrawRectangle(0 + offset, 0 + offset, float64(screenWidth) - (10 + offset), float64(screenHeight) - (10 + offset));
	dc.SetLineWidth(1);
	dc.Stroke();

	dc.SetRGBA255(r, g, b, 255);
	dc.DrawStringWrapped(
		text,
		float64(dc.Width()), // x
		float64(dc.Height()) / 2, // y
		/*
			The anchor point is x - w * ax, y - h * ay, 
				where w, h is the size of the image. 
			Use ax=0.5, ay=0.5 to center 
				the image at the specified point.
		*/
		1.1, // ax (anchor x)
		0.5, // ay (anchor y)
		float64(dc.Width()) / 2, // width
		1.2, // line spacing
		gg.AlignRight,
	);
	return dc;
}
