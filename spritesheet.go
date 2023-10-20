package codec

import (
	"bytes"
	"encoding/binary"
)

type SpritesheetData struct {
	Width  int32
	Height int32
	Orig   OrigData
	Area   AreaData
}

type OrigData struct {
	X int32
	Y int32
}

type AreaData struct {
	PixelWidth  int32
	PixelHeight int32
}

func (ssdata *SpritesheetData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, ssdata.Width)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, ssdata.Height)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, ssdata.Orig.X)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, ssdata.Orig.Y)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, ssdata.Area.PixelWidth)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, ssdata.Area.PixelHeight)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func SpritesheetDataFromBytes(data []byte) (*SpritesheetData, error) {
	buffer := bytes.NewBuffer(data)
	ssdata := &SpritesheetData{}

	err := binary.Read(buffer, binary.BigEndian, &ssdata.Width)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &ssdata.Height)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &ssdata.Orig.X)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &ssdata.Orig.Y)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &ssdata.Area.PixelWidth)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &ssdata.Area.PixelHeight)

	if err != nil {
		return nil, err
	}

	return ssdata, nil
}
