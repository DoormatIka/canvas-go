package styles

import (
	"image"
	"image/gif"
	"sort"

	"sync"

	"golang.org/x/image/font"

	"github.com/disintegration/gift"
	"github.com/fogleman/gg"

	"canvas/lib/utils"
)

type GifFrame struct {
	palettedImage *image.Paletted
	delay int
	disposal byte
	index int
}

func ModifyMinimalistGif(src *gif.GIF, font *font.Face) *gif.GIF {
	newGif := &gif.GIF{};

	var fontMutex sync.Mutex;
	var wg sync.WaitGroup;

	frameChan := make(chan GifFrame, len(src.Image));
	// making a goroutine frame by frame.
	// refactor this later to avoid the creation of the goroutine being a bottleneck

	average_luminosity, _ := utils.GetAverageBrightnessOfImage[*image.Paletted](src.Image[0], src.Config.Width, src.Config.Height);
	for i := 0; i < len(src.Image); i++ {
		wg.Add(1);

		go func (i int) {
			defer wg.Done();

			screenResolution := image.Rect(0, 0, src.Config.Width, src.Config.Height);

			img := src.Image[i];
			delay := src.Delay[i];
			disposal := src.Disposal[i];

			fontMutex.Lock();
			dc := ComposeMinimalistFrame(img, *font, "You're amazing at what you do.", screenResolution, average_luminosity);
			fontMutex.Unlock();

			dcImg := dc.Image();
			// Initialize the octree with a color depth of 4
			octree := utils.NewOctree(4) // Adjust the color depth as needed
			utils.BuildTree(dcImg, octree)
			// Reduce the octree to fit a 256 color palette
			octree.Reduce()
			// Build the palette from the reduced octree
			octree.BuildPalette();
			// Convert the image to a paletted image using the octree
			palettedImage := octree.ConvertToPaletted(dcImg)

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
func ComposeMinimalistFrame(
	img *image.Paletted,
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
