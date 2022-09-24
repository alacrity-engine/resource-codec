package codec

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"io"
)

var (
	decompression = map[CompressionAlgorithm]func([]byte, int) ([]byte, error){
		CompressionAlgorithmLZWOrderLSBLitWidth8: DecompressLZWOrderLSBLitWidth8,
	}
)

func Decompress(in []byte, sourceSize int) ([]byte, error) {
	decompressionAlgorithm, ok := decompression[ConsentedCompressionAlgorithm]

	if !ok {
		return nil, fmt.Errorf(
			"decompression algorithm '%s' not found",
			ConsentedCompressionAlgorithm)
	}

	return decompressionAlgorithm(in, sourceSize)
}

func DecompressLZWOrderLSBLitWidth8(in []byte, sourceSize int) ([]byte, error) {
	buffer := bytes.NewReader(in)
	source := make([]byte, sourceSize)

	reader := lzw.NewReader(buffer, lzw.LSB, 8)
	defer reader.Close()
	total := 0
	read := 1
	var err error

	for err != io.EOF && read > 0 {
		read, err = reader.Read(source[total:])
		total += read
	}

	return source, nil
}
