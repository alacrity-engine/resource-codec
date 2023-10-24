package codec

import (
	"bytes"
	"encoding/binary"

	"github.com/alacrity-engine/core/math/geometry"
)

type AtlasData struct {
	Glyphs    map[rune]GlyphData
	SymbolSet *PictureData
	Size      int32
	MaxHeight float64
	FontName  string
}

type GlyphData struct {
	Dot     geometry.Vec
	Frame   geometry.Rect
	Advance float64
}

type CompressedFrames struct {
	FrameCount     int32
	OrigDataLength int32
	Data           []byte
}

type CompressedAtlasData struct {
	CompressedFramesData *CompressedFrames
	CompressedSymbolSet  *CompressedPictureData
	Size                 int32
	MaxHeight            float64
	FontName             string
}

func glyphsDictToBytes(glyphs map[rune]GlyphData) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	for symbol, glyph := range glyphs {
		err := binary.Write(buffer, binary.BigEndian, symbol)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Dot.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Dot.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Frame.Min.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Frame.Min.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Frame.Max.X)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Frame.Max.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Write(buffer, binary.BigEndian, glyph.Advance)

		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func glyphsDictFromBytes(data []byte, frameCount int32) (map[rune]GlyphData, error) {
	frames := make(map[rune]GlyphData, frameCount)
	buffer := bytes.NewBuffer(data)

	for i := 0; i < int(frameCount); i++ {
		var symbol rune

		err := binary.Read(buffer, binary.BigEndian, &symbol)

		if err != nil {
			return nil, err
		}

		var dot geometry.Vec

		err = binary.Read(buffer, binary.BigEndian, &dot.X)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &dot.Y)

		if err != nil {
			return nil, err
		}

		var frame geometry.Rect

		err = binary.Read(buffer, binary.BigEndian, &frame.Min.X)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &frame.Min.Y)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &frame.Max.X)

		if err != nil {
			return nil, err
		}

		err = binary.Read(buffer, binary.BigEndian, &frame.Max.Y)

		if err != nil {
			return nil, err
		}

		var advance float64

		err = binary.Read(buffer, binary.BigEndian, &advance)

		if err != nil {
			return nil, err
		}

		frames[symbol] = GlyphData{
			Dot:     dot,
			Frame:   frame,
			Advance: advance,
		}
	}

	return frames, nil
}

func (ad *AtlasData) Compress() (*CompressedAtlasData, error) {
	framesData, err := glyphsDictToBytes(ad.Glyphs)

	if err != nil {
		return nil, err
	}

	compressedFramesData, err := Compress(framesData)

	if err != nil {
		return nil, err
	}

	compressedFrames := &CompressedFrames{
		FrameCount:     int32(len(ad.Glyphs)),
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
		Size:                 ad.Size,
		MaxHeight:            ad.MaxHeight,
		FontName:             ad.FontName,
	}, nil
}

func (cad *CompressedAtlasData) Decompress() (*AtlasData, error) {
	framesData, err := Decompress(
		cad.CompressedFramesData.Data,
		int(cad.CompressedFramesData.OrigDataLength))

	if err != nil {
		return nil, err
	}

	frames, err := glyphsDictFromBytes(framesData,
		cad.CompressedFramesData.FrameCount)

	if err != nil {
		return nil, err
	}

	picture, err := cad.CompressedSymbolSet.Decompress()

	if err != nil {
		return nil, err
	}

	return &AtlasData{
		Glyphs:    frames,
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

	// Write everything else.
	err = binary.Write(buffer, binary.BigEndian, cad.Size)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, cad.MaxHeight)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(cad.FontName)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(cad.FontName)

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

	// Read everything else.
	err = binary.Read(buffer, binary.BigEndian, &cad.Size)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &cad.MaxHeight)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	fieldData = make([]byte, length)
	_, err = buffer.Read(fieldData)

	if err != nil {
		return nil, err
	}

	fontName := string(fieldData)
	cad.FontName = fontName

	return cad, nil
}
