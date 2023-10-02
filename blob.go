package codec

import (
	"bytes"
	"encoding/binary"
	"io"
)

func SerializeBlobs(blobs [][]byte) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	for _, blob := range blobs {
		err := binary.Write(buffer, binary.BigEndian, int32(len(blob)))

		if err != nil {
			return nil, err
		}

		_, err = buffer.Write(blob)

		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func DeserializeBlobs(blobsData []byte) ([][]byte, error) {
	buffer := bytes.NewBuffer(blobsData)
	var err error
	var length int32
	var blobs [][]byte

	for {
		err = binary.Read(buffer, binary.BigEndian, &length)

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		data := make([]byte, length)
		_, err = buffer.Read(data)

		if err != nil {
			return nil, err
		}

		blobs = append(blobs, data)
	}

	return blobs, nil
}
