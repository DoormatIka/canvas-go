package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"os"
	// "regexp"
	"runtime/pprof"
	"runtime/trace"

	"canvas/lib/styles"
	"canvas/lib/utils"

	"github.com/fogleman/gg"
)

func openImage(fp string) image.Image {
	image_file, err := os.Open(fp);
	if err != nil {
		panic(err);
	}
	im := bufio.NewReader(image_file);
	decoded_im, _, err := image.Decode(im);
	if err != nil {
		panic(err);
	}
	if err := image_file.Close(); err != nil {
		panic(err);
	}
	return decoded_im;
}

func outputPalette(palette color.Palette, filename string) {
    pixelSize := 10
    paletteImg := image.NewRGBA(image.Rect(0, 0, 8*pixelSize, ((len(palette)+7)/8)*pixelSize))
    for i, c := range palette {
        x := (i % 8) * pixelSize
        y := (i / 8) * pixelSize
        for dx := 0; dx < pixelSize; dx++ {
            for dy := 0; dy < pixelSize; dy++ {
                r, g, b, a := c.RGBA()
                paletteImg.Set(x+dx, y+dy, color.RGBA{
                    R: uint8(r >> 8),
                    G: uint8(g >> 8),
                    B: uint8(b >> 8),
                    A: uint8(a >> 8),
                })
            }
        }
    }
    paletteFile, err := os.Create("./images/minimalist/gifs/debug/" + filename + ".png");
    if err != nil {
        fmt.Println("Error creating palette image:", err)
        return
    }
    defer paletteFile.Close()

    err = png.Encode(paletteFile, paletteImg)
    if err != nil {
        fmt.Println("Error encoding palette image:", err)
        return
    }

    fmt.Println("Palette image saved as " + filename);
}

func runGifForMinimalist(filename string) {
	font, err := gg.LoadFontFace("./fonts/Lora-Italic.ttf", 25)
	if err != nil {
		panic(err);
	}

	defer utils.Timer(filename)();

	sky_gif_file, err := os.Open("./images/" + filename);
	println(sky_gif_file.Name());
	if err != nil {
		panic(err);
	}
	sky_gif := bufio.NewReader(sky_gif_file);

	g, err := gif.DecodeAll(sky_gif);
	if err != nil {
		panic(err);
	}
	f, err := os.Create("./images/minimalist/gifs/out_" + filename);
	defer f.Close();

	fmt.Printf("Number of frames: %v\n", len(g.Image));

	modified_gif := styles.ModifyMinimalistGif(g, &font, "Death.");
	if err := gif.EncodeAll(f, modified_gif); err != nil {
		panic(err);
	}
}


func runImageForQuote() {
	big_font, err := gg.LoadFontFace("./fonts/Mirador-Bold.ttf", 25 * 2);
	if err != nil {
		panic(err);
	}
	small_font, err := gg.LoadFontFace("./fonts/Mirador-BookItalic.ttf", 10 * 2);
	if err != nil {
		panic(err);
	}
	defer utils.Timer("sky gif")();

	gradient := openImage("./images/quote/qgradient.png");
	pfp := openImage("./images/pfp2.png");

	quote_img := styles.ModifyQuoteImage(&pfp, &gradient, &big_font, &small_font);
	quote_img.SavePNG("./images/quote/out_Moon.png");
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file");
var memorytrace = flag.String("memorytrace", "", "write memory trace to file");
func main() {
	flag.Parse();
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile);
		if err != nil {
			log.Fatal(err);
		}
		pprof.StartCPUProfile(f);
		defer pprof.StopCPUProfile();
	}
	if *memorytrace != "" {
		f, err := os.Create(*memorytrace);
		if err != nil {
			log.Fatal(err);
		}
		trace.Start(f);
		defer trace.Stop();
	}

	runImageForQuote();
	/*
	files, err := os.ReadDir("./images/");
	if err != nil {
		log.Fatal(err);
	}
	reg := regexp.MustCompile(".*\\.gif");
	for _, file := range files {
		if file.IsDir() {
			continue;
		}
		if reg.Match([]byte(file.Name())) {
			runGifForMinimalist(file.Name());
		}
		println();
	}
	*/
}
