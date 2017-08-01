package main

import (
	"regexp"
	"log"
	"strings"
)

type casabot_command struct {
	UserID string
	Command string
	Args []string
}

func parseCommand(msg string) *casabot_command {
	command_regexp := regexp.MustCompile(`^(?i)<@([A-Z0-9]+)> +([a-z]+) (.*)$`)

	if command_regexp.MatchString(msg) {
		log.Printf("Command: %s\n", msg)

		casabot_command := casabot_command{}
		command_parsed := command_regexp.FindStringSubmatch(msg)

		casabot_command.UserID = command_parsed[1]
		casabot_command.Command = command_parsed[2]
		casabot_command.Args = strings.Split(command_parsed[3], " ")

		return &casabot_command
	}

	return nil
}