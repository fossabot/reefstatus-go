package profilux

type ConnectionSettings struct {
	Address           string
	Port              int
	Timeout           int
	ControllerAddress int
}

func CreateDefaultConnectionSettings() (connectionSettings ConnectionSettings) {
	connectionSettings.Address = "192.168.3.10"
	connectionSettings.Port = 10001
	connectionSettings.Timeout = 5000
	connectionSettings.ControllerAddress = 1
	return
}
