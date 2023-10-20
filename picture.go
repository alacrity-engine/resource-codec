package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/alacrity-engine/core/math/geometry"
)

type PictureData struct {
	Width         int32
	Height        int32
	Pix           []byte
	Hash          []byte
	PixFormat     PixFormat
	HashAlgorithm HashAlgorithm
}

type CompressedPictureData struct {
	Width                 int32
	Height                int32
	OriginalPixSize       int32
	CompressedPix         []byte
	OriginalHash          []byte
	OriginalPixFormat     PixFormat
	OriginalHashAlgorithm HashAlgorithm
	CompressionAlgorithm  CompressionAlgorithm
}

// GetSpritesheetFrames returns the set of rectangles
// corresponding to the frames of the spritesheet.
func (pic *PictureData) GetSpritesheetFrames(ss *SpritesheetData) ([]geometry.Rect, error) {
	orig := OrigData{
		X: ss.Orig.X,
		Y: ss.Orig.Y + ss.Area.PixelHeight,
	}

	if orig.X+ss.Area.PixelWidth > pic.Width ||
		orig.Y-ss.Area.PixelHeight < 0 {
		return nil, fmt.Errorf(
			"the spritesheet cannot be applied to the picture")
	}

	frames := make([]geometry.Rect, 0)
	pixelWidth := float64(ss.Area.PixelWidth)
	pixelHeight := float64(ss.Area.PixelHeight)
	dw := pixelWidth / float64(ss.Width)
	dh := pixelHeight / float64(ss.Height)

	for y := float64(orig.Y); y > float64(ss.Orig.Y); y -= dh {
		for x := float64(orig.X); x < float64(orig.X)+pixelWidth; x += dw {
			frame := geometry.R(x, y-dh, x+dw, y)
			frames = append(frames, frame)
		}
	}

	return frames, nil
}

// GetSpritesheetFrames returns the set of rectangles
// corresponding to the frames of the spritesheet.
func (cpic *CompressedPictureData) GetSpritesheetFrames(ss *SpritesheetData) ([]geometry.Rect, error) {
	orig := OrigData{
		X: ss.Orig.X,
		Y: ss.Orig.Y + ss.Area.PixelHeight,
	}

	if orig.X+ss.Area.PixelWidth > cpic.Width ||
		orig.Y-ss.Area.PixelHeight < 0 {
		return nil, fmt.Errorf(
			"the spritesheet cannot be applied to the picture")
	}

	frames := make([]geometry.Rect, 0)
	pixelWidth := float64(ss.Area.PixelWidth)
	pixelHeight := float64(ss.Area.PixelHeight)
	dw := pixelWidth / float64(ss.Width)
	dh := pixelHeight / float64(ss.Height)

	for y := float64(orig.Y); y > float64(ss.Orig.Y); y -= dh {
		for x := float64(orig.X); x < float64(orig.X)+pixelWidth; x += dw {
			frame := geometry.R(x, y-dh, x+dw, y)
			frames = append(frames, frame)
		}
	}

	return frames, nil
}

func (picture *PictureData) Compress() (*CompressedPictureData, error) {
	compressedPix, err := Compress(picture.Pix)

	if err != nil {
		return nil, err
	}

	return &CompressedPictureData{
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

func (compressedPicture *CompressedPictureData) Decompress() (*PictureData, error) {
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

	if !sliceEqual(decompressedHash, compressedPicture.OriginalHash) {
		return nil, fmt.Errorf(
			"data corruption error: expected %s but got %s (%s)",
			hex.EncodeToString(compressedPicture.OriginalHash),
			hex.EncodeToString(decompressedHash),
			compressedPicture.OriginalHashAlgorithm)
	}

	return &PictureData{
		Width:         compressedPicture.Width,
		Height:        compressedPicture.Height,
		Pix:           decompressedPix,
		Hash:          decompressedHash,
		PixFormat:     compressedPicture.OriginalPixFormat,
		HashAlgorithm: compressedPicture.OriginalHashAlgorithm,
	}, nil
}

func (compressedPicture *CompressedPictureData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, compressedPicture.Width)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, compressedPicture.Height)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, compressedPicture.OriginalPixSize)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(compressedPicture.CompressedPix)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write(compressedPicture.CompressedPix)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(compressedPicture.OriginalHash)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write(compressedPicture.OriginalHash)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(compressedPicture.OriginalPixFormat))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(compressedPicture.OriginalHashAlgorithm))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(compressedPicture.CompressionAlgorithm))

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CompressedPictureFromBytes(data []byte) (*CompressedPictureData, error) {
	buffer := bytes.NewBuffer(data)
	compressedPicture := &CompressedPictureData{}

	err := binary.Read(buffer, binary.BigEndian, &compressedPicture.Width)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &compressedPicture.Height)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &compressedPicture.OriginalPixSize)

	if err != nil {
		return nil, err
	}

	var compressedPixLength int32
	err = binary.Read(buffer, binary.BigEndian, &compressedPixLength)

	if err != nil {
		return nil, err
	}

	compressedPix := make([]byte, compressedPixLength)
	_, err = buffer.Read(compressedPix)

	if err != nil {
		return nil, err
	}

	compressedPicture.CompressedPix = compressedPix

	var originalHashLength int32
	err = binary.Read(buffer, binary.BigEndian, &originalHashLength)

	if err != nil {
		return nil, err
	}

	originalHash := make([]byte, originalHashLength)
	_, err = buffer.Read(originalHash)

	if err != nil {
		return nil, err
	}

	compressedPicture.OriginalHash = originalHash

	var originalPixFormat int32
	err = binary.Read(buffer, binary.BigEndian, &originalPixFormat)

	if err != nil {
		return nil, err
	}

	compressedPicture.OriginalPixFormat = PixFormat(originalPixFormat)

	var originalHashAlgorithm int32
	err = binary.Read(buffer, binary.BigEndian, &originalHashAlgorithm)

	if err != nil {
		return nil, err
	}

	compressedPicture.OriginalHashAlgorithm = HashAlgorithm(originalHashAlgorithm)

	var compressionAlgorithm int32
	err = binary.Read(buffer, binary.BigEndian, &compressionAlgorithm)

	if err != nil {
		return nil, err
	}

	compressedPicture.CompressionAlgorithm = CompressionAlgorithm(compressionAlgorithm)

	return compressedPicture, nil
}

func NewPictureFromImage(img image.Image) (*PictureData, error) {
	imgRGBA := image.NewRGBA(image.Rect(0, 0,
		img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(),
		img.Bounds(), img.Bounds().Min, draw.Src)
	reversePix(imgRGBA.Pix)
	mirror(imgRGBA)

	hash, err := Hash(imgRGBA.Pix)

	if err != nil {
		return nil, err
	}

	return &PictureData{
		Width:         int32(imgRGBA.Bounds().Dx()),
		Height:        int32(imgRGBA.Bounds().Dy()),
		Pix:           imgRGBA.Pix,
		Hash:          hash,
		PixFormat:     ConsentedPixFormat,
		HashAlgorithm: ConsentedHashAlgorithm,
	}, nil
}
