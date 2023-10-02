package codec

import (
	"bytes"
	"encoding/binary"
)

type BatchData struct {
	Name      string
	CanvasID  string
	TextureID string
	ZMin      float32
	ZMax      float32
}

func (bdata *BatchData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, int32(len(bdata.Name)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(bdata.Name))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(bdata.CanvasID)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(bdata.CanvasID))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, int32(len(bdata.TextureID)))

	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(bdata.TextureID))

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, bdata.ZMin)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, bdata.ZMax)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func BatchDataFromBytes(data []byte) (*BatchData, error) {
	buffer := bytes.NewBuffer(data)
	batchData := &BatchData{}
	var length int32

	err := binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	strData := make([]byte, length)
	_, err = buffer.Read(strData)

	if err != nil {
		return nil, err
	}

	batchData.Name = string(strData)

	err = binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	strData = make([]byte, length)
	_, err = buffer.Read(strData)

	if err != nil {
		return nil, err
	}

	batchData.CanvasID = string(strData)

	err = binary.Read(buffer, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	strData = make([]byte, length)
	_, err = buffer.Read(strData)

	if err != nil {
		return nil, err
	}

	batchData.TextureID = string(strData)

	err = binary.Read(buffer, binary.BigEndian, &batchData.ZMin)

	if err != nil {
		return nil, err
	}

	err = binary.Read(buffer, binary.BigEndian, &batchData.ZMax)

	if err != nil {
		return nil, err
	}

	return batchData, nil
}
