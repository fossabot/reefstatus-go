package native

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/profilux/settings"
	"net"
)

type connection struct {
	conn net.Conn
}

func newConnection(settings settings.ConnectionSettings) (*connection, error) {
	log.Printf("Connecting to %s:%d", settings.Address, settings.Port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", settings.Address, settings.Port))
	if err != nil {
		return nil, err
	}

	return &connection{conn: conn}, nil
}

func (c connection) Write(data []byte) error {
	log.Printf("Write Buffer")
	_, err := c.conn.Write(data)
	return err
}
func (c connection) Read(size int) ([]byte, int, error) {
	log.Printf("Read Buffer")
	buf := make([]byte, size)
	read, err := c.conn.Read(buf)
	log.Printf("Read %d %v", read, buf)
	return buf, read, err
}
func (c connection) Disconnect() {
	c.conn.Close()
}
