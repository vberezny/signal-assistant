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
// from using dbus.Object.Call().
// TODO: further investigate using org.asamk.Signal.sendMessage method.
func (s *Signal) SendMessage(msg *Message) (err error) {
	if msg.PhoneNumber[0:2] != "+1" && len(msg.PhoneNumber) != 12 {
		err = errors.New(fmt.Sprintf("Unable to send message, phone number format incorrect: %v", msg.PhoneNumber))
		return
	}
	args := []string{"--dbus", "send", "-m", msg.Text, msg.PhoneNumber}
	if msg.Attachment != "" {
		args = append(args, "-a", msg.Attachment)
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
	Attachment  string
	MessageType string
}

func NewMessage(t time.Time, phone, text, attachment, msgType string) *Message {
	return &Message{
		Time:        t,
		PhoneNumber: phone,
		Text:        text,
		Attachment:  attachment,
		MessageType: msgType,
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
	attachment, ok := signal.Body[4].(string)
	if !ok {
		// A message without an attachment is acceptable. Just log for now.
		log.Printf("Failed to convert attachment path to string, %v", signal.Body[4])
		attachment = ""
	}
	msg = NewMessage(t, phone, text, attachment, signal.Name)
	return
}
