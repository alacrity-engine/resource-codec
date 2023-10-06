package codec

import (
	"bytes"
	"encoding/binary"
)

type CanvasData struct {
	Name  string
	DrawZ int32
}

func (cdata *CanvasData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, int32(len(cdata.Name)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(cdata.Name))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, cdata.DrawZ)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CanvasDataFromBytes(data []byte) (*CanvasData, error) {
	buffer := bytes.NewBuffer(data)
	cdata := &CanvasData{}

	var length int32
	err := binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	nameData := make([]byte, length)
	_, err = buffer.Read(nameData)

	if err != nil {
		return nil, err
	}

	cdata.Name = string(nameData)

	err = binary.Read(buffer, binary.BigEndian, &cdata.DrawZ)

	if err != nil {
		return nil, err
	}

	return cdata, nil
}
