package codec

var (
	ConsentedHashAlgorithm        = HashAlgorithmKeccak256                   // ConsentedHashAlgorithm is the hash algorithm chosen for overall integrity checking.
	ConsentedPixFormat            = PixFormatRGBA                            // ConsentedPixFormat is the pixel data format chosen for retrieved pictures.
	ConsentedCompressionAlgorithm = CompressionAlgorithmLZWOrderLSBLitWidth8 // ConsentedCompressionAlgorithm is the compression algorithm chosen for picture compressing.
)
