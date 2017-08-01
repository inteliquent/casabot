package main

import (
	"github.com/nlopes/slack"
	"regexp"
	"strings"
)

type casabot_command struct {
	UserID    string
	Command   string
	ChannelID string
	Args      []string
	Verb      string
}

// Shifts the first argument off the Args list
func (command *casabot_command) verb(verbs []string) string {
	for _, verb := range verbs {
		if strings.ToLower(command.Args[0]) == verb {
			command.Args = command.Args[1:]
			command.Verb = verb
			return verb
		}
	}
	return ""
}

func parseCommand(ev *slack.MessageEvent) *casabot_command {
	command_regexp := regexp.MustCompile(`^(?i)<@([A-Z0-9]+)> +([a-z]+) ?(.*)?$`)

	if command_regexp.MatchString(ev.Msg.Text) {
		command_parsed := command_regexp.FindStringSubmatch(ev.Msg.Text)

		casabot_command := casabot_command{}

		casabot_command.UserID = command_parsed[1]
		casabot_command.Command = command_parsed[2]
		casabot_command.Args = strings.Split(command_parsed[3], " ")
		casabot_command.ChannelID = ev.Channel

		return &casabot_command
	}
	return nil
}
