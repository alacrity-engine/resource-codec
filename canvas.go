package codec

import (
	"bytes"
	"encoding/binary"

	"github.com/go-gl/mathgl/mgl32"
)

type CanvasData struct {
	Name       string
	DrawZ      int
	Projection mgl32.Mat4
}

func (cdata *CanvasData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, len(cdata.Name))

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

	err = binary.Write(buffer, binary.BigEndian, cdata.Projection)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func CanvasDataFromBytes(data []byte) (*CanvasData, error) {
	buffer := bytes.NewBuffer(data)
	cdata := &CanvasData{}

	var length int
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

	err = binary.Read(buffer, binary.BigEndian, &cdata.Projection)

	if err != nil {
		return nil, err
	}

	return cdata, nil
}
