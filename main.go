package main

import (
	"os"

	"github.com/vberezny/signal-assistant/assistant"
)

func main() {
	assistant := assistant.NewAssistant(os.Getenv("OWNER_NUMBER"))
	assistant.Run()
}
