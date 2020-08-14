# Signal Assistant
This project is in its infancy.

## Dependencies
- signal-cli https://github.com/AsamK/signal-cli
- godbus/dbus https://github.com/godbus/dbus

## Run
This application expects the `OWNER_NUMBER` env variable to be set to a valid Signal
Phone Number. Must have the country code.

```
OWNER_NUMBER=+11234567890 go run main.go
```
