package main

import (
	"strings"
)

// Command bot command structure
type Command struct {
	name string
	args []string
}

// ParseCommand parsing text commands
func ParseCommand(message string) Command {
	var command Command
	trimmedMessage := strings.TrimSpace(message)
	if !strings.HasPrefix(trimmedMessage, "/") {
		return command
	}

	splited := strings.Split(trimmedMessage, " ")
	command.name = splited[0][1:]

	for _, elem := range splited[1:] {
		if elem != "" {
			command.args = append(command.args, elem)
		}
	}

	return command
}
