// Package transformer does image transformations (grayscale and more)
package transformer

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func InvertedFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			red := uint8(r >> 8)
			green := uint8(g >> 8)
			blue := uint8(b >> 8)
			color := color.RGBA{255 - red, 255 - green, 255 - blue, 255}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}

func Run() error {
	file, err := os.Open("./example.jpg")
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	newImg := InvertedFilter(img)

	result, err := os.Create("result.jpg")
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer result.Close()

	err = jpeg.Encode(result, newImg, &jpeg.Options{Quality: 90})
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}
