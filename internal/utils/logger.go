package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(message string) {
	Logger.Println("INFO: " + message)
}

func LogError(message string) {
	Logger.Println("ERROR: " + message)
}
