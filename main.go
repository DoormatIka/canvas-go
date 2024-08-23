package main

import (
	// "image"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"

	"canvas/lib/styles"
	"canvas/lib/utils"

	"github.com/fogleman/gg"

	_ "image/png"
)

var gradient = utils.OpenImage("./images/quote/qgradient.png");
var big_classic_font, _ = gg.LoadFontFace("./fonts/Mirador-SemiBold.ttf", 25 * 2);
var small_classic_font, _ = gg.LoadFontFace("./fonts/Mirador-BookItalic.ttf", 15 * 2);

var big_classicgif_font, _ = gg.LoadFontFace("./fonts/Mirador-SemiBold.ttf", 25);
var small_classicgif_font, _ = gg.LoadFontFace("./fonts/Mirador-BookItalic.ttf", 15);

var gifminimalist_font, _ = gg.LoadFontFace("./fonts/Lora-Italic.ttf", 25)


func getImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url);
	if err != nil {
		return nil, err;
	}
	defer resp.Body.Close();

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status);
	}
	img, _, err := image.Decode(resp.Body);
	if err != nil {
		return nil, err;
	}

	return img, nil;
}

type Meta struct {
	Url string `json:"avatar_url"`
	Author string `json:"author"`
	Text string `json:"text"`
}

func sendClassicImage(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body);
	var meta Meta;


	if err := json.Unmarshal(reqBody, &meta); err != nil {
		http.Error(w, "Failed to parse metadata.", http.StatusBadRequest);
		return;
	}
	img, err := getImageFromURL(meta.Url);
	if err != nil {
		http.Error(w, "Can't get image from URL. " + err.Error(), http.StatusBadRequest);
		return;
	}
	imgData := styles.ModifyClassicImage(meta.Text, meta.Author, img, gradient, &big_classic_font, &small_classic_font);
	w.Header().Set("Content-Type", "image/png");
	imgData.EncodePNG(w);
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"));
}

func main() {
	http.HandleFunc("/ping", ping);
	http.HandleFunc("/quote", sendClassicImage);
	http.ListenAndServe(":8080", nil);
	println("Started server on localhost:8080");
}
