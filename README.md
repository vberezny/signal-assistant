# Signal Assistant
This project is in its infancy.

## Dependencies
- signal-cli https://github.com/AsamK/signal-cli
- godbus/dbus https://github.com/godbus/dbus

## Run
This application expects the `OWNER_NUMBER` env variable to be set to a valid Signal
Phone Number including the country code. `ASSISTANT_FOLDER` must be set to a valid directory 
to be used as a shared folder by the assistant to store and retrieve attachments. 
```
OWNER_NUMBER=+11234567890 ASSISTANT_FOLDER=/home/user/shared-directory/ go run main.go
```
