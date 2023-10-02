package codec

import "encoding/gob"

func init() {
	gob.Register(GameObjectPointerData{})
	gob.Register(ComponentPointerData{})
	gob.Register(ResourcePointerData{})
}

type GameObjectPointerData struct {
	Name string
}

type ComponentPointerData struct {
	GmobName string
	CompType string
}

type ResourcePointerData struct {
	ResourceType string
	ResourceID   string
}
