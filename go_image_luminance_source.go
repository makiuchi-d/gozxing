package gozxing

import (
	"image"
)

func NewLuminanceSourceFromImage(img image.Image) LuminanceSource {
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y

	luminance := make([]byte, width*height)
	index := 0
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			luminance[index] = byte((r + 2*g + b) * 255 / (4 * 0xffff))
			index++
		}
	}

	return &RGBLuminanceSource{
		LuminanceSourceBase{width, height},
		luminance,
		width,
		height,
		0,
		0,
	}
}
