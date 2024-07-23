package utils

import (
	"image"
)


func GetAverageBrightnessOfImage[K *image.Paletted | image.RGBA](img *image.Paletted, w int, h int) (uint32, uint32) {
	var average_luminosity uint32 = 0;
	var pixels_sampled uint32 = 0;

	for x := w / 2; x < w; x += w / 8 {
		for y := h / 4; y < h; y += y / 8 {
			r, g, b, _ := (*img).At(x, y).RGBA();
			average_luminosity += uint32((0.299 * float64(r) + 0.587 * float64(g) + 0.114 * float64(b)) / 256);
			pixels_sampled++;
		}
	}
	average_luminosity /= pixels_sampled;

	return average_luminosity, pixels_sampled;
}
