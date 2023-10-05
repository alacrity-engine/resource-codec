package codec

import (
	"bytes"
	"encoding/gob"

	"github.com/alacrity-engine/core/math/geometry"
)

type PrefabData struct {
	Name          string
	TransformRoot *TransformData
}

type TransformData struct {
	Position geometry.Vec
	Angle    float64
	Scale    geometry.Vec
	Gmob     *GameObjectData
	Children []*TransformData
}

type GameObjectData struct {
	Name       string
	ZUpdate    float64
	Components []*ComponentData
	Sprite     *SpriteData
	Draw       bool
}

type ComponentData struct {
	TypeName string
	Active   bool
	Data     map[string]interface{}
}

type SpriteData struct {
	ColorMask       []float32
	TargetArea      geometry.Rect
	ZDraw           float32
	VertexDrawMode  uint32
	TextureDrawMode uint32
	ColorDrawMode   uint32
	ShaderProgramID string
	TextureID       string
	CanvasID        string
	BatchID         string
}

func (pdata *PrefabData) ToBytes() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buffer)

	err := enc.Encode(pdata)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func PrefabDataFromBytes(data []byte) (*PrefabData, error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	pdata := &PrefabData{}

	err := dec.Decode(pdata)

	if err != nil {
		return nil, err
	}

	return pdata, nil
}
