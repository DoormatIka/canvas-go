package utils

import (
	// "errors"
	"image"
	// "image/color"
	"image/draw"
)

// Ripped out of Go's drawPaletted function to modify it.

func clamp(i int32) int32 {
    // Clamp to 0 if i is negative
    clamped := i &^ (i >> 31)
    // Clamp to 0xffff if i is greater than 0xffff
    clamped |= (0xffff - clamped) >> 31
    return clamped & 0xffff
}

func sqDiff(x, y int32) uint32 {
	d := uint32(x - y);
	return (d * d) >> 2
}

// An attempt to optimize drawPaletted using a bin best first method. (kd tree for now)
//
// This function assumes that you already dithered the source image to suit my purposes.
func GoDrawPaletted(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, floydSteinberg bool) {
	// If dst is an *image.Paletted, we have a fast path for dst.Set and
	// dst.At. The dst.Set equivalent is a batch version of the algorithm
	// used by color.Palette's Index method in image/color/color.go,	
	palette, pix, stride := [][4]int32(nil), []byte(nil), 0
	if p, ok := dst.(*image.Paletted); ok {
		palette = make([][4]int32, len(p.Palette))
		for i, col := range p.Palette {
			r, g, b, a := col.RGBA()
			palette[i][0] = int32(r)
			palette[i][1] = int32(g)
			palette[i][2] = int32(b)
			palette[i][3] = int32(a)
		}
		// p.Pix passes a reference into pix, that's how its being modified.
		pix, stride = p.Pix[p.PixOffset(r.Min.X, r.Min.Y):], p.Stride
	}

	// Loop over each source pixel.
	for y := 0; y != r.Dy(); y++ {
		for x := 0; x != r.Dx(); x++ {
			// er, eg and eb are the pixel's R,G,B values plus the
			// optional Floyd-Steinberg error.
			sr, sg, sb, sa := src.At(sp.X+x, sp.Y+y).RGBA()
			er, eg, eb, ea := int32(sr), int32(sg), int32(sb), int32(sa)

			// Find the closest palette color in Euclidean R,G,B,A space:
			// the one that minimizes sum-squared-difference.
			// TODO(nigeltao): consider smarter algorithms.
			bestIndex, bestSum := 0, uint32(1<<32-1)
			for index, p := range palette {
				sum := sqDiff(er, p[0]) + sqDiff(eg, p[1]) + sqDiff(eb, p[2]) + sqDiff(ea, p[3])
				if sum < bestSum {
					bestIndex, bestSum = index, sum
					if sum == 0 {
						break
					}
				}
			}

			// modifies the pixels.
			pix[y*stride+x] = byte(bestIndex)
		}
	}
}
