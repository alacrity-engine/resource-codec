package codec_test

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"testing"

	codec "github.com/alacrity-engine/resource-codec"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed test.jpg
	TestImage []byte
)

func TestCompressLZWOrderLSBLitWidth8(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(TestImage))
	assert.Nil(t, err)
	imgRGBA := image.NewRGBA(image.Rect(0, 0,
		img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(),
		img.Bounds(), img.Bounds().Min, draw.Src)

	compressedPix, err := codec.CompressLZWOrderLSBLitWidth8(imgRGBA.Pix)
	assert.Nil(t, err)
	assert.Equal(t, 4096, len(compressedPix))
}

func TestDecompressLZWOrderLSBLitWidth8(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(TestImage))
	assert.Nil(t, err)
	imgRGBA := image.NewRGBA(image.Rect(0, 0,
		img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(),
		img.Bounds(), img.Bounds().Min, draw.Src)

	compressedPix, err := codec.CompressLZWOrderLSBLitWidth8(imgRGBA.Pix)
	assert.Nil(t, err)
	assert.Equal(t, 4096, len(compressedPix))

	decompressedPix, err := codec.DecompressLZWOrderLSBLitWidth8(compressedPix,
		imgRGBA.Rect.Dx()*imgRGBA.Rect.Dy()*4)
	assert.Nil(t, err)
	assert.Equal(t, imgRGBA.Rect.Dx()*imgRGBA.Rect.Dy()*4, len(decompressedPix))

	hashOriginal, err := codec.Hash(imgRGBA.Pix)
	assert.Nil(t, err)
	hashDecompressed, err := codec.Hash(decompressedPix)
	assert.Nil(t, err)
	assert.ElementsMatch(t, hashOriginal, hashDecompressed)
}
