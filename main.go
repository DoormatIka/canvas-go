package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"runtime/trace"

	"canvas/lib/styles"
	"canvas/lib/utils"

	"github.com/fogleman/gg"
)

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
	modified_gif := styles.ModifyMinimalistGif(g, &font, "The sky looks nice doesn't it?");
	if err := gif.EncodeAll(f, modified_gif); err != nil {
		panic(err);
	}
}

func runImageForQuote() {
	font, err := gg.LoadFontFace("./fonts/Lora-Italic.ttf", 25)
	if err != nil {
		panic(err);
	}
	defer utils.Timer("sky gif")();

	image_file, err := os.Open("./images/Moon.png");
	if err != nil {
		panic(err);
	}
	im := bufio.NewReader(image_file);
	decoded_im, _, err := image.Decode(im);
	quote_img := styles.ModifyQuoteImage(decoded_im, &font);
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

	// runImageForQuote();
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
	// utils.RunOctQuant();
}
