package settings

import (
	"github.com/cjburchell/reefstatus-go/common"
)

const ProtocolHTTP = "HTTP"
const ProtocolSocket = "Socket"

type ConnectionSettings struct {
	Address           string
	Port              int
	Timeout           int
	ControllerAddress int
	Protocol          string
}

func NewConnectionSettings() (connection ConnectionSettings) {
	connection.Address = common.GetEnv("PROFILUX_ADDRESS", "192.168.3.10")
	connection.Port = common.GetEnvInt("PROFILUX_PORT", 80)
	connection.Protocol = common.GetEnv("PROFILUX_PROTOCOL", "HTTP")
	connection.Timeout = 5000
	connection.ControllerAddress = 1
	return
}

var Connection = NewConnectionSettings()
