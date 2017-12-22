package main

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/profilux"
)

func main() {
	var controller Controller
	settings := profilux.CreateDefaultConnectionSettings()
	profiluxController, err := profilux.NewController(settings)
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	controller.Update(profiluxController)

	fmt.Printf("%v", controller)
}
