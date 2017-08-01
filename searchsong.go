package main

import (
	"github.com/inteliquent/casatunes"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

func casa_SearchSong(slack_api *slack.Client, command *casabot_command) {
	CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
	casa_api := casatunes.New(CASA_ENDPOINT)

	message_parameters := slack.NewPostMessageParameters()
	message_parameters.AsUser = true

	channelID := command.ChannelID

	search_text := strings.Join(command.Args, " ")

	slack_attachment_fields := []slack.AttachmentField{}

	media_collection, err := casa_api.MediaSearchMC("070cb4f4bdba71f00352207de5049ddd", search_text)

	if err != nil {
		log.Fatal(err)
	}

	media_item := casatunes.RESTMediaItem{}
	switch command.Verb {
	case "album", "playlist":
		for _, item := range media_collection.MediaItems {
			if item.GroupName == "Playlists" || item.GroupName == "Albums" {
				media_item = item
				slack_attachment_fields = append(slack_attachment_fields, slack.AttachmentField{Title: "Title", Value: media_item.Title})
				break
			}
		}
	case "title", "song":
		for _, item := range media_collection.MediaItems {
			if item.GroupName == "Tracks" {
				media_item = item
				slack_attachment_fields = append(slack_attachment_fields, slack.AttachmentField{Title: "Title", Value: media_item.Title})
				slack_attachment_fields = append(slack_attachment_fields, slack.AttachmentField{Title: "Artist", Value: media_item.Artists, Short: true})
				slack_attachment_fields = append(slack_attachment_fields, slack.AttachmentField{Title: "Album", Value: media_item.Album, Short: true})
				break
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	message := ""
	if media_item.Title != "" {
		attachment := slack.Attachment{
			Title:     "Search results for '" + search_text + "'",
			TitleLink: CASA_ENDPOINT,
			ThumbURL:  media_item.ArtworkURI,
			Color:     "#009933",
			Fields:    slack_attachment_fields,
		}

		message_parameters.Attachments = []slack.Attachment{attachment}
	} else {
		message = "I couldn't find anything for '" + search_text + "'"
	}
	slack_api.PostMessage(
		channelID,
		message,
		message_parameters,
	)
}
