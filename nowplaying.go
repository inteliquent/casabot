package main

import (
  "log"
  "os"
  "github.com/nlopes/slack"
  "github.com/inteliquent/casatunes"
)

func casa_NowPlaying(slack_api *slack.Client, ev *slack.MessageEvent) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  casa_api := casatunes.New(CASA_ENDPOINT)
  channelID := ev.Channel

  nowplaying, err := casa_api.NowPlaying("0")
  if err != nil {
    log.Fatal(err)
  }

  attachment := slack.Attachment{
    Title: "Now Playing",
    TitleLink: CASA_ENDPOINT,
    ThumbURL: nowplaying.CurrSong.ArtworkURI,
    Fields: []slack.AttachmentField{
      slack.AttachmentField{
        Title: "Title",
        Value: nowplaying.CurrSong.Title,
      },
      slack.AttachmentField{
        Title: "Artist",
        Value: nowplaying.CurrSong.Artists,
        Short: true,
      },
      slack.AttachmentField{
        Title: "Album",
        Value: nowplaying.CurrSong.Album,
        Short: true,
      },
    },
  }

  message_parameters.Attachments = []slack.Attachment{attachment}

  slack_api.PostMessage(
    channelID,
    "",
    message_parameters,
  )
}
