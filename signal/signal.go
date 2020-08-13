package signal

import (
	"errors"
	"log"
	"os/exec"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	SIGNAL_CLI_DBUS_SERVICE = "org.asamk.Signal"
)

// TODO: better error handling
// TODO: better logging

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
	signals := make(chan *dbus.Signal, 10)
	conn := connectDBus(signals)
	defer conn.Close()

	for signal := range signals {
		// Read receipts are of no interest to this application.
		if signal.Name == "org.asamk.Signal.ReceiptReceived" {
			continue
		}
		message := newMessageFromSignal(signal)
		s.Messages <- message
	}
}

// SendMessage() uses the signal-cli command in dbus mode. This is due to
// the method org.asamk.Signal.sendMessage not working.
func (s *Signal) SendMessage(m *Message) {
	// TODO: further investigate using org.asamk.Signal.sendMessage method.
	args := []string{"--dbus", "send", "-m", m.Text, m.PhoneNumber}
	if m.Attachment != "" {
		args = append(args, "-a", m.Attachment)
	}
	cmd := exec.Command("signal-cli", args...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Message sent")
	log.Print(out)
}

// Establish a connection to the org.asamk.Signal DBus interface on the Session Bus.
func connectDBus(signals chan<- *dbus.Signal) (conn *dbus.Conn) {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	if err := verifyConnection(conn); err != nil {
		log.Fatal(err)
	}

	options := dbus.WithMatchSender(SIGNAL_CLI_DBUS_SERVICE)
	err = conn.AddMatchSignal(options)
	if err != nil {
		log.Fatal(err)
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

// TODO: Rename?
// TODO: Text -> Command?
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
func newMessageFromSignal(signal *dbus.Signal) *Message {
	log.Printf("Body: %+v", signal)
	utc, ok := signal.Body[0].(int64)
	if !ok {
		log.Fatalf("Failed to convert time to int64, %v\n", signal.Body[0])
	}
	t := time.Unix(utc, 0)
	phone, ok := signal.Body[1].(string)
	if !ok {
		log.Fatalf("Failed to convert phone number to string, %v\n", signal.Body[1])
	}
	text, ok := signal.Body[3].(string)
	if !ok {
		log.Printf("Failed to convert message text to string, %v\n", signal.Body[3])
		text = ""
	}
	attachment, ok := signal.Body[4].(string)
	if !ok {
		log.Printf("Failed to convert attachment path to string, %v", signal.Body[4])
		attachment = ""
	}
	return NewMessage(t, phone, text, attachment, signal.Name)
}
