package codec

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
)

type Picture struct {
	Width         int32
	Height        int32
	Pix           []byte
	Hash          []byte
	PixFormat     PixFormat
	HashAlgorithm HashAlgorithm
}

type CompressedPicture struct {
	Width                 int32
	Height                int32
	CompressedPix         []byte
	OriginalHash          []byte
	OriginalPixFormat     PixFormat
	OriginalHashAlgorithm HashAlgorithm
	CompressionAlgorithm  CompressionAlgorithm
}

func (picture *Picture) Compress() (*CompressedPicture, error) {
	compressedPix, err := Compress(picture.Pix)

	if err != nil {
		return nil, err
	}

	return &CompressedPicture{
		Width:                 picture.Width,
		Height:                picture.Height,
		CompressedPix:         compressedPix,
		OriginalHash:          picture.Hash,
		OriginalPixFormat:     picture.PixFormat,
		OriginalHashAlgorithm: picture.HashAlgorithm,
		CompressionAlgorithm:  ConsentedCompressionAlgorithm,
	}, nil
}

func NewPictureFromImage(img image.Image) (*Picture, error) {
	imgRGBA := image.NewRGBA(image.Rect(0, 0,
		img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(),
		img.Bounds(), img.Bounds().Min, draw.Src)

	hash, err := Hash(imgRGBA.Pix)

	if err != nil {
		return nil, err
	}

	return &Picture{
		Width:         int32(imgRGBA.Bounds().Dx()),
		Height:        int32(imgRGBA.Bounds().Dy()),
		Pix:           imgRGBA.Pix,
		Hash:          hash,
		PixFormat:     ConsentedPixFormat,
		HashAlgorithm: ConsentedHashAlgorithm,
	}, nil
}
