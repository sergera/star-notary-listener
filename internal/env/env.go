package env

import (
	"log"
	"os"
	"strconv"
)

var InfuraWebsocketURL string
var ContractAddress string
var ConfirmedThreshold uint64
var OrphanThreshold uint64

func Init() {
	setInfuraWebsocketURL()
	setContractAddress()
	setConfirmedThreshold()
	setOrphanThreshold()
}

func setInfuraWebsocketURL() {
	infuraWebsocketURL, exists := os.LookupEnv("INFURA_WEBSOCKET_URL")
	if !exists {
		log.Fatal("Infura websocket URL environment variable not found")
	}

	InfuraWebsocketURL = infuraWebsocketURL
}

func setContractAddress() {
	contractAddress, exists := os.LookupEnv("CONTRACT_ADDRESS")
	if !exists {
		log.Fatal("Contract address environment variable not found")
	}

	ContractAddress = contractAddress
}

func setConfirmedThreshold() {
	confirmedThresholdString, exists := os.LookupEnv("CONFIRMED_THRESHOLD")
	if !exists {
		log.Fatal("Confirmed threshold environment variable not found")
	}

	confirmedThreshold, err := strconv.ParseUint(confirmedThresholdString, 10, 64)
	if err != nil {
		log.Fatal("Could not convert confirmed threshold environment variable to uint")
	}

	ConfirmedThreshold = confirmedThreshold
}

func setOrphanThreshold() {
	orphanThresholdString, exists := os.LookupEnv("ORPHAN_THRESHOLD")
	if !exists {
		log.Fatal("Orphan threshold environment variable not found")
	}

	orphanThrehold, err := strconv.ParseUint(orphanThresholdString, 10, 64)
	if err != nil {
		log.Fatal("Could not convert orphan threshold environment variable to uint")
	}

	OrphanThreshold = orphanThrehold
}
