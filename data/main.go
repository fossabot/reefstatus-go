package main

import (
	"encoding/json"
	"fmt"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/profilux"
)

func main() {
	controller := newController()
	settings := profilux.CreateDefaultConnectionSettings()
	profiluxController, err := profilux.NewController(settings)
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	controller.Update(profiluxController)

	result, _ := json.MarshalIndent(controller, "", "    ")
	fmt.Printf("%s", result)
}
