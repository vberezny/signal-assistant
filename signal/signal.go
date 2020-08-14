package signal

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	SIGNAL_CLI_DBUS_SERVICE = "org.asamk.Signal"
)

type Signal struct {
	Messages chan *Message
}

func NewSignal() *Signal {
	messages := make(chan *Message, 10)
	return &Signal{messages}
}

// Listen() establishes a connection to the DBus service and listens for
// incoming Signal messages.
func (s *Signal) Listen() {
	// TODO: check if service is started and start it if not.
	signals := make(chan *dbus.Signal, 10)
	conn, err := connectDBus(signals)
	if err != nil {
		log.Fatalf("Failed to connect to DBus %v", err)
	}
	defer conn.Close()

	for signal := range signals {
		// Read receipts are of no interest to this application.
		if signal.Name == "org.asamk.Signal.ReceiptReceived" {
			continue
		}
		message, err := newMessageFromSignal(signal)
		if err != nil {
			log.Printf("Failed to parse new message from signal %v", err)
			continue
		}
		s.Messages <- message
	}
}

// SendMessage() uses the signal-cli command in dbus mode. This is due to
// the method org.asamk.Signal.sendMessage not working as expected when called
// using dbus.Object.Call().
// TODO: further investigate using org.asamk.Signal.sendMessage method.
func (s *Signal) SendMessage(msg *Message) (err error) {
	if msg.PhoneNumber[0:2] != "+1" && len(msg.PhoneNumber) != 12 {
		err = errors.New(fmt.Sprintf("Unable to send message, phone number format incorrect: %v", msg.PhoneNumber))
		return
	}
	args := []string{"--dbus", "send", "-m", msg.Text, msg.PhoneNumber}
	if len(msg.Attachments) > 0 {
		args = append(args, "-a")
		for _, attachment := range msg.Attachments {
			args = append(args, attachment)
		}
	}
	cmd := exec.Command("signal-cli", args...)
	_, err = cmd.Output()
	if err != nil {
		return
	}
	log.Println("Message sent")
	return
}

// Establish a connection to the org.asamk.Signal DBus interface on the Session Bus.
func connectDBus(signals chan<- *dbus.Signal) (conn *dbus.Conn, err error) {
	conn, err = dbus.SessionBus()
	if err != nil {
		return
	}
	if err = verifyConnection(conn); err != nil {
		return
	}
	options := dbus.WithMatchSender(SIGNAL_CLI_DBUS_SERVICE)
	err = conn.AddMatchSignal(options)
	if err != nil {
		return
	}
	conn.Signal(signals)
	return
}

// Verify that org.asamk.Signal is available via the DBus interface.
func verifyConnection(conn *dbus.Conn) (err error) {
	s := []string{}
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return
	}
	found := false
	for _, v := range s {
		if v == SIGNAL_CLI_DBUS_SERVICE {
			found = true
			break
		}
	}
	if !found {
		err = errors.New("signal-cli connection not found on Dbus.")
	}
	return
}

type Message struct {
	Time        time.Time
	PhoneNumber string
	Text        string
	Attachments []string
}

func NewMessage(t time.Time, phone, text string, attachments []string) *Message {
	return &Message{
		Time:        t,
		PhoneNumber: phone,
		Text:        text,
		Attachments: attachments,
	}
}

// Helper method that transforms a *dbus.Signal to a *Message.
func newMessageFromSignal(signal *dbus.Signal) (msg *Message, err error) {
	utc, ok := signal.Body[0].(int64)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert time to int64, %v\n", signal.Body[0]))
		return
	}
	t := time.Unix(utc, 0)
	phone, ok := signal.Body[1].(string)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert phone number to string, %v\n", signal.Body[1]))
		return
	}
	text, ok := signal.Body[3].(string)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert message text to string, %v\n", signal.Body[3]))
		return
	}
	attachments, ok := signal.Body[4].([]string)
	if !ok {
		err = errors.New(fmt.Sprintf("failed to convert attachment path to slice of strings, %v", signal.Body[4]))
		return
	}
	msg = NewMessage(t, phone, text, attachments)
	return
}
