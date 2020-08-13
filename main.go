package main

import "github.com/vberezny/signal-assistant/assistant"

func main() {
	assistant := assistant.NewAssistant()
	assistant.Run()
}
