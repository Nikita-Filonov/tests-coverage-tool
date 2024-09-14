package logger

import (
	"log"
)

func StartMakeReport(report string) {
	log.Printf("Starting to make %s report file", report)
}

func SuccessfullyMadeReport(report string) {
	log.Printf("Successfully made %s report file", report)
}

func ErrorMakingReport(report string) {
	log.Printf("Error making %s report file", report)
}

func FatalSavingReport(report string, err error) {
	log.Fatalf("Error saving %s report: %v", report, err)
}

func EnvVariableEmptySkipping(variable string) {
	log.Printf("Env variable '%s' empty, skipping", variable)
}

func FatalBuildingNewClient(client string, err error) {
	log.Fatalf("Error building %s: %v", client, err)
}

func FatalGettingEntity(client string, err error) {
	log.Fatalf("Error building %s: %v", client, err)
}
