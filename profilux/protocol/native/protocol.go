package native

import (
	"errors"
	"fmt"
	"github.com/cjburchell/reefstatus-go/common/log"
	protocol2 "github.com/cjburchell/reefstatus-go/profilux/protocol"
	"github.com/cjburchell/reefstatus-go/profilux/settings"
)

type protocol struct {
	Connection *connection
	Address    int
}

var protocolError = errors.New("protocol error")

func codeError(code byte) error {
	return fmt.Errorf("code error %d", code)
}

func NewProtocol(settings settings.ConnectionSettings) (protocol2.IProtocol, error) {
	con, err := newConnection(settings)
	if err != nil {
		return nil, err
	}

	var p protocol

	p.Address = settings.ControllerAddress
	p.Connection = con

	return &p, nil
}

func (protocol protocol) Disconnect() {
	protocol.Connection.Disconnect()
}

func (protocol protocol) SendData(code, data int) error {
	sendCommandInt(code, data, protocol.Connection, protocol.Address)

	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return err
		}

		err = verifyAckPacket(reply)
		if err == nil {
			return nil
		}

		if err != protocolError {
			return err
		}
	}
}

func (protocol protocol) SendText(code int, data string) error {
	sendTextCommand(code, data, protocol.Connection, protocol.Address)

	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return err
		}

		err = verifyAckPacket(reply)
		if err == nil {
			return nil
		}

		if err != protocolError {
			return err
		}
	}
}

func (protocol protocol) GetDataText(code int) (string, error) {
	sendCommand(code, protocol.Connection, protocol.Address)
	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return "", err
		}

		err = verifyDataPacket(reply, code)
		if err == nil {
			return getMessageString(reply), nil
		}

		if err != protocolError {
			return "", err
		}
	}
}

func (protocol protocol) GetData(code int) (int, error) {

	log.Printf("Sending GetData %d", code)
	sendCommand(code, protocol.Connection, protocol.Address)
	for {

		log.Printf("Reading Repsonce")
		reply, err := protocol.readPacket()
		if err != nil {
			return 0, err
		}

		err = verifyDataPacket(reply, code)
		if err == nil {
			log.Printf("Got valid data")
			return getMessageData(reply), nil
		}

		if err != protocolError {
			return 0, err
		}
	}
}

func (protocol protocol) GetDataShortArray(code int) ([]int, error) {
	sendCommand(code, protocol.Connection, protocol.Address)

	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return nil, err
		}

		err = verifyDataPacket(reply, code)
		if err == nil {
			return getMessageDataShortArray(reply), nil
		}

		if err != protocolError {
			return nil, err
		}
	}
}

func (protocol protocol) GetDataByteArray(code int) ([]byte, error) {
	sendCommand(code, protocol.Connection, protocol.Address)

	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return nil, err
		}

		err = verifyDataPacket(reply, code)
		if err == nil {
			return getMessageBytes(reply), nil
		}

		if err != protocolError {
			return nil, err
		}
	}
}

func (protocol protocol) GetDataTwoByteArray(code int) ([]int, error) {
	sendCommand(code, protocol.Connection, protocol.Address)

	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return nil, err
		}

		err = verifyDataPacket(reply, code)
		if err == nil {
			return getMessageDataTwoByteArray(reply), nil
		}

		if err != protocolError {
			return nil, err
		}
	}
}

func (protocol protocol) GetDataBoolArray(code int) ([]bool, error) {
	sendCommand(code, protocol.Connection, protocol.Address)

	for {
		reply, err := protocol.readPacket()
		if err != nil {
			return nil, err
		}

		err = verifyDataPacket(reply, code)
		if err == nil {
			return getMessageBools(reply), nil
		}

		if err != protocolError {
			return nil, err
		}
	}
}

func (protocol protocol) readPacket() (reply []byte, err error) {
	for {
		data, size, err := protocol.Connection.Read(1)
		if err != nil {
			return nil, err
		}

		if size == 0 {
			break
		}

		reply = append(reply, data...)

		if atEndOfPacket(reply) {
			break
		}
	}

	return
}

func verifyDataPacket(reply []byte, code int) error {
	if len(reply) < 4 {
		// strange packet size!
		log.Warn("Expecting Packet size of at least 4")
		return protocolError
	}

	if reply[4] == EOT {
		log.Warn("Unexpected Message: Empty Reply")
		return protocolError
	}

	if reply[4] == STX {
		if reply[5] == NAK {
			errorCode := reply[6]
			return codeError(errorCode)
		}

		if reply[5] == ACK {
			log.Warn("Unexpected Message: ACK")
			return protocolError
		}

		// should be ok we must now look for the code and verify it
		replyCode := getGetMessageCode(reply)
		if replyCode != code {
			log.Warnf("Unexpected Message: Wrong Code Expecting %d Got %d", code, replyCode)
			return protocolError
		}
	} else {
		log.Warn("Unknown message type")
		return protocolError
	}

	return nil
}

func verifyAckPacket(reply []byte) error {
	if len(reply) < 4 {
		// strange packet size!
		log.Warnf("Expecting Packet size of at least 4")
		return protocolError
	}

	if reply[4] == EOT {
		log.Warn("Unexpected Message: Empty Reply")
		return protocolError
	}

	if reply[4] == STX {
		if reply[5] == NAK {
			errorCode := reply[6]
			return codeError(errorCode)
		}

		if reply[5] != ACK {
			replyCode := getGetMessageCode(reply)
			log.Warnf("Unexpected Message: Code %d", replyCode)
			return protocolError
		}

	} else {
		log.Warn("Unknown message type")
		return protocolError
	}

	return nil
}
