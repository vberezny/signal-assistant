package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	SEND_MESSAGE_ADDRESS = "org.asamk.Signal.sendMessage"
	SIGNAL_CLI_DBUS_NAME = "org.asamk.Signal"
)

// TODO: probably rename.
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

func newMessageFromSignal(signal *dbus.Signal) *Message {
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

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := checkConnection(conn); err != nil {
		log.Fatal(err)
	}

	options := dbus.WithMatchSender(SIGNAL_CLI_DBUS_NAME)
	err = conn.AddMatchSignal(options)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *dbus.Signal, 10)
	conn.Signal(ch)
	for signal := range ch {
		message := newMessageFromSignal(signal)
		if message.Text == "command" {
			fmt.Println("valid command")
			obj := conn.Object(SIGNAL_CLI_DBUS_NAME, "/org/asamk/Signal")
			fmt.Println(obj.Destination())
			call := obj.Call(SEND_MESSAGE_ADDRESS, 0, "text")
			if call.Err != nil {
				log.Fatal(call.Err)
			}
		} else {
			fmt.Println("invalid command")
		}
	}
}

func checkConnection(conn *dbus.Conn) (err error) {
	var s []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return
	}
	found := false
	for _, v := range s {
		if v == SIGNAL_CLI_DBUS_NAME {
			found = true
			break
		}
	}
	if !found {
		err = errors.New("signal-cli connection not found on Dbus.")
	}
	return
}
