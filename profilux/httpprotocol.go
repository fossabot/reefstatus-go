package profilux

import (
	"errors"
	"fmt"
	"github.com/cjburchell/reefstatus-go/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type httpProtocol struct {
	address string
}

func newHttpProtocol(settings ConnectionSettings) (iProtocol, error) {
	var p httpProtocol

	p.address = fmt.Sprintf("%s:%d", settings.Address, settings.Port)

	return &p, nil
}

func (p *httpProtocol) SendData(code, data int) error {

	url := fmt.Sprintf("http://%s/communication.php?dir=enq&code=%d&data=%d", p.address, code, data)

	log.Debugf("SendData: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(reply) == "Access Denied" {
		return errors.New("Access Denied")
	}

	command := getCommandFromReply(reply)
	dataParam := getDataFromReply(reply)

	if command != code {
		return fmt.Errorf("unexpected comand reply: %d", command)
	}

	if strings.HasPrefix(dataParam, "NACK") {
		return fmt.Errorf("error in command: %d", command)
	}

	if !strings.HasPrefix(dataParam, "ACK") {
		return errors.New("Unexpected message: Missing ACK")
	}

	return nil
}

func getCommandFromReply(reply []byte) int {
	code, _ := strconv.Atoi(strings.SplitN(strings.SplitN(string(reply), "&", 2)[0], "=", 2)[1])
	return code
}

func getDataFromReply(reply []byte) string {
	return strings.SplitN(strings.SplitN(string(reply), "&", 2)[1], "=", 2)[1]
}

func (p *httpProtocol) GetDataText(code int) (string, error) {
	reply, err := p.getRawData(code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(*reply), nil
}

func (p *httpProtocol) GetData(code int) (int, error) {
	data, err := p.getRawData(code)
	if err != nil {
		log.Errorf(err, "GetData %s", err.Error())
		return 0, err
	}

	result, err := strconv.Atoi(*data)
	if err != nil {
		log.Errorf(err, "GetData: %s", err.Error())
		return 0, err
	}

	return result, nil
}

func (p *httpProtocol) getRawData(code int) (*string, error) {
	url := fmt.Sprintf("http://%s/communication.php?dir=enq&code=%d", p.address, code)

	log.Debugf("Get %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	command := getCommandFromReply(reply)
	dataParam := getDataFromReply(reply)

	if command != code {
		return nil, fmt.Errorf("unexpected comand reply: %d, %s", command, reply)
	}

	if strings.HasPrefix(dataParam, "NACK") {
		return nil, fmt.Errorf("error in command: %d", command)
	}

	return &dataParam, nil
}

func (p *httpProtocol) Disconnect() {
}
