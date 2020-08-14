package assistant

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vberezny/signal-assistant/signal"
)

type command string

const (
	// TODO: env variable.
	sharedFolder = "/home/vlad/signal-assistant-shared/"

	store command = "!store" // Store an attachment in the shared folder. Must provide name.
	get           = "!get"   // Get a file by name from the shared folder.
	list          = "!list"  // List all files in the shared folder and send a reply.
	man           = "!man"   // List all commands and their options (manual).
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

	for msg := range a.cli.Messages {
		err := a.validateMessage(msg)
		if err != nil {
			a.errorHandler("Failed to validate message", err)
		}
		err = a.executeCommand(msg)
		if err != nil {
			a.errorHandler("Failed to execute command, ", err)
		}
	}
}

func (a *Assistant) validateMessage(msg *signal.Message) error {
	if msg.PhoneNumber != a.owner {
		return errors.New(fmt.Sprintf("Message arrived from unknown number %v", msg.PhoneNumber))
	}
	if len(msg.Text) == 0 || string(msg.Text[0]) != "!" {
		return errors.New(fmt.Sprintf("Invalid command format. Message Text: %v", msg.Text))
	}
	return nil
}

func (a *Assistant) executeCommand(msg *signal.Message) (err error) {
	args := strings.Split(msg.Text, " ")
	command := command(args[0])
	switch command {
	case store:
		if len(args[1:]) < 2 {
			err = errors.New(fmt.Sprintf("Expected 2 arguments, got %v", len(args[1:])))
			return
		}
		err = a.storeAttachments(msg, args[1], args[2])
		if err != nil {
			return
		}
	case list:
		log.Print("TODO")
	default:
		err = a.sendMessage("Invalid command, type !man to see a list of available commands", nil)
		if err != nil {
			log.Printf("Failed to send message, %v", err)
		}
	}
	return
}

func (a *Assistant) storeAttachments(msg *signal.Message, fileName, fileExtension string) (err error) {
	if len(msg.Attachments) == 0 {
		err = errors.New("No attachments to store.")
		return
	}
	// TODO: add support for multiple attachments.
	if len(msg.Attachments) > 1 {
		err = errors.New("Only one attachment supported at the moment.")
		return
	}
	// TODO: detect mime type and use it to generate a fileExtension.
	fullPath := sharedFolder + fileName + "." + fileExtension
	err = copy(msg.Attachments[0], fullPath)
	if err == nil {
		err = a.sendMessage("Saved attachment at "+fullPath, nil)
		if err != nil {
			log.Printf("Failed to send message, %v", err)
		}
	}
	return
}

func (a *Assistant) sendMessage(text string, attachments []string) (err error) {
	msg := signal.NewMessage(time.Now(), a.owner, text, attachments)
	err = a.cli.SendMessage(msg)
	return
}

// A helper method to log an error and send a notification to the Assistant owner.
func (a *Assistant) errorHandler(message string, err error) {
	txt := fmt.Sprintf("%v %v", message, err)
	log.Printf(txt)
	// Notify the owner of any errors.
	err = a.sendMessage(txt, nil)
	if err != nil {
		log.Printf("Failed to send message, %v", err)
	}
}

// Copy the src file to dst. Any existing file will be overwritten.
func copy(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}
	err = out.Close()
	return
}
