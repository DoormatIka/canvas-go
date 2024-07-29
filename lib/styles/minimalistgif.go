package styles

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"sync"

	"golang.org/x/image/font"

	"github.com/fogleman/gg"

	"canvas/lib/utils"
)

type GifFrame struct {
	palettedImage *image.Paletted
	delay int
	disposal byte
	index int
}

func ModifyMinimalistGif(src *gif.GIF, font *font.Face, text string) *gif.GIF {
	newGif := &gif.GIF{};

	fontMutex := sync.Mutex{};

    // Initialize the octree quantizer
    quantizer := utils.NewOctreeQuantizer()
    // Add colors from each frame to the quantizer
	utils.AddColorsToQuantizer(quantizer, src);
    // Generate the palette
    colorCount := 64
    palette := quantizer.MakePalette(colorCount)
    // Convert []Color to color.Palette
    colorPalette := utils.ConvertToColorPalette(palette)
    // Add a transparent color to the end of the palette
    colorPalette = append(colorPalette, color.RGBA{0, 0, 0, 0})
    // Create a new GIF with quantized frames

	average_luminosity, _ := utils.GetAverageBrightnessOfPalettedImage(src.Image[0], src.Config.Width, src.Config.Height);
	// making a goroutine frame by frame.
	// refactor this later to avoid the creation of the goroutine being a bottleneck
	for i := 0; i < len(src.Image); i++ {
		screenResolution := image.Rect(0, 0, src.Config.Width, src.Config.Height);

		img := src.Image[i];
		delay := src.Delay[i];
		disposal := src.Disposal[i];

		fontMutex.Lock()
		dc := ComposeMinimalistFrameGif(img, *font, text, screenResolution, average_luminosity);
		fontMutex.Unlock()

		dcImg := dc.Image();
		bounds := dcImg.Bounds();
		quantizedFrame := image.NewPaletted(bounds, colorPalette);
		transparentIndex := len(colorPalette) - 1;

		if i > 0 {
			switch src.Disposal[i-1] {
            case gif.DisposalPrevious:
                draw.Draw(quantizedFrame, bounds, newGif.Image[i-1], image.Point{}, draw.Over)
            case gif.DisposalBackground:
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
                    for x := bounds.Min.X; x < bounds.Max.X; x++ {
                        quantizedFrame.SetColorIndex(x, y, uint8(transparentIndex))
                    }
                }
			}
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := dcImg.At(x, y).RGBA()
				color := utils.NewColor(int(r>>8), int(g>>8), int(b>>8), int(a>>8))
				if a == 0 {
					quantizedFrame.SetColorIndex(x, y, uint8(transparentIndex))
				} else {
					index := quantizer.GetPaletteIndex(color)
					quantizedFrame.SetColorIndex(x, y, uint8(index))
				}
			}
		}

		newGif.Image = append(newGif.Image, quantizedFrame);
		newGif.Delay = append(newGif.Delay, delay);
		newGif.Disposal = append(newGif.Disposal, disposal);
	}
	newGif.LoopCount = src.LoopCount;
	newGif.Config.Height = src.Config.Height;
	newGif.Config.Width = src.Config.Width;
	newGif.Config.ColorModel = src.Config.ColorModel;
	newGif.BackgroundIndex = src.BackgroundIndex;

	return newGif;
}

// this automatically adapts to the image resolutions
func ComposeMinimalistFrameGif(
	img *image.Paletted,
	font font.Face,
	text string, 
	resolution image.Rectangle,
	average_luminosity uint32,
) *gg.Context {
	gifWidth := resolution.Max.X;
	gifHeight := resolution.Max.Y;

	imgXOrigin := img.Bounds().Min.X;
	imgYOrigin := img.Bounds().Min.Y;
	imgWidth := img.Bounds().Max.X;
	imgHeight := img.Bounds().Max.Y;

	var dc *gg.Context;
	if 0 != imgXOrigin || 0 != imgYOrigin || imgWidth != gifWidth || imgHeight != gifHeight {
		dc = gg.NewContext(gifWidth, gifHeight);
		dc.DrawImage(img, 0, 0); // significantly slower.
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
	dc.DrawRectangle(0 + offset, 0 + offset, float64(gifWidth) - (10 + offset), float64(gifHeight) - (10 + offset));
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
