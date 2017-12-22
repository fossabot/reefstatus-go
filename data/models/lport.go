package models

const (
	LValueMin = 18.0
	LValueMax = 255.00
)

type LPort struct {
	DeviceInfo
	Value float64
}

func (lport *LPort) SetValue(value int) {
	if value < LValueMin {
		lport.Value = 0
	} else {
		lport.Value = (float64(value) - LValueMin) / (LValueMax - LValueMin) * 100.0
	}
}

func NewLPort() *LPort {
	var lPort LPort
	lPort.Type = "LPort"
	lPort.Units = "%"
	return &lPort
}
