package assistant

import (
	"errors"
	"fmt"
	"log"

	"github.com/vberezny/signal-assistant/signal"
)

type Command string

const (
	store    Command = "!store"
	get              = "!get"
	commands         = "!man" // List all commands (manual).
)

type Assistant struct {
	cli   *signal.Signal
	owner string
}

func NewAssistant(owner string) *Assistant {
	return &Assistant{
		cli:   signal.NewSignal(),
		owner: owner,
	}
}

func (a *Assistant) Run() {
	go a.cli.Listen()

	for message := range a.cli.Messages {
		err := a.processMessage(message)
		if err != nil {
			log.Printf("Error while processing message: %v", err)
		}
	}
}

// TODO: hashmap of valid commands?
func (a *Assistant) processMessage(msg *signal.Message) (err error) {
	if msg.PhoneNumber != a.owner {
		err = errors.New(fmt.Sprintf("message arrived from unknown number %v", msg.PhoneNumber))
		return
	}
	if string(msg.Text[0]) != "!" {
		// TODO: reply saying invalid command.
		err = errors.New(fmt.Sprintf("invalid command format, must start with !. Message Text: %v", msg.Text))
		return
	}
	command := Command(msg.Text)
	if command == store {
		log.Print("valid command")
		msg.Text = "you passed"
		err = a.cli.SendMessage(msg)
		if err != nil {
			return
		}
	} else {
		log.Print("invalid command")
	}
	return
}
