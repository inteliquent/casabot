package main

import (
  "log"
  "time"
  "os"
  "github.com/inteliquent/casatunes"
  "github.com/nlopes/slack"
)

var (
  boombox_started bool = false
  bb_quit chan int = make(chan int, 1)
)

func casa_BoomBox(slack_api *slack.Client, ev *slack.MessageEvent) {
  //CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  //casa_api := casatunes.New(CASA_ENDPOINT)

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  user_input := regexp_boombox.FindStringSubmatch(ev.Text)[1]

  checkInterval, err := time.ParseDuration("2s")
  bb_ticker := time.NewTicker(checkInterval)

  if err != nil {
    log.Fatal(err)
  }

  switch user_input {
  case "start":
    if !boombox_started {
      log.Printf("Starting BoomBox in channel [%s]", ev.Channel)
      slack_api.PostMessage(
        ev.Channel,
        "OK - I'll post every song change in this channel.",
        message_parameters,
      )
      go boomBox(bb_ticker.C, bb_quit, slack_api, ev)
      boombox_started = true
    }
  case "stop":
    if boombox_started {
      log.Printf("Stopping BoomBox in channel [%s]", ev.Channel)
      slack_api.PostMessage(
        ev.Channel,
        "OK - I'll stop posting song changes to this channel.",
        message_parameters,
      )
      bb_quit <- 0
      boombox_started = false
    }
  }
}

func boomBox(tick <-chan time.Time, quit chan int, slack_api *slack.Client, ev *slack.MessageEvent) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  casa_api := casatunes.New(CASA_ENDPOINT)

  np_check := &casatunes.RESTNowPlayingMediaItem{}
  np_current, err := casa_api.NowPlaying("0")
  if err != nil {
    log.Fatal(err)
  }

  for {
    select {
    case <-tick:
      np_check, err = casa_api.NowPlaying("0")
      if err != nil {
        log.Fatal(err)
      }
      if np_check.CurrSong.ID != np_current.CurrSong.ID {
        casa_NowPlaying(slack_api, ev)
        np_current = np_check
      }
    case <-quit:
      return
    }
  }
}
