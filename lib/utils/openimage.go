package utils;

import (
	"image"
	"os"
	"bufio"
)

func OpenImage(fp string) image.Image {
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
