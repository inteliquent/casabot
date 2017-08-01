package main

import (
	"github.com/inteliquent/casatunes"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

func casa_PlayAlbum(slack_api *slack.Client, command *casabot_command) {
	CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
	casa_api := casatunes.New(CASA_ENDPOINT)

	message_parameters := slack.NewPostMessageParameters()
	message_parameters.AsUser = true

	channelID := command.ChannelID

	search_text := strings.Join(command.Args, " ")

	media_collection, err := casa_api.MediaSearchMC("070cb4f4bdba71f00352207de5049ddd", search_text)

	if err != nil {
		log.Fatal(err)
	}

	media_item := casatunes.RESTMediaItem{}
	for _, item := range media_collection.MediaItems {
		if item.GroupName == "Albums" || item.GroupName == "Playlists" {
			media_item = item
			break
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	message := ""

	if media_item.Title != "" {
		// Add the album / playlist & start playing
		_, err := casa_api.MediaSourcesPlayMedia("0", media_item.ID, "playnow")
		if err != nil {
			log.Fatal(err)
		}
		message = "Putting on album, '" + media_item.Title + "'"
	} else {
		message = "I couldn't find anything for '" + search_text + "'"
	}

	slack_api.PostMessage(
		channelID,
		message,
		message_parameters,
	)
}
