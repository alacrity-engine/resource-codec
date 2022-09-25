package codec

type PixFormat int

const (
	PixFormatRGBA PixFormat = iota
	PixFormatRGB
	PixFormatCMYK
)

func (pixFormat PixFormat) String() string {
	switch pixFormat {
	case PixFormatRGBA:
		return "RGBA"

	case PixFormatRGB:
		return "RGB"

	case PixFormatCMYK:
		return "CMYK"

	default:
		return ""
	}
}
