package codec

import "encoding/gob"

func init() {
	gob.Register(GameObjectPointerData{})
	gob.Register(ComponentPointerData{})
	gob.Register(ResourcePointerData{})
	gob.Register(BatchPointerData{})
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

type BatchPointerData struct {
	CanvasID string
	BatchID  string
}
