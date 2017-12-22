package profilux

import (
	"errors"
	"fmt"
	"github.com/cjburchell/reefstatus-go/common"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/profilux/types"
	"strconv"
	"time"
)

type protocol struct {
	Connection *connection
	Address    int
}

var ProtocolError = errors.New("protocol error")

func codeError(code byte) error {
	return fmt.Errorf("code error %d", code)
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

		if err != ProtocolError {
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

		if err != ProtocolError {
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

		if err != ProtocolError {
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

		if err != ProtocolError {
			return 0, err
		}
	}
}

func (protocol protocol) GetDataDate(code int) (time.Time, error) {
	result, err := protocol.GetData(code)
	if err != nil {
		return time.Now(), err
	}

	timeString := strconv.Itoa(result)

	if len(timeString) == 6 {
		yearValue, _ := strconv.Atoi(timeString[len(timeString)-2:])
		monthValue, _ := strconv.Atoi(timeString[len(timeString)-4 : len(timeString)-2])
		dateValue, _ := strconv.Atoi(timeString[:len(timeString)-4])
		return time.Date(yearValue+2000, time.Month(monthValue), dateValue, 0, 0, 0, 0, time.UTC), nil
	} else if len(timeString) == 7 {
		yearValue, _ := strconv.Atoi(timeString[len(timeString)-3:])
		monthValue, _ := strconv.Atoi(timeString[len(timeString)-5 : len(timeString)-3])
		dateValue, _ := strconv.Atoi(timeString[:len(timeString)-5])
		return time.Date(yearValue+2000, time.Month(monthValue), dateValue, 0, 0, 0, 0, time.UTC), nil
	}

	return time.Now(), err
}

func (protocol protocol) GetDataEnum(code int, convert func(int) string) (string, error) {
	result, err := protocol.GetData(code)
	if err != nil {
		return "", err
	}

	return convert(result), nil
}

func (protocol protocol) GetDataCurrentState(code int) (types.CurrentState, error) {
	result, err := protocol.GetData(code)
	if err != nil {
		return "", err
	}

	return types.GetCurrentState(result), nil
}

func (protocol protocol) GetDataFloat(code int, multiplier float64) (float64, error) {
	result, err := protocol.GetData(code)
	if err != nil {
		return 0, err
	}

	return float64(result) * multiplier, nil
}

func (protocol protocol) GetDataMultiplier(code int, multiplier int) (int, error) {
	result, err := protocol.GetData(code)
	if err != nil {
		return 0, err
	}

	return result * multiplier, nil
}

func (protocol protocol) GetDataBool(code int) (bool, error) {
	result, err := protocol.GetData(code)
	if err != nil {
		return false, err
	}

	return result != 0, nil
}

func (protocol protocol) GetDataFloatAndRound(code int, multiplier float64, digits int) (float64, error) {
	result, err := protocol.GetDataFloat(code, multiplier)
	if err != nil {
		return 0, err
	}

	return common.Round(result, digits), nil
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

		if err != ProtocolError {
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

		if err != ProtocolError {
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

		if err != ProtocolError {
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

		if err != ProtocolError {
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
		return ProtocolError
	}

	if reply[4] == EOT {
		log.Warn("Unexpected Message: Empty Reply")
		return ProtocolError
	}

	if reply[4] == STX {
		if reply[5] == NAK {
			errorCode := reply[6]
			return codeError(errorCode)
		}

		if reply[5] == ACK {
			log.Warn("Unexpected Message: ACK")
			return ProtocolError
		}

		// should be ok we must now look for the code and verify it
		replyCode := getGetMessageCode(reply)
		if replyCode != code {
			log.Warnf("Unexpected Message: Wrong Code Expecting %d Got %d", code, replyCode)
			return ProtocolError
		}
	} else {
		log.Warn("Unknown message type")
		return ProtocolError
	}

	return nil
}

func verifyAckPacket(reply []byte) error {
	if len(reply) < 4 {
		// strange packet size!
		log.Warnf("Expecting Packet size of at least 4")
		return ProtocolError
	}

	if reply[4] == EOT {
		log.Warn("Unexpected Message: Empty Reply")
		return ProtocolError
	}

	if reply[4] == STX {
		if reply[5] == NAK {
			errorCode := reply[6]
			return codeError(errorCode)
		}

		if reply[5] != ACK {
			replyCode := getGetMessageCode(reply)
			log.Warnf("Unexpected Message: Code %d", replyCode)
			return ProtocolError
		}

	} else {
		log.Warn("Unknown message type")
		return ProtocolError
	}

	return nil
}
