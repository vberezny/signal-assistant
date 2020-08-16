package assistant

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vberezny/signal-assistant/signal"
)

type Assistant struct {
	cli          *signal.Signal
	owner        string
	commands     []command
	sharedFolder string
}

func NewAssistant(owner, sharedFolder string) *Assistant {
	commands := getAllCommands()
	return &Assistant{
		cli:          signal.NewSignal(),
		owner:        owner,
		commands:     commands,
		sharedFolder: sharedFolder,
	}
}

func (a *Assistant) Run() {
	go a.cli.Listen()

	for msg := range a.cli.Messages {
		err := a.validateMessage(msg)
		if err != nil {
			a.errorHandler("Failed to validate message,", err)
			continue
		}
		err = a.executeCommand(msg)
		if err != nil {
			a.errorHandler("Failed to execute command,", err)
		}
	}
}

func (a *Assistant) validateMessage(msg *signal.Message) error {
	if msg.PhoneNumber != a.owner {
		return errors.New(fmt.Sprintf("Message arrived from unknown number %v", msg.PhoneNumber))
	}
	if len(msg.Text) == 0 || string(msg.Text[0]) != "!" {
		return errors.New(fmt.Sprintf("Invalid command format. Must start with !. Message Text: %v", msg.Text))
	}
	return nil
}

func (a *Assistant) executeCommand(msg *signal.Message) (err error) {
	splitMessage := strings.Split(msg.Text, " ")
	args := []string{}
	if len(splitMessage) > 1 {
		args = splitMessage[1:]
	}
	command := command{
		cmd:  splitMessage[0],
		args: args,
	}
	switch command.cmd {
	case STORE:
		if len(command.args) < 2 {
			err = errors.New(fmt.Sprintf("Expected 2 arguments, got %v", len(command.args)))
			return
		}
		err = a.storeAttachments(msg.Attachments, command.args[0], command.args[1])
		if err != nil {
			return
		}
	case LIST:
		err = a.listSharedFiles()
		if err != nil {
			return
		}
	case GET:
		if len(command.args) < 1 {
			err = errors.New(fmt.Sprintf("Expected at least 1 argument, got %v", len(command.args)))
			return
		}
		err = a.returnFile(command.args)
		if err != nil {
			return
		}
	case MAN:
		err = a.returnManual()
		if err != nil {
			return
		}
	default:
		err = errors.New("Invalid command, type !man to see a list of available commands.")
	}
	return
}

// Handles !store command.
func (a *Assistant) storeAttachments(attachments []string, fileName, fileExtension string) (err error) {
	if len(attachments) == 0 {
		err = errors.New("No attachments to store.")
		return
	}
	// TODO: add support for multiple attachments.
	if len(attachments) > 1 {
		err = errors.New("Only one attachment supported at the moment.")
		return
	}
	// TODO: detect mime type and use it to generate a fileExtension.
	fullPath := a.sharedFolder + fileName + "." + fileExtension
	err = copy(attachments[0], fullPath)
	if err == nil {
		// No need to handle err as it will bubble up to the error handler.
		err = a.sendMessage("Saved attachment at "+fullPath, nil)
	}
	return
}

// Handles !list command.
func (a *Assistant) listSharedFiles() (err error) {
	files, err := ioutil.ReadDir(a.sharedFolder)
	if err != nil {
		return
	}
	output := "All files in " + a.sharedFolder + "\n"
	for _, f := range files {
		output += fmt.Sprintf("%v, size: %v\n", f.Name(), f.Size())
	}
	err = a.sendMessage(output, nil)
	return
}

// Handles !get command.
func (a *Assistant) returnFile(fileNames []string) (err error) {
	attachments := []string{}
	message := ""
	for _, fileName := range fileNames {
		fullPath := a.sharedFolder + fileName
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			// fileName exists.
			attachments = append(attachments, fullPath)
		} else {
			message += fmt.Sprintf("Unable to find %v.", fullPath)
		}
	}
	err = a.sendMessage(message, attachments)
	return
}

// Handles !man command.
func (a *Assistant) returnManual() (err error) {
	manual := getCommandManual()
	err = a.sendMessage(manual, nil)
	return
}

// Wraps signal.SendMessage.
func (a *Assistant) sendMessage(text string, attachments []string) (err error) {
	msg := signal.NewMessage(time.Now(), a.owner, text, attachments)
	err = a.cli.SendMessage(msg)
	return
}

// A helper method to log the error and send a notification to the Assistant owner.
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
