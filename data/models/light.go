package models

type Light struct {
	BaseInfo
	Value                   float64
	Channel                 int
	IsDimmable              bool
	OperationHours          int
	MaxOperationHours       int
	EnableMaxOperationHours bool
}

func (light Light) IsLightOn() bool {
	return light.Value != 0
}

func (light Light) IsOverMaxOperationHours() bool {
	return light.EnableMaxOperationHours && light.MaxOperationHours < (light.OperationHours/60.0)
}

func CreateLight() Light {
	var light Light
	light.Type = "Light"
	light.Units = "%"
	return light
}
