package codec

import (
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// hash is the map of all the available
	// hashing algorithms.
	hash = map[HashAlgorithm]func([]byte) ([]byte, error){
		HashAlgorithmKeccak256: HashKeccak256,
		HashAlgorithmSHA256:    HashSHA256,
	}
)

// Hash computes the hash sum of the
// given data using the consented
// hashing algorithm.
func Hash(in []byte) ([]byte, error) {
	hashAlgorithm, ok := hash[ConsentedHashAlgorithm]

	if !ok {
		return nil, fmt.Errorf(
			"hash algorithm '%s' not found",
			ConsentedHashAlgorithm)
	}

	return hashAlgorithm(in)
}

// HashKeccak256 is the Keccak256 hashing algorithm.
func HashKeccak256(in []byte) ([]byte, error) {
	hash := crypto.NewKeccakState()
	total := 0
	var written int
	var err error

	for written, err = hash.Write(in[total:]); written > 0 && err == nil; written, err = hash.Write(in[total:]) {
		total += written
	}

	if err != nil {
		return nil, err
	}

	sum := hash.Sum(nil)

	return sum, nil
}

// HashSHA256 is the SHA-256 hashing algorithm.
func HashSHA256(in []byte) ([]byte, error) {
	sum := sha256.Sum256(in)

	return sum[:], nil
}
