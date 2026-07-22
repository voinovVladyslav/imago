// Package transformer does image transformations (grayscale and more)
package transformer

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func invertedFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
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

func grayscaleFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}
	return newImg
}

func greenlessFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, _, b, _ := img.At(x, y).RGBA()
			red := uint8(r >> 8)
			blue := uint8(b >> 8)
			color := color.RGBA{red, 0, blue, 255}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}

func bluelessFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, _, _ := img.At(x, y).RGBA()
			red := uint8(r >> 8)
			green := uint8(g >> 8)
			color := color.RGBA{red, green, 0, 255}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}

func redlessFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, g, b, _ := img.At(x, y).RGBA()
			green := uint8(g >> 8)
			blue := uint8(b >> 8)
			color := color.RGBA{0, green, blue, 255}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}

func redFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			red := uint8(r >> 8)
			color := color.RGBA{red, 0, 0, 255}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}
func greenFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, g, _, _ := img.At(x, y).RGBA()
			green := uint8(g >> 8)
			color := color.RGBA{0, green, 0, 255}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}
func blueFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, b, _ := img.At(x, y).RGBA()
			blue := uint8(b >> 8)
			color := color.RGBA{0, 0, blue, 255}
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
	newImg := redFilter(img)

	result, err := os.Create("result.jpg")
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	defer result.Close()

	err = jpeg.Encode(result, newImg, &jpeg.Options{Quality: 100})
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}
