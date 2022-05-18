package env

import (
	"log"
	"os"
	"strconv"
)

var InfuraProjectID string
var ContractAddress string
var ConfirmationsThreshold uint64

func Init() {
	setInfuraProjectID()
	setContractAddress()
	setConfirmationsThreshold()
}

func setInfuraProjectID() {
	infuraProjectId, exists := os.LookupEnv("INFURA_PROJECT_ID")
	if !exists {
		log.Fatal("Infura project id environment variable not found")
	}

	InfuraProjectID = infuraProjectId
}

func setContractAddress() {
	contractAddress, exists := os.LookupEnv("CONTRACT_ADDRESS")
	if !exists {
		log.Fatal("Contract address environment variable not found")
	}

	ContractAddress = contractAddress
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
