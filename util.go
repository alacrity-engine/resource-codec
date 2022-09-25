package codec

import "image"

// sliceEqual is a generic function for
// slice element-wise equality comparison.
func sliceEqual[T comparable](left, right []T) bool {
	if len(left) != len(right) {
		return false
	}

	for i, elem := range left {
		if elem != right[i] {
			return false
		}
	}

	return true
}

// reversePix reverses the pixel data
// of the image.
func reversePix(arr []byte) {
	start := 0
	end := len(arr) - 4

	for start < end {
		for i := 0; i < 4; i++ {
			temp := arr[start+i]
			arr[start+i] = arr[end+i]
			arr[end+i] = temp
		}

		start += 4
		end -= 4
	}
}

// mirror mirrors the image
// about the vertical axis.
func mirror(img *image.RGBA) {
	for i := 0; i < img.Rect.Dy(); i++ {
		reversePix(img.Pix[4*i*img.Rect.Dx() : 4*(i+1)*img.Rect.Dx()])
	}
}
