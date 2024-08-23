package styles

import (
	"image"
	"image/gif"

	"golang.org/x/image/font"
	"github.com/disintegration/gift"
	"canvas/lib/utils"
)

// text string, author string, src *image.Image, gradient *image.Image, font *font.Face, small_font *font.Face
func ModifyClassicGif(src *gif.GIF, font *font.Face, small_font *font.Face, text string, author string, gradient *image.Image) *gif.GIF {
	newGif := &gif.GIF{};

	width := int(float32(src.Config.Height) * 1.77778);
	grad := image.NewRGBA(image.Rect(0, 0, width, src.Config.Height));
	resizer := gift.New(gift.Resize(0, src.Config.Height, gift.LinearResampling))
	resizer.Draw(grad, *gradient);

    quantizer := utils.NewOctreeQuantizer()
	utils.AddColorsToQuantizer(quantizer, src);
    colorCount := 256; // colors. 256 before.
	palette := quantizer.MakePalette(colorCount)
	colorPalette := utils.ConvertToColorPalette(palette);
	screenResolution := image.Rect(0, 0, width, src.Config.Height);

	// the reused image.
	reusedImage := image.NewPaletted(screenResolution, colorPalette);

	for i := 0; i < len(src.Image); i++ {
		img := src.Image[i];
		delay := src.Delay[i];
		disposal := src.Disposal[i];

		regularImage := image.NewPaletted(screenResolution, colorPalette);

		dc := composeClassicImage(img, grad, *font, *small_font, screenResolution, text, author, 400, 9);
		dcImg := dc.Image().(*image.RGBA);
		bounds := dcImg.Bounds();

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				i := y*dcImg.Stride + x*4
				// Extract the RGBA values directly from the Pix slice
				r := dcImg.Pix[i+0]
				g := dcImg.Pix[i+1]
				b := dcImg.Pix[i+2]
				a := dcImg.Pix[i+3]
				if a == 0 {
					continue;
				}

				color := utils.NewColor(int(r), int(g), int(b), int(a));
				index := quantizer.GetPaletteIndex(color);
				if disposal == gif.DisposalNone {
					i := reusedImage.PixOffset(x, y);
					reusedImage.Pix[i] = uint8(index);
				} else {
					regularImage.SetColorIndex(x, y, uint8(index));
				}
			}
		}

		if disposal == gif.DisposalNone {
			copiedReusedImage := image.NewPaletted(screenResolution, colorPalette)
			copy(copiedReusedImage.Pix, reusedImage.Pix)
			newGif.Image = append(newGif.Image, copiedReusedImage);
		} else {
			newGif.Image = append(newGif.Image, regularImage);
		}
		newGif.Delay = append(newGif.Delay, delay);
		newGif.Disposal = append(newGif.Disposal, 1);
	}

	newGif.LoopCount = src.LoopCount;
	newGif.Config.Height = src.Config.Height;
	newGif.Config.Width = width;
	newGif.Config.ColorModel = src.Config.ColorModel;
	newGif.BackgroundIndex = src.BackgroundIndex;

	return newGif;
}
