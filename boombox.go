package main

import (
  "log"
  "time"
  "os"
  "github.com/inteliquent/casatunes"
  "github.com/nlopes/slack"
)
func casa_BoomBox(slack_api *slack.Client, ev *slack.MessageEvent) {
  //CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  //casa_api := casatunes.New(CASA_ENDPOINT)

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  userInput := regexp_boombox.FindStringSubmatch(ev.Text)[1]

  checkInterval, err := time.ParseDuration("2s")

  if err != nil {
    log.Fatal(err)
  }

  switch userInput {
  case "start":
    log.Println("Starting BoomBox")
    bbTicker := time.NewTicker(checkInterval)
    bbQuit := make(chan int, 1)
    go boomBox(bbTicker.C, bbQuit, slack_api, ev)
  case "stop":
  }
  /*
  slack_api.PostMessage(
    ev.Channel,
    "",
    message_parameters,
  )
*/
}

func boomBox(tick <-chan time.Time, quit chan int, slack_api *slack.Client, ev *slack.MessageEvent) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  casa_api := casatunes.New(CASA_ENDPOINT)

  npCurrent, err := casa_api.NowPlaying("0")
  if err != nil {
    log.Fatal(err)
  }
  npCheck := &casatunes.RESTNowPlayingMediaItem{}
  for {
    select {
    case <-tick:
      npCheck, err = casa_api.NowPlaying("0")
      if err != nil {
        log.Fatal(err)
      }
      if npCheck.CurrSong.ID != npCurrent.CurrSong.ID {
        casa_NowPlaying(slack_api, ev)
        npCurrent = npCheck
      }
    case <-quit:
      log.Printf("Stopping boombox in channel %s", ev.Channel)
      return
    }
  }
}
