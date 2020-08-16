package main

import (
	"log"
	"os"

	"github.com/vberezny/signal-assistant/assistant"
)

func main() {
	phoneNumber := os.Getenv("OWNER_NUMBER")
	if phoneNumber == "" {
		log.Fatal("Phone number not provided.")
	}
	sharedFolder := os.Getenv("ASSISTANT_FOLDER")
	if sharedFolder == "" {
		log.Fatal("Shared folder not provided")
	}
	_, err := os.Stat(sharedFolder)
	if os.IsNotExist(err) {
		log.Fatal("Provided shared folder does not exists")
	}
	// The Assistant expects the shared folder to have the trailing slash included.
	if string(sharedFolder[len(sharedFolder)-1]) != "/" {
		sharedFolder += "/"
	}
	assistant := assistant.NewAssistant(phoneNumber, sharedFolder)
	assistant.Run()
}
