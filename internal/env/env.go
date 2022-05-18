package env

import (
	"log"
	"os"
	"strconv"
)

var InfuraProjectID string
var ConfirmationsThreshold uint64

func Init() {
	setInfuraProjectID()
	setConfirmationsThreshold()
}

func setInfuraProjectID() {
	infuraProjectId, exists := os.LookupEnv("INFURA_PROJECT_ID")
	if !exists {
		log.Fatal("Infura project id environment variable not found")
	}

	InfuraProjectID = infuraProjectId
}

func setConfirmationsThreshold() {
	confirmationsThresholdString, exists := os.LookupEnv("CONFIRMATIONS_THRESHOLD")
	if !exists {
		log.Fatal("Confirmations threshold environment variable not found")
	}

	confirmationsThreshold, err := strconv.ParseUint(confirmationsThresholdString, 10, 64)
	if err != nil {
		log.Fatal("Could not convert onfirmations threshold environment variable to int")
	}

	ConfirmationsThreshold = confirmationsThreshold
}
