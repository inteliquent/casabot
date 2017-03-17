package main

import (
  "os"
  "github.com/nlopes/slack"
  "github.com/inteliquent/casatunes"
  "log"
)

func casa_PlaySong(slack_api *slack.Client, ev *slack.MessageEvent) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  casa_api := casatunes.New(CASA_ENDPOINT)

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  channelID := ev.Channel

  search_text := regexp_playsong.FindStringSubmatch(ev.Text)[1]

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

  user, err := slack_api.GetUserInfo(ev.User)

  if err != nil {
    log.Fatal(err)
  }

  if media_item.Title != "" {
    attachment := slack.Attachment{
      Title: user.RealName + " is now playing",
      TitleLink: CASA_ENDPOINT,
      ThumbURL: media_item.ArtworkURI,
      Color: "#aeffa0",
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

    _, err := casa_api.MediaSourcesPlayMedia("0", media_item.ID)

    if err != nil {
      log.Fatal(err)
    }

    slack_api.PostMessage(
      channelID,
      "",
      message_parameters,
    )
  }
}
