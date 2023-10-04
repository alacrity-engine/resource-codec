package codec_test

import (
	"bytes"
	"image"
	"testing"

	"github.com/alacrity-engine/core/math/geometry"
	codec "github.com/alacrity-engine/resource-codec"
	"github.com/stretchr/testify/assert"
)

func TestSerializePicture(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(TestImage))
	assert.Nil(t, err)
	picture, err := codec.NewPictureFromImage(img)
	assert.Nil(t, err)
	compressedPicture, err := picture.Compress()
	assert.Nil(t, err)

	data, err := compressedPicture.ToBytes()
	assert.Nil(t, err)
	assert.Equal(t, 7219, len(data))
}

func TestDeserializePicture(t *testing.T) {
	img, _, err := image.Decode(bytes.NewReader(TestImage))
	assert.Nil(t, err)
	picture, err := codec.NewPictureFromImage(img)
	assert.Nil(t, err)
	compressedPicture, err := picture.Compress()
	assert.Nil(t, err)

	data, err := compressedPicture.ToBytes()
	assert.Nil(t, err)
	assert.Equal(t, 7219, len(data))

	deserializedPicture, err := codec.CompressedPictureFromBytes(data)
	assert.Nil(t, err)
	assert.Equal(t, 7155, len(deserializedPicture.CompressedPix))
	assert.Equal(t, codec.HashAlgorithmKeccak256, deserializedPicture.OriginalHashAlgorithm)
	assert.Equal(t, codec.CompressionAlgorithmLZWOrderLSBLitWidth8, deserializedPicture.CompressionAlgorithm)
}

func TestSerializeAnimationData(t *testing.T) {
	animData := &codec.AnimationData{
		TextureID: "cirno-player",
		Frames: []geometry.Rect{
			geometry.R(0, 0, 32, 32),
			geometry.R(32, 0, 64, 32),
		},
		Durations: []int32{60, 120},
	}

	data, err := animData.ToBytes()
	assert.Nil(t, err)
	assert.Equal(t, 92, len(data))
}

func TestDeserializeAnimationData(t *testing.T) {
	animData := &codec.AnimationData{
		TextureID: "cirno-player",
		Frames: []geometry.Rect{
			geometry.R(0, 0, 32, 32),
			geometry.R(32, 0, 64, 32),
		},
		Durations: []int32{60, 120},
	}

	data, err := animData.ToBytes()
	assert.Nil(t, err)
	assert.Equal(t, 92, len(data))

	restoredAnimData, err := codec.AnimationDataFromBytes(data)
	assert.Nil(t, err)
	assert.Equal(t, animData.TextureID, restoredAnimData.TextureID)
	assert.ElementsMatch(t, animData.Frames, restoredAnimData.Frames)
	assert.ElementsMatch(t, animData.Durations, restoredAnimData.Durations)
}
