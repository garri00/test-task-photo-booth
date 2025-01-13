package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

func ResizeImageB64(dataB64, extension string, percentage uint) (string, error) {
	var resizedImageB64 string

	b, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		return "", fmt.Errorf("base64.StdEncoding.DecodeString() failed: %w", err)
	}

	switch extension {
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(b))
		if err != nil {
			return "", fmt.Errorf("jpeg.Decode() failed: %w", err)
		}

		width, height := getResizedImageBounds(img, percentage)
		resImag := resize.Resize(width, height, img, resize.Lanczos3)

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, resImag, &jpeg.Options{Quality: 20}); err != nil {
			return "", fmt.Errorf("jpeg.Decode() failed: %w", err)
		}

		resizedImageB64 = base64.StdEncoding.EncodeToString(buf.Bytes())

	case "image/png":
		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			return "", fmt.Errorf("jpeg.Decode() failed: %w", err)
		}

		width, height := getResizedImageBounds(img, percentage)
		resImag := resize.Resize(width, height, img, resize.Lanczos3)

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, resImag, &jpeg.Options{Quality: 20}); err != nil {
			return "", fmt.Errorf("jpeg.Encode() failed: %w", err)
		}

		resizedImageB64 = base64.StdEncoding.EncodeToString(buf.Bytes())

	default:
		return "", fmt.Errorf("unknown extension")
	}

	return resizedImageB64, nil
}

func getResizedImageBounds(img image.Image, percentage uint) (uint, uint) {
	width := uint(float64(img.Bounds().Dx()) * float64(percentage) / 100)
	height := uint(float64(img.Bounds().Dy()) * float64(percentage) / 100)

	return width, height
}
