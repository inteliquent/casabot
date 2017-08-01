package main

import (
	"fmt"
	"github.com/inteliquent/casatunes"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

var (
	message string
)

func casa_PlayerAction(slack_api *slack.Client, command *casabot_command) {
	CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
	casa_api := casatunes.New(CASA_ENDPOINT)

	message_parameters := slack.NewPostMessageParameters()
	message_parameters.AsUser = true

	channelID := command.ChannelID

	playerAction := command.Command

	if playerAction == "resume" {
		playerAction = "play"
	}

	err := casa_api.SourcesPlayerAction("0", playerAction)

	if err != nil {
		log.Println(err)
		message = fmt.Sprintf("Oops! %s", err)
	} else {
		if err != nil {
			log.Fatal(err)
		}

		switch strings.ToLower(playerAction) {
		case "play":
			message = "resumed playback!"

		case "pause":
			message = "paused playback!"

		case "previous":
			message = "_*BACK!!*_"

		case "next":
			message = "_*NEXT!!*_"
		}
	}

	slack_api.PostMessage(
		channelID,
		message,
		message_parameters,
	)
}
