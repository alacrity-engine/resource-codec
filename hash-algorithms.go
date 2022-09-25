package codec

type HashAlgorithm int

const (
	HashAlgorithmKeccak256 HashAlgorithm = iota // HashAlgorithmKeccak256 is the identifier of the Keccak256 hashing algorithm.
	HashAlgorithmSHA256                         // HashAlgorithmSHA256 is the identifier of the SHA-256 hashing algorithm.
)

// String returns the name of the hashing algorithm.
func (alg HashAlgorithm) String() string {
	switch alg {
	case HashAlgorithmKeccak256:
		return "Keccak256"

	case HashAlgorithmSHA256:
		return "SHA-256"

	default:
		return ""
	}
}
