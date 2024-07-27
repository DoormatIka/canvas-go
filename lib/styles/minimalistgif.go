package styles

import (
	"image"
	"image/gif"
	"sort"
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
	var wg sync.WaitGroup;

	frameChan := make(chan GifFrame, len(src.Image));
	// making a goroutine frame by frame.
	// refactor this later to avoid the creation of the goroutine being a bottleneck

	average_luminosity, _ := utils.GetAverageBrightnessOfPalettedImage(src.Image[0], src.Config.Width, src.Config.Height);
	for i := 0; i < len(src.Image); i++ {
		wg.Add(1);

		screenResolution := image.Rect(0, 0, src.Config.Width, src.Config.Height);
		go func (i int) {
			defer wg.Done();

			img := src.Image[i];
			delay := src.Delay[i];
			disposal := src.Disposal[i];

			fontMutex.Lock()
			dc := ComposeMinimalistFrameGif(img, *font, text, screenResolution, average_luminosity);
			fontMutex.Unlock()

			dcImg := dc.Image();
			// Initialize the octree with a color depth of 4
			hexatree := utils.NewHexaTree(8) // Adjust the color depth as needed
			utils.BuildTree(dcImg, hexatree);
			hexatree.Reduce();
			// Build the palette from the reduced hexatree (reduces from image automatically)
			hexatree.BuildPalette();
			// Convert the image to a paletted image using the hexatree
			palettedImage := hexatree.ConvertToPaletted(dcImg)

			frameChan <- GifFrame {palettedImage: palettedImage, delay: delay, disposal: disposal, index: i}
		}(i)
	}
	wg.Wait();
	close(frameChan);

	frames := []GifFrame{};

	// three loops in a row. there has to be a better way!
	for v := range frameChan {
		frames = append(frames, v);
	}
	sort.SliceStable(frames, func(i, j int) bool {
		return frames[i].index < frames[j].index;
	})
	for _, v := range frames {
		newGif.Image = append(newGif.Image, v.palettedImage);
		newGif.Delay = append(newGif.Delay, v.delay);
		newGif.Disposal = append(newGif.Disposal, v.disposal);
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
