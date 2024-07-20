package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
)

func gradient_test() {
	dc := gg.NewContext(500, 500)

	grad := gg.NewLinearGradient(0, 0, float64(dc.Width()), float64(dc.Height()))
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.SetFillStyle(grad)
	// dc.SetStrokeStyle(grad);
	// dc.SetLineWidth(10);
	// dc.Stroke();
	dc.Fill()

	dc.DrawCircle(float64(dc.Height())/2, float64(dc.Width())/2, 20)
	dc.SetColor(color.White)
	dc.Fill()

	dc.SavePNG("out_gradient.png")
	fmt.Printf("Hello there.\n")
}

func sizer(filename string, size int) error {
	canvas, err := gg.LoadImage("./images/" + filename)
	if err != nil {
		return err
	}
	filter := gift.New(gift.ResizeToFit(size, size, gift.LanczosResampling))
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	filter.Draw(dst, canvas)

	f, err := os.Create(fmt.Sprintf("images/%v_youmu.png", size))
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, dst)
	if err != nil {
		fmt.Printf("Error encoding the png.")
		return err
	}
	return nil
}

func write_to_png(dst *image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, *dst)
	if err != nil {
		fmt.Printf("Error encoding the png.")
		return err
	}
	return nil
}

func mask(filename string) (*gg.Context, error) {
	canvas, err := gg.LoadImage("./images/" + filename)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContext(512, 512)
	filter := gift.New(
		gift.Resize(512, 0, gift.BoxResampling),
		gift.Contrast(-20),
	)
	// dc.DrawRoundedRectangle(0, 0, 512, 512, 64*2)
	dc.DrawCircle(float64(dc.Height()) / 2, float64(dc.Width()) / 2, 256);
	dc.Clip() // gets the clip from DrawRoundedRectangle??

	dst := image.NewRGBA(filter.Bounds((canvas).Bounds()))
	filter.Draw(dst, canvas)

	/*
		img := dst.SubImage(dst.Rect);
		write_to_png(&img, "dst.png");
	*/

	dc.DrawImage(dst, 0, 0)
	return dc, nil
}

func main() {
	files := []string{"100_youmu.png", "200_youmu.png", "300_youmu.png", "400_youmu.png", "500_youmu.png", "full_youmu.jpg"}
	for _, file := range files {
		dc, err := mask(file)
		if err != nil {
			fmt.Printf("%v couldn't be masked. %v", file, err)
			return
		}
		dc.SavePNG(fmt.Sprintf("./images/mask/mask_%v", file))
	}
}
