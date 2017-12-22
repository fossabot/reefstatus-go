package profilux

import (
	"fmt"
	"github.com/pkg/errors"
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
	resp, err := http.Get(fmt.Sprintf("http://%s/communication.php?dir=enq&code%d&data%d", p.address, code, data))
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
	code, _ := strconv.Atoi(strings.SplitN(strings.SplitN(string(reply), "&", 2)[1], "=", 2)[1])
	return code
}

func getDataFromReply(reply []byte) string {
	return strings.Split(strings.Split(string(reply), "&")[0], "=")[1]
}

func (p *httpProtocol) GetDataText(code int) (string, error) {
	reply, err := p.GetData(code)
	if err != nil {
		return "", err
	}

	var data string

	for reply != 0 {
		value := byte(reply & 0xff)
		reply >>= 4
		data += string(value)
	}

	return data, nil
}

func (p *httpProtocol) GetData(code int) (int, error) {
	data, err := p.getRawData(code)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *httpProtocol) getRawData(code int) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/communication.php?dir=enq&code%d", p.address, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (p *httpProtocol) Disconnect() {
}
