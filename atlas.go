package codec

import (
	"bytes"
	"encoding/binary"

	"github.com/alacrity-engine/core/math/geometry"
)

type AtlasData struct {
	Frames    map[rune]geometry.Rect
	SymbolSet *PictureData
}

type CompressedFrames struct {
	FrameCount     int32
	OrigDataLength int32
	Data           []byte
}

type CompressedAtlasData struct {
	CompressedFramesData *CompressedFrames
	CompressedSymbolSet  *CompressedPictureData
}

func framesDictToBytes(frames map[rune]geometry.Rect) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	for symbol, rect := range frames {
		err := binary.Write(buffer, binary.BigEndian, symbol)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, rect.Min.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, rect.Min.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, rect.Max.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, rect.Max.Y)

		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func framesDictFromBytes(data []byte, frameCount int32) (map[rune]geometry.Rect, error) {
	frames := make(map[rune]geometry.Rect, frameCount)
	buffer := bytes.NewBuffer(data)

	for i := 0; i < int(frameCount); i++ {
		var symbol rune

		err := binary.Read(buffer, binary.BigEndian, &symbol)

		if err != nil {
			return nil, err
		}

		var rect geometry.Rect

		err = binary.Read(buffer, binary.BigEndian, &rect.Min.X)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &rect.Min.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &rect.Max.X)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &rect.Max.Y)

		if err != nil {
			return nil, err
		}

		frames[symbol] = rect
	}

	return frames, nil
}

func (ad *AtlasData) Compress() (*CompressedAtlasData, error) {
	framesData, err := framesDictToBytes(ad.Frames)

	if err != nil {
		return nil, err
	}

	compressedFramesData, err := Compress(framesData)

	if err != nil {
		return nil, err
	}

	compressedFrames := &CompressedFrames{
		FrameCount:     int32(len(ad.Frames)),
		OrigDataLength: int32(len(framesData)),
		Data:           compressedFramesData,
	}

	compressedSymbolSet, err := ad.SymbolSet.Compress()

	if err != nil {
		return nil, err
	}

	return &CompressedAtlasData{
		CompressedFramesData: compressedFrames,
		CompressedSymbolSet:  compressedSymbolSet,
	}, nil
}

func (cad *CompressedAtlasData) Decompress() (*AtlasData, error) {
	framesData, err := Decompress(
		cad.CompressedFramesData.Data,
		int(cad.CompressedFramesData.OrigDataLength))

	if err != nil {
		return nil, err
	}

	frames, err := framesDictFromBytes(framesData,
		cad.CompressedFramesData.FrameCount)

	if err != nil {
		return nil, err
	}

	picture, err := cad.CompressedSymbolSet.Decompress()

	if err != nil {
		return nil, err
	}

	return &AtlasData{
		Frames:    frames,
		SymbolSet: picture,
	}, nil
}

func (cf *CompressedFrames) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, cf.FrameCount)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, cf.OrigDataLength)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(cf.Data)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write(cf.Data)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CompressedFramesFromBytes(data []byte) (*CompressedFrames, error) {
	buffer := bytes.NewBuffer(data)
	cf := &CompressedFrames{}

	err := binary.Read(buffer, binary.BigEndian, &cf.FrameCount)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &cf.OrigDataLength)

	if err != nil {
		return nil, err
	}

	var length int32
	err = binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	cf.Data = make([]byte, length)
	_, err = buffer.Read(cf.Data)

	if err != nil {
		return nil, err
	}

	return cf, nil
}

func (cad *CompressedAtlasData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	// Write the character index.
	data, err := cad.CompressedFramesData.ToBytes()

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(data)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write(data)

	if err != nil {
		return nil, err
	}

	// Write the picture.
	data, err = cad.CompressedSymbolSet.ToBytes()

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(data)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write(data)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CompressedAtlasDataFromBytes(data []byte) (*CompressedAtlasData, error) {
	buffer := bytes.NewBuffer(data)
	cad := &CompressedAtlasData{}

	// Read the compressed frames data,
	var length int32
	err := binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	fieldData := make([]byte, length)
	_, err = buffer.Read(fieldData)

	if err != nil {
		return nil, err
	}

	cad.CompressedFramesData, err = CompressedFramesFromBytes(fieldData)

	if err != nil {
		return nil, err
	}

	// Read the font picture.
	err = binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	fieldData = make([]byte, length)
	_, err = buffer.Read(fieldData)

	if err != nil {
		return nil, err
	}

	cad.CompressedSymbolSet, err = CompressedPictureFromBytes(fieldData)

	if err != nil {
		return nil, err
	}

	return cad, nil
}
