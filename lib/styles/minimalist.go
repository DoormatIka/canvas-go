package styles

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"sync"

	"golang.org/x/image/font"

	"github.com/disintegration/gift"
	"github.com/ericpauley/go-quantize/quantize"
	"github.com/fogleman/gg"

	"canvas/lib/utils"
)

type Frame struct {
	palettedImage *image.Paletted
	delay int
	disposal byte
}

func inLoopMinimalistGif(
	wg *sync.WaitGroup,
	frameChan chan Frame,
	src *gif.GIF,
	font font.Face,
	i int,
) {
	defer wg.Done();

	quantizer := quantize.MedianCutQuantizer{};
	screenResolution := image.Rect(0, 0, src.Config.Width, src.Config.Height);

	img := src.Image[i];
	delay := src.Delay[i];
	disposal := src.Disposal[i];

	dc := ComposeMinimalistFrameGif(img, font, "You're beautiful.", screenResolution);

	dc_img := dc.Image();
	bounds := dc_img.Bounds();
	img_palette := quantizer.Quantize(make(color.Palette, 0, 256), dc_img);
	palettedImage := image.NewPaletted(bounds, img_palette);
	draw.Draw(palettedImage, bounds, dc_img, bounds.Min, draw.Src);

	frameChan <- Frame {palettedImage, delay, disposal}
}

func ModifyMinimalistGif(src *gif.GIF, font *font.Face) *gif.GIF {
	newGif := &gif.GIF{};

	var wg sync.WaitGroup;

	frameChan := make(chan Frame, len(src.Image) + 1);
	// making a goroutine frame by frame.
	// refactor this later to avoid the creation of the goroutine being a bottleneck
	
	for i := 0; i < len(src.Image); i++ {
		wg.Add(1);
		go inLoopMinimalistGif(&wg, frameChan, src, *font, i);
	}
	wg.Wait();
	close(frameChan);

	for v := range frameChan {
		newGif.Image = append(newGif.Image, v.palettedImage);
		newGif.Delay = append(newGif.Delay, v.delay);
		newGif.Disposal = append(newGif.Disposal, v.disposal);
	}

	return newGif;
}

// this is the frame function optimized for gif frames
// this automatically adapts to the image resolutions
func ComposeMinimalistFrameGif(
	img *image.Paletted,
	font font.Face,
	text string, 
	resolution image.Rectangle,
) *gg.Context {
	screenWidth := resolution.Max.X;
	screenHeight := resolution.Max.Y;

	dc := gg.NewContext(screenWidth, screenHeight);

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
		dst := image.NewRGBA(resizer.Bounds(image.Rect(0, 0, screenWidth, screenHeight)))
		resizer.Draw(dst, img)
		dc.DrawImage(dst, 0, 0);
	} else {
		dc.DrawImage(img, 0, 0);
	}

	dc_img := dc.Image();
	average_luminosity, _ := utils.GetAverageBrightnessOfImage(&dc_img, screenWidth, screenHeight);

	var r, g, b int;
	if average_luminosity > 130 {
		r = 0;
		g = 0;
		b = 0;
	} else {
		r = 255;
		g = 255;
		b = 255;
	}

	/*
	dc.SetRGBA255(r, g, b, 200);
	dc.SetFontFace(font);
	dc.DrawString(fmt.Sprintf("lum: %v", average_luminosity), 0, float64(screenHeight) / 2);
	dc.SetRGBA255(r, g, b, 200);
	dc.SetFontFace(font);
	dc.DrawString(fmt.Sprintf("pixels: %v", pixels_sampled), 0, float64(screenHeight) / 2 + 50);
	*/

	offset := 10.0;
	dc.SetRGBA255(r, g, b, 255);
	dc.DrawRectangle(0 + offset, 0 + offset, float64(screenWidth) - (10 + offset), float64(screenHeight) - (10 + offset));
	dc.SetLineWidth(2);
	dc.Stroke();

	dc.SetRGBA255(r, g, b, 200);
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


// this is meant for images. provides greater quality with less focus on speed
func ComposeMinimalistFrameImg(img image.Image, font font.Face, text string) *gg.Context {
	screenWidth := 512 * 2;
	screenHeight := 512 * 2;

	dc := gg.NewContext(screenWidth, screenHeight);

	imgWidth := img.Bounds().Dx();
	imgHeight := img.Bounds().Dy()

	var newWidth, newHeight int 
	// it should handle multiple image resolutions.
	if imgHeight > imgWidth {
		newHeight = 0
		newWidth = screenWidth
	} else {
		newHeight = screenHeight
		newWidth = 0
	}

	filter := gift.New(
		gift.Resize(newWidth, newHeight, gift.LinearResampling),
	)
	dst := image.NewRGBA(filter.Bounds(img.Bounds()))
	filter.Draw(dst, img)
	dc.DrawImage(dst, 0, 0);

	dc.SetRGB(1, 1, 1);
	dc.DrawRectangle(0, 0, float64(screenWidth), float64(screenHeight));
	dc.SetLineWidth(10);
	dc.Stroke();

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 200})
	dc.SetFontFace(font);
	
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
