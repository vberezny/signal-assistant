# Signal Assistant
A chat bot/assistant for Signal https://www.signal.org/.

## Description
This project uses signal-cli https://github.com/AsamK/signal-cli as a backend. signal-cli is a command
line client and Dbus daemon built using the Java implementation of the Signal protocol.
Using a match signal on the `org.asamk.Signal` DBus service, the assistant listens to all
incoming messages and parses their content for commands. Currently all commands are prefixed with `!`.
In addition to commands, the assistant also supports alerts that can be registered with the application
and performed at pre-determined intervals.

## Setup
To set up the assistant you will need a working phone number for one time sms verification.
If you don't have access to Google Voice I recommend a temporary phone number service that allows you to
recieve verification codes. Temporary numbers work fine for this case because messages sent over the
Signal protocol use phone numbers as user identifiers, nothing more.

### Download/Install signal-cli
Visit https://github.com/AsamK/signal-cli for details.

### Configure signal-cli
1. `signal-cli -u ASSISTANT_NUMBER register`
2. Wait for the verification code to arrive via sms then:
3. `signal-cli -u ASSISTANT_NUMBER verify CODE`
4. Test signal-cli:
```signal-cli -u ASSISTANT_NUMBER send -m "This is a message" OWNER_NUMBER```

### Build and install the binary
1. `go build`
2. `go install`, ensure that your `$GOBIN` is set.
3. Define the following environment variables using your preferred method, as long as they are visible to the `signal-assistant` binary:
    - `ASSISTANT_NUMBER` - the assistant phone number.
    - `OWNER_NUMBER` - the owners phone number. Must be a valid Signal user.
    - `ASSISTANT_FOLDER` - the directory that acts as a shared folder for all attachment related commands.

### Run
Once the binary is installed it should be as simple as running: `signal-assistant`.

## Developing
Either export the env variables beforehand or run it like this:
```
ASSISTANT_NUMBER=+11234567891 OWNER_NUMBER=+11234567890 ASSISTANT_FOLDER=/home/user/shared-directory/ go run main.go
```

## Current Features

### Commands:
`!store` - Store an attachment in a shared folder.

`!get` - Retrieve a file from a shared folder.

`!list` - List the contents of the shared folder.

`!man` - List all commands and their detailed descriptions

### Alerts:
`crypto` - a basic alert mostly for demonstration purposes, checks the price of Bitcoin and Ethereum
and replies if the price action satisfies a set of pre-determined conditions.

## Future Plans
1. Currently the assistant starts signal-cli on the session bus (DBus). Signal-cli supports usage over the system bus
and with some code and assumption changes the assistant can be configured to run there as well.
1. More commands and alerts.
2. More extensible command/alert interfaces to allow for greater customization.
