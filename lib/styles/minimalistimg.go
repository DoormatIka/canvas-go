package styles

import (
	"image"

	"golang.org/x/image/font"

	"github.com/disintegration/gift"
	"github.com/fogleman/gg"

	"canvas/lib/utils"
)

// for use in images.
func ModifyMinimalistRGBA(src *image.RGBA, font *font.Face) *image.Image {
	width, height := src.Rect.Max.X, src.Rect.Max.Y;

	average_luminosity, _ := utils.GetAverageBrightnessOfRGBA(src, width, height);
	screenResolution := image.Rect(0, 0, width, height);

	dc := ComposeMinimalistFrameImage(src, *font, "You're amazing at what you do.", screenResolution, average_luminosity);
	dcImg := dc.Image();

	return &dcImg;
}

// this automatically adapts to the image resolutions
func ComposeMinimalistFrameImage(
	img *image.RGBA,
	font font.Face,
	text string, 
	resolution image.Rectangle,
	average_luminosity uint32,
) *gg.Context {
	screenWidth := resolution.Max.X;
	screenHeight := resolution.Max.Y;

	var dc *gg.Context;
	if screenHeight > 500 || screenWidth > 500 {
		var newWidth, newHeight int 
		// it should handle multiple image resolutions.
		if screenHeight > screenWidth {
			newHeight = 0
			newWidth = screenWidth
		} else {
			newHeight = screenHeight
			newWidth = 0
		}
		resizer := gift.New( // downscaling
			gift.Resize(newWidth, newHeight, gift.LanczosResampling),
		)
		dst := image.NewRGBA(resizer.Bounds(image.Rect(0, 0, screenWidth, screenHeight)));
		resizer.Draw(dst, img);
		dc = gg.NewContextForImage(dst);
	} else {
		dc = gg.NewContextForImage(img);
	}

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
