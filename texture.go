package codec

import (
	"bytes"
	"encoding/binary"
)

type TextureData struct {
	Name      string
	PictureID string
	Filtering uint32
}

func (tdata *TextureData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, len(tdata.Name))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(tdata.Name))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, len(tdata.PictureID))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(tdata.PictureID))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, tdata.Filtering)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func TextureDataFromBytes(data []byte) (*TextureData, error) {
	buffer := bytes.NewBuffer(data)
	tdata := &TextureData{}

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

	tdata.Name = string(nameData)

	err = binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	pictureIDData := make([]byte, length)
	_, err = buffer.Read(pictureIDData)

	if err != nil {
		return nil, err
	}

	tdata.PictureID = string(pictureIDData)

	err = binary.Read(buffer, binary.BigEndian, &tdata.Filtering)

	if err != nil {
		return nil, err
	}

	return tdata, nil
}
