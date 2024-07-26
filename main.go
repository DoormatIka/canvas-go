package main

import (
	"bufio"
	"flag"

	"image/gif"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"canvas/lib/styles"
	"canvas/lib/utils"

	"github.com/fogleman/gg"
)

func importGif(filename string) *bufio.Reader {
	inputFile, err := os.Open(filename);
	if err != nil {
		panic(err);
	}
	r := bufio.NewReader(inputFile);
	return r;
}

func runGif() {
	font, err := gg.LoadFontFace("./fonts/Lora-Italic.ttf", 25)
	if err != nil {
		panic(err);
	}
	defer utils.Timer("sky gif")();

	sky_gif := importGif("./images/gifs/night-sky-time-lapse.gif");
	g, err := gif.DecodeAll(sky_gif);
	if err != nil {
		panic(err);
	}
	f, err := os.Create("./images/gifs/minimalist/out_night-sky-time-lapse.gif");
	defer f.Close();

	modified_gif := styles.ModifyMinimalistGif(g, &font);
	if err := gif.EncodeAll(f, modified_gif); err != nil {
		panic(err);
	}
	/*
	entries, err := os.ReadDir("./images/gifs/");
	if err != nil {
		panic(err);
	}

	for _, v := range entries {
		if v.IsDir() {
			continue;
		}

		sky_gif := importGif("./images/gifs/" + v.Name());
		g, err := gif.DecodeAll(sky_gif);
		if err != nil {
			panic(err);
		}
		f, err := os.Create("./images/gifs/minimalist/out_" + v.Name());
		defer f.Close();

		if err != nil {
			panic(err);
		}
		modified_gif := styles.ModifyMinimalistGif(g, &font);
		if err := gif.EncodeAll(f, modified_gif); err != nil {
			panic(err);
		}
	}
	*/
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

	runGif();
}
