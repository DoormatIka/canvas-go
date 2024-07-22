package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"

	"canvas/lib"

	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
)

// timer returns a function that prints the name argument and
// the elapsed time between the call to timer and the call to
// the returned function. The returned function is intended to
// be used in a defer statement:
//
//   defer timer("sum")()
func timer(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

func gradient_test() {
	dc := gg.NewContext(500, 500)

	grad := gg.NewLinearGradient(0, 0, float64(dc.Width()), float64(dc.Height()))
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.SetFillStyle(grad)
	// dc.SetStrokeStyle(grad);
	// dc.SetLineWidth(10);
	// dc.Stroke();
	dc.Fill()

	dc.DrawCircle(float64(dc.Height())/2, float64(dc.Width())/2, 20)
	dc.SetColor(color.White)
	dc.Fill()

	dc.SavePNG("out_gradient.png")
	fmt.Printf("Hello there.\n")
}

func sizer(filename string, size int) error {
	canvas, err := gg.LoadImage("./images/" + filename)
	if err != nil {
		return err
	}
	filter := gift.New(gift.ResizeToFit(size, size, gift.LanczosResampling))
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	filter.Draw(dst, canvas)

	f, err := os.Create(fmt.Sprintf("images/%v_youmu.png", size))
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, dst)
	if err != nil {
		fmt.Printf("Error encoding the png.")
		return err
	}
	return nil
}

func write_to_png(dst *image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, *dst)
	if err != nil {
		fmt.Printf("Error encoding the png.")
		return err
	}
	return nil
}

func mask(filename string) (*gg.Context, error) {
	canvas, err := gg.LoadImage("./images/" + filename)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContext(512, 512);
	filter := gift.New(
		gift.Resize(512, 0, gift.BoxResampling),
		gift.Contrast(-20),
	)
	// dc.DrawRoundedRectangle(0, 0, 512, 512, 64*2)
	dc.DrawCircle(float64(dc.Height()) / 2, float64(dc.Width()) / 2, 256);
	dc.Clip() // gets the clip from DrawRoundedRectangle??

	dst := image.NewRGBA(filter.Bounds((canvas).Bounds()))
	filter.Draw(dst, canvas)

	/*
		img := dst.SubImage(dst.Rect);
		write_to_png(&img, "dst.png");
	*/

	dc.DrawImage(dst, 0, 0)
	return dc, nil
}

// this is the frame function optimized for gif frames
// this automatically adapts to the image resolutions
func minimalist_frame_gif(img *image.Paletted, font font.Face, text string) *gg.Context {
	screenWidth := img.Bounds().Max.X;
	screenHeight := img.Bounds().Max.Y;

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
		dst := image.NewRGBA(resizer.Bounds(img.Bounds()))
		resizer.Draw(dst, img)
		dc.DrawImage(dst, 0, 0);
	} else {
		dc.DrawImage(img, 0, 0);
	}

	var average_luminosity uint32 = 0;
	var pixels_sampled uint32 = 0;
	underlying_img := dc.Image();
	for x := dc.Width() / 2; x < dc.Width(); x += 50 {
		for y := dc.Height() / 4; y < dc.Height(); y += 50 {
			r, g, b, _ := underlying_img.At(x, y).RGBA();
			average_luminosity += uint32((0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)) / 256);
			pixels_sampled++;
		}
	}
	average_luminosity /= pixels_sampled;

	dc.SetFontFace(font);
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 100})
	dc.DrawString(fmt.Sprintf("lum: %v", average_luminosity), 0, float64(screenHeight) / 2);
	dc.SetFontFace(font);
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 100})
	dc.DrawString(fmt.Sprintf("pixels: %v", pixels_sampled), 0, float64(screenHeight) / 2 + 50);

	dc.SetRGB(1, 1, 1);
	dc.DrawRectangle(0, 0, float64(screenWidth), float64(screenHeight));
	// dc.DrawRectangle(float64(dc.Width()) / 2, float64(dc.Height()) / 4, float64(dc.Width()), float64(dc.Height()) / 2);
	dc.SetLineWidth(10);
	dc.Stroke();

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 200})
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
func minimalist_frame_img(img image.Image, font font.Face, text string) *gg.Context {
	screenWidth := 512 * 2;
	screenHeight := 512 * 2;

	dc := gg.NewContext(screenWidth, screenHeight);

	imgWidth := img.Bounds().Dx()
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
		strings.Repeat(text, 1),
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

func import_gif(filename string) *bufio.Reader {
	inputFile, err := os.Open(filename);
	if err != nil {
		panic(err);
	}
	r := bufio.NewReader(inputFile);
	return r;
}

func main() {
	font, err := gg.LoadFontFace("./fonts/Lora-Italic.ttf", 45)
	if err != nil {
		panic(err);
	}
	sky_gif := import_gif("./images/okina-matara-junko.gif");
	g, err := gif.DecodeAll(sky_gif);
	if err != nil {
		panic(err);
	}
	
	f, err := os.Create("./images/res_bright.gif");
	if err != nil {
		panic(err);
	}
	defer f.Close();

	defer timer("gif")();

	newGif := &gif.GIF{};
	quantizer := mediancut.MedianCutQuantizer{NumColor: 256};

	for i := 0; i < len(g.Image); i++ {
		img := g.Image[i]; // PalettedImage because
							// this image was grabbed from a gif.
		delay := g.Delay[i];

		dc := minimalist_frame_gif(img, font, "You're beautiful.");

		dc_img := dc.Image();
		bounds := dc_img.Bounds();
		palettedImage := image.NewPaletted(bounds, nil);
		quantizer.Quantize(palettedImage, bounds, dc_img, bounds.Min);

		newGif.Image = append(newGif.Image, palettedImage);
		newGif.Delay = append(newGif.Delay, delay);
	}

	giferr := gif.EncodeAll(f, newGif);
	if giferr != nil {
		panic(giferr)
	}
	

	/*
	img, err := gg.LoadImage("./images/Moon.png");
	if err != nil {
		return;
	}
	dc.SavePNG("./images/quote/image_minimalist_frame_2.png");
	files := []string{"100_youmu.png", "200_youmu.png", "300_youmu.png", "400_youmu.png", "500_youmu.png", "full_youmu.jpg"}
	for _, file := range files {
		dc, err := mask(file)
		if err != nil {
			fmt.Printf("%v couldn't be masked. %v", file, err)
			return
		}
		dc.SavePNG(fmt.Sprintf("./images/mask/mask_%v", file))
	}
	*/
}
