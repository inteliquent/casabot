package main

import (
  "os"
  "github.com/nlopes/slack"
  "github.com/inteliquent/casatunes"
  "log"
  "fmt"
  "strings"
)

var (
  message string
)

func casa_PlayerAction(slack_api *slack.Client, ev *slack.MessageEvent) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  casa_api := casatunes.New(CASA_ENDPOINT)

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  channelID := ev.Channel

  playerAction := regexp_playeraction.FindStringSubmatch(ev.Text)[1]

  err := casa_api.SourcesPlayerAction("0", playerAction)

  if err != nil {
    log.Println(err)
    message = fmt.Sprintf("Oops! %s", err)
  } else {
    user, err := slack_api.GetUserInfo(ev.User)

    if err != nil {
      log.Fatal(err)
    }

    switch strings.ToLower(playerAction) {
      case "play":
        message = fmt.Sprintf("%s has resumed playback!", user.RealName)

      case "pause":
        message = fmt.Sprintf("%s has paused playback!", user.RealName)

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
