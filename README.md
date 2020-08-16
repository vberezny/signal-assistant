# Signal Assistant
A chat bot/assistant for Signal https://www.signal.org/.

## Description
This project uses signal-cli https://github.com/AsamK/signal-cli as a backend. signal-cli is a command
line client and Dbus daemon built using the Java implementation of the Signal protocol.
Using a match signal on the `org.asamk.Signal` DBus service, the assistant listens to all
incoming messages and parses their content for commands. Currently all commands are prefixed with `!`.
In addition to commands, the assistant also supports alerts that can be registered with the application
and performed at pre-determined intervals.

### Current Features

#### Commands:
`!store` - Store an attachment in a shared folder.

`!get` - Retrieve a file from a shared folder.

`!list` - List the contents of the shared folder.

`!man` - List all commands and their detailed descriptions

#### Alerts:
`crypto` - a basic alert mostly for demonstration purposes, checks the price of Bitcoin and Ethereum
and replies if the price action satisfies a set of pre-determined conditions.

### Future Plans
1. More commands and alerts
2. More extensible command/alert interfacesto allow for greater customization

## Setup
TODO

## Developing
This application expects the `OWNER_NUMBER` env variable to be set to a valid Signal
Phone Number including the country code. `ASSISTANT_FOLDER` must be set to a valid directory 
to be used as a shared folder by the assistant to store and retrieve attachments. 
```
OWNER_NUMBER=+11234567890 ASSISTANT_FOLDER=/home/user/shared-directory/ go run main.go
```
