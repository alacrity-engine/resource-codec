package codec

import (
	"bytes"
	"encoding/binary"
)

type SpritesheetData struct {
	Width  int32
	Height int32
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

	return ssdata, nil
}
