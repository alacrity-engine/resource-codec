package codec

type CompressionAlgorithm int

const (
	CompressionAlgorithmLZWOrderLSBLitWidth8 CompressionAlgorithm = iota // CompressionAlgorithmLZWOrderLSBLitWidth8 stands for the LZW compression algorithm with LSB for order and 8 for literal width.
)

// String returns the name of the compression algorithm.
func (alg CompressionAlgorithm) String() string {
	switch alg {
	case CompressionAlgorithmLZWOrderLSBLitWidth8:
		return "LZW-LSB-8"

	default:
		return ""
	}
}
