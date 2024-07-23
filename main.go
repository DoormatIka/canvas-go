package main

import (
	"bufio"
	"flag"
	"fmt"

	"image/gif"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"canvas/lib/styles"

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

func importGif(filename string) *bufio.Reader {
	inputFile, err := os.Open(filename);
	if err != nil {
		panic(err);
	}
	r := bufio.NewReader(inputFile);
	return r;
}


var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file");
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

	font, err := gg.LoadFontFace("./fonts/Lora-Italic.ttf", 45)
	if err != nil {
		panic(err);
	}
	sky_gif := importGif("./images/okina-matara-junko.gif")
	g, err := gif.DecodeAll(sky_gif);
	if err != nil {
		panic(err);
	}
	f, err := os.Create("./images/res_speen.gif");
	if err != nil {
		panic(err);
	}
	defer f.Close();
	
	defer timer("gif")();

	modified_gif := styles.ModifyMinimalistGif(g, &font);
	if err := gif.EncodeAll(f, modified_gif); err != nil {
		panic(err);
	}
}
