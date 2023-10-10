package codec

import (
	"bytes"
	"encoding/binary"
)

type ListHeader struct {
	Count int32
}

func (lh *ListHeader) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	err := binary.Write(buffer, binary.BigEndian, lh.Count)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ListHeaderFromBytes(data []byte) (*ListHeader, error) {
	buffer := bytes.NewBuffer(data)
	lh := &ListHeader{}

	err := binary.Read(buffer, binary.BigEndian, &lh.Count)

	if err != nil {
		return nil, err
	}

	return lh, nil
}
