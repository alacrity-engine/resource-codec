package codec_test

import (
	"testing"

	"github.com/alacrity-engine/core/math/geometry"
	codec "github.com/alacrity-engine/resource-codec"
	"github.com/stretchr/testify/assert"
)

func TestGetSpritesheetFrames(t *testing.T) {
	pic := &codec.PictureData{
		Width:  512,
		Height: 256,
	}
	ss := &codec.SpritesheetData{
		Width:  3,
		Height: 1,
		Orig: codec.OrigData{
			X: 64,
			Y: 128,
		},
		Area: codec.AreaData{
			PixelWidth:  3 * 64,
			PixelHeight: 128,
		},
	}

	frames, err := pic.GetSpritesheetFrames(ss)
	assert.Nil(t, err)
	assert.Len(t, frames, 3)
	assert.ElementsMatch(t, frames, []geometry.Rect{
		geometry.R(64, 128, 128, 256),
		geometry.R(128, 128, 192, 256),
		geometry.R(192, 128, 256, 256),
	})
}
