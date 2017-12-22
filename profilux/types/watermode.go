package types

type WaterMode int

const (
	WaterModeReady    = "Ready"
	WaterModeDraining = "Draining"
	WaterModeFilling  = "Filling"
)

var waterModeMap = map[int]string{
	0: WaterModeReady,
	1: WaterModeDraining,
	2: WaterModeFilling,
}

func GetWaterMode(value int) string {
	if val, ok := waterModeMap[value]; ok {
		return val
	}

	return "Unknown"
}
