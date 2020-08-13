package assistant

import (
	"log"

	"github.com/vberezny/signal-assistant/signal"
)

const (
	command1 = "test"
	command2 = "test2"
)

// TODO: name the member?
type Assistant struct {
	*signal.Signal
}

func NewAssistant() *Assistant {
	return &Assistant{signal.NewSignal()}
}

func (a *Assistant) Run() {
	go a.Listen()

	for message := range a.Messages {
		log.Print(message)
		a.processMessage(message)
	}
}

func (a *Assistant) processMessage(message *signal.Message) {
	// TODO: hashmap of valid commands?
	if message.Text == command1 {
		log.Print("valid command")
		message.Text = "you passed"
		a.SendMessage(message)
	} else {
		// TODO: if the message came from the right number, reply back saying something
		// about the command being wrong or malformed.
		log.Print("invalid command")
	}
}
