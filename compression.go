package codec

import (
	"bytes"
	"compress/lzw"
	"fmt"
)

var (
	// compression is the map containing all the
	// available compression functions by their identifiers.
	compression = map[CompressionAlgorithm]func([]byte) ([]byte, error){
		CompressionAlgorithmLZWOrderLSBLitWidth8: CompressLZWOrderLSBLitWidth8,
	}
)

// Compress compresses the given data using
// the consented compression function.
func Compress(in []byte) ([]byte, error) {
	compressionAlgorithm, ok := compression[CompressionAlgorithmLZWOrderLSBLitWidth8]

	if !ok {
		return nil, fmt.Errorf(
			"compression algorithm '%s' not found",
			ConsentedCompressionAlgorithm)
	}

	return compressionAlgorithm(in)
}

// CompressLZWOrderLSBLitWidth8 uses the
// LZWOrderLSBLitWidth8 compression algorithm
// to compress the given data.
func CompressLZWOrderLSBLitWidth8(in []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := lzw.NewWriter(&buffer, lzw.LSB, 8)
	total := 0
	var written int
	var err error

	for written, err = writer.Write(in[total:]); written > 0 && err == nil; written, err = writer.Write(in[total:]) {
		total += written
	}

	if err != nil {
		return nil, err
	}

	// It seems it's not okay to
	// defer the Close() call here.
	err = writer.Close()

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
