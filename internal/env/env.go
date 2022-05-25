package env

import (
	"log"
	"os"
	"strconv"
)

var InfuraWebsocketURL string
var ContractAddress string
var ConfirmedThreshold uint64
var OrphanedThreshold uint64
var SleepIntervalSeconds uint64
var LogPath string

func Init() {
	setInfuraWebsocketURL()
	setContractAddress()
	setConfirmedThreshold()
	setOrphanedThreshold()
	setSleepIntervalSeconds()
	setLogPath()
}

func setInfuraWebsocketURL() {
	infuraWebsocketURL, exists := os.LookupEnv("INFURA_WEBSOCKET_URL")
	if !exists {
		log.Panic("Infura websocket URL environment variable not found")
	}

	InfuraWebsocketURL = infuraWebsocketURL
}

func setContractAddress() {
	contractAddress, exists := os.LookupEnv("CONTRACT_ADDRESS")
	if !exists {
		log.Panic("Contract address environment variable not found")
	}

	ContractAddress = contractAddress
}

func setConfirmedThreshold() {
	confirmedThresholdString, exists := os.LookupEnv("CONFIRMED_THRESHOLD")
	if !exists {
		log.Panic("Confirmed threshold environment variable not found")
	}

	confirmedThreshold, err := strconv.ParseUint(confirmedThresholdString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert confirmed threshold environment variable to uint: %+v\n", err)
	}

	ConfirmedThreshold = confirmedThreshold
}

func setOrphanedThreshold() {
	orphanedThresholdString, exists := os.LookupEnv("ORPHANED_THRESHOLD")
	if !exists {
		log.Panic("Orphaned threshold environment variable not found")
	}

	orphanedThrehold, err := strconv.ParseUint(orphanedThresholdString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert orphaned threshold environment variable to uint: %+v\n", err)
	}

	OrphanedThreshold = orphanedThrehold
}

func setSleepIntervalSeconds() {
	sleepIntervalSecondsString, exists := os.LookupEnv("SLEEP_INTERVAL_SECONDS")
	if !exists {
		log.Panic("Sleep interval seconds environment variable not found")
	}

	sleepIntervalSeconds, err := strconv.ParseUint(sleepIntervalSecondsString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert sleep interval seconds environment variable to uint: %+v\n", err)
	}

	SleepIntervalSeconds = sleepIntervalSeconds
}

func setLogPath() {
	logPath, exists := os.LookupEnv("LOG_PATH")
	if !exists {
		log.Panic("Log path environment variable not found")
	}

	LogPath = logPath
}
