package assistant

import "fmt"

const (
	STORE = "!store"
	GET   = "!get"
	LIST  = "!list"
	MAN   = "!man"

	STORE_DESCRIPTION = "Store an attachment in the shared folder.\n" +
		"Must provide 2 arguments: name and file type, separated by spaces.\n" +
		"Currently only supports 1 attachment at a time."
	GET_DESCRIPTION = "Retrieve a file from the shared folder.\n" +
		"Must provide at least 1 argument: the complete file name.\n" +
		"If more than 1 argument is provided the assistant will attempt\n" +
		"to return all files listed. Arguments must be space separated.\n" +
		"Currently multiple attachments can only be returned if they are\n" +
		"of a similar file type (such as png and jpg, png and pdf will only return the first one)."
	LIST_DESCRIPTION = "List all files in the shared folder."
	MAN_DESCRIPTION  = "List all commands and their descriptions."
)

var commands = map[string]string{
	STORE: STORE_DESCRIPTION,
	GET:   GET_DESCRIPTION,
	LIST:  LIST_DESCRIPTION,
	MAN:   MAN_DESCRIPTION,
}

type command struct {
	cmd         string
	description string
	args        []string
}

func getAllCommands() (allCommands []command) {
	for cmd, description := range commands {
		c := command{
			cmd:         cmd,
			description: description,
		}
		allCommands = append(allCommands, c)
	}
	return
}

func getCommandManual() (manual string) {
	for cmd, description := range commands {
		manual += fmt.Sprintf("%v ---\n %v\n", cmd, description)
	}
	return
}
