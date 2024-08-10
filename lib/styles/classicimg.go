package styles

import (
	"image"
	"image/color"

	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// for use in images.
func ModifyQuoteImage(src *image.Image, gradient *image.Image, font *font.Face, small_font *font.Face) *gg.Context {
	screenResolution := image.Rect(0, 0, 1280, 720);

	long := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur urna ligula, pellentesque eu risus nec, aliquam malesuada nulla. Duis vel suscipit velit.";
	// short := "hey there buddy."
	return ComposeQuoteFrameImage(src, gradient, *font, *small_font, screenResolution, long);
}

// img should be below 720x720 to reap LanczosResampling's speed in upscaling.
// gradient should be exactly 1280x720.
func ComposeQuoteFrameImage(
	img *image.Image,
	gradient *image.Image,
	font font.Face,
	small_font font.Face,
	resolution image.Rectangle,
	text string,
) *gg.Context { // 720p
	resized := *image.NewRGBA(resolution);
	g := gift.New(
		gift.Resize(720, 0, gift.LanczosResampling), // fastest upscaling.
		gift.Grayscale(),
		gift.Brightness(-20),
	);
	g.Draw(&resized, *img);

	text_x := float64(resolution.Max.X / 2);
	text_y := float64(resolution.Max.Y / 2);
	wrap_width := float64(600);
	
	dc := gg.NewContext(resolution.Max.X, resolution.Max.Y);
	dc.DrawImage(&resized, 0, 0);
	dc.DrawImage(*gradient, 0, 0);

	dc.SetFontFace(font);

	// unclean code below.
	// there's a lot of constants and other stuff to manually arrange the elements into something that looks good.

	// sees the height of the wordwrapped text.
	// gg doesn't provide functionality to see the height of the word wrapped text so i resorted to this.
	_, font_height := dc.MeasureString(text);
	// i offsetted it by 400 because of some discrepancy between [WordWrap and MeasureString] and [DrawStringWrapped's ay=0.5].
	// check DrawStringWrapped below.
	s := dc.WordWrap(text, wrap_width + 400);
	var total_font_height float64 = 0;
	for range s {
		total_font_height += font_height;
	}
	dc.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255});
	dc.SetLineWidth(3);
	dc.DrawLine(0, text_x + total_font_height, 1280, text_x + total_font_height);
	dc.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 255});
	dc.DrawPoint(text_x, text_y, 5);

	dc.SetColor(color.White);
	// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the text.
	// this works properly dw about this. (call this "quote string")
	dc.DrawStringWrapped(text, text_x, text_y, 0, 0.5, wrap_width, 1, gg.AlignRight);

	dc.SetColor(color.White);
	dc.SetFontFace(small_font);
	// the offending code (call this "author")
	// i wanted this to attach itself below the DrawStringWrapped above, (text_y + total_font_height).
	// however, i forget to account for the anchor y (ay=0.5), so I needed to put multiple constants everywhere to make it semi decent.
	dc.DrawStringWrapped("- alice", text_x + 300, text_y + total_font_height, 0, 0, 300, 1, gg.AlignRight);

	// a cleaner solution would be to:
	// 		grab the "size of the text (w, h)"  // small this is your task.
	//		and adjust the "author"'s y value  	// i'll figure this out.
	// 			to match "quote string"'s actual y value.
	return dc;
}
