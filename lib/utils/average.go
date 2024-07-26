package utils

import (
	"image"
)

func GetAverageBrightnessOfRGBA(img *image.RGBA, w int, h int) (uint32, uint32) {
	var average_luminosity uint32 = 0;
	var pixels_sampled uint32 = 0;
	w_interval := w / 8
	h_interval := h / 8

	for x := 0; x < w; x += w_interval {
		for y := 0; y < h; y += h_interval {
			r, g, b, _ := (*img).At(x, y).RGBA();
			average_luminosity += uint32((0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)) / 256);
			pixels_sampled++;
		}
	}
	average_luminosity /= pixels_sampled;

	return average_luminosity, pixels_sampled;
}

func GetAverageBrightnessOfPalettedImage(img *image.Paletted, w int, h int) (uint32, uint32) {
	var average_luminosity uint32 = 0;
	var pixels_sampled uint32 = 0;
	w_interval := w / 8
	h_interval := h / 8

	for x := 0; x < w; x += w_interval {
		for y := 0; y < h; y += h_interval {
			r, g, b, _ := (*img).At(x, y).RGBA();
			average_luminosity += uint32((0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)) / 256);
			pixels_sampled++;
		}
	}
	average_luminosity /= pixels_sampled;

	return average_luminosity, pixels_sampled;
}
