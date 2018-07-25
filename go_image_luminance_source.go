package gozxing

import (
	"image"
)

func NewLuminanceSourceFromImage(img image.Image) LuminanceSource {
	rect := img.Bounds()
	top := rect.Min.Y
	left := rect.Min.X
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y

	luminance := make([]byte, width*height)
	for y := 0; y < height; y++ {
		offset := y * width
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			luminance[offset+x] = byte((r + 2*g + b) * 255 / (4 * 0xffff))
		}
	}

	return &RGBLuminanceSource{
		LuminanceSourceBase{width, height},
		luminance,
		width,
		height,
		top,
		left,
	}
}
