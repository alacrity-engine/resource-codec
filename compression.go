package codec

import (
	"bytes"
	"compress/lzw"
	"fmt"
)

var (
	compression = map[CompressionAlgorithm]func([]byte) ([]byte, error){
		CompressionAlgorithmLZWOrderLSBLitWidth8: CompressLZWOrderLSBLitWidth8,
	}
)

func Compress(in []byte) ([]byte, error) {
	compressionAlgorithm, ok := compression[CompressionAlgorithmLZWOrderLSBLitWidth8]

	if !ok {
		return nil, fmt.Errorf(
			"compression algorithm '%s' not found",
			ConsentedCompressionAlgorithm)
	}

	return compressionAlgorithm(in)
}

func CompressLZWOrderLSBLitWidth8(in []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := lzw.NewWriter(&buffer, lzw.LSB, 8)
	defer writer.Close()
	total := 0
	var written int
	var err error

	for written, err = writer.Write(in[total:]); written > 0 && err == nil; written, err = writer.Write(in[total:]) {
		total += written
	}

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
