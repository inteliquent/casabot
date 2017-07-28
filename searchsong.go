package main

import (
  "os"
  "github.com/nlopes/slack"
  "github.com/inteliquent/casatunes"
  "log"
)

func casa_SearchSong(slack_api *slack.Client, ev *slack.MessageEvent) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  casa_api := casatunes.New(CASA_ENDPOINT)

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  channelID := ev.Channel

  search_text := regexp_searchsong.FindStringSubmatch(ev.Text)[1]

  media_collection, err := casa_api.MediaSearchMC("070cb4f4bdba71f00352207de5049ddd", search_text)

  if err != nil {
    log.Fatal(err)
  }

  media_item := casatunes.RESTMediaItem{}
  for _, item := range media_collection.MediaItems {
    if item.GroupName == "Tracks" {
      media_item = item
      break
    }
  }

  if err != nil {
    log.Fatal(err)
  }

  message := ""
  if media_item.Title != "" {
    attachment := slack.Attachment{
      Title: "Search results for '" + search_text + "'",
      TitleLink: CASA_ENDPOINT,
      ThumbURL: media_item.ArtworkURI,
      Color: "#009933",
      Fields: []slack.AttachmentField{
        slack.AttachmentField{
          Title: "Title",
          Value: media_item.Title,
        },
        slack.AttachmentField{
          Title: "Artist",
          Value: media_item.Artists,
          Short: true,
        },
        slack.AttachmentField{
          Title: "Album",
          Value: media_item.Album,
          Short: true,
        },
      },
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
