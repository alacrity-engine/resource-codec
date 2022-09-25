package codec

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
)

// TODO: add the CBOR serialization
// to save pictures to .res files.

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
	OriginalPixSize       int32
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
		OriginalPixSize:       int32(len(picture.Pix)),
		CompressedPix:         compressedPix,
		OriginalHash:          picture.Hash,
		OriginalPixFormat:     picture.PixFormat,
		OriginalHashAlgorithm: picture.HashAlgorithm,
		CompressionAlgorithm:  ConsentedCompressionAlgorithm,
	}, nil
}

func (compressedPicture *CompressedPicture) Decompress() (*Picture, error) {
	previousCompressionAlgorithm := ConsentedCompressionAlgorithm
	defer func() { ConsentedCompressionAlgorithm = previousCompressionAlgorithm }()
	ConsentedCompressionAlgorithm = compressedPicture.CompressionAlgorithm

	previousHashAlgorithm := ConsentedHashAlgorithm
	defer func() { ConsentedHashAlgorithm = previousHashAlgorithm }()
	ConsentedHashAlgorithm = compressedPicture.OriginalHashAlgorithm

	decompressedPix, err := Decompress(compressedPicture.CompressedPix,
		int(compressedPicture.OriginalPixSize))

	if err != nil {
		return nil, err
	}

	decompressedHash, err := Hash(decompressedPix)

	if err != nil {
		return nil, err
	}

	if !SliceEqual(decompressedHash, compressedPicture.OriginalHash) {
		return nil, fmt.Errorf(
			"data corruption error: expected %s but got %s (%s)",
			hex.EncodeToString(compressedPicture.OriginalHash),
			hex.EncodeToString(decompressedHash),
			compressedPicture.OriginalHashAlgorithm)
	}

	return &Picture{
		Width:         compressedPicture.Width,
		Height:        compressedPicture.Height,
		Pix:           decompressedPix,
		Hash:          decompressedHash,
		PixFormat:     compressedPicture.OriginalPixFormat,
		HashAlgorithm: compressedPicture.OriginalHashAlgorithm,
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
