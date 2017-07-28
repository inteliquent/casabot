package main

import (
  "errors"
  "log"
  "time"
  "os"
  "github.com/inteliquent/casatunes"
  "github.com/nlopes/slack"
  "fmt"
)

type boomBox struct {
  Channels []string
}

// Add a Slack channel to the boomBox object. Return an error
// if the boomBox object already contains that channel.
func (boombox *boomBox) addChannel(channel string, slack_api *slack.Client) {
  var err error
  slack_message_parameters := slack.NewPostMessageParameters()
  slack_message_parameters.AsUser = true

  if len(boombox.Channels) > 0 {
    for _, c := range(boombox.Channels) {
      if c == channel {
        err = errors.New("Channel already in BoomBox!")
      }
    }
  }

  if err != nil {
    log.Printf(
      "Failed to start BoomBox in channel [%s]: %s",
      channel,
      err,
    )
    slack_api.PostMessage(
      channel,
      fmt.Sprint(err),
      slack_message_parameters,
    )
  } else {
    boombox.Channels = append(boombox.Channels, channel)
    log.Printf("BoomBox started in channel [%s]", channel)
    slack_api.PostMessage(
      channel,
      "BoomBox started! I'll report future song changes to this channel.",
      slack_message_parameters,
    )
  }
}

// Remove a Slack channel from the boomBox object. Return an error
// if the channel does not exist within the boomBox object.
func (boombox *boomBox) removeChannel(channel string, slack_api *slack.Client) {
  slack_message_parameters := slack.NewPostMessageParameters()
  slack_message_parameters.AsUser = true

  if len(boombox.Channels) > 0 {
    for i, c := range(boombox.Channels) {
      if c == channel {
        boombox.Channels = append(boombox.Channels[:i], boombox.Channels[i+1:]...)
        log.Printf("BoomBox stopped in channel [%s]", channel)
        slack_api.PostMessage(
          channel,
          "BoomBox stopped! I'll stop reporting song changes to this channel.",
          slack_message_parameters,
        )
        return
      }
    }
  }
  log.Printf(
    "Failed to stop BoomBox in channel [%s]: %s",
    channel,
    "Channel does not exist in BoomBox!",
  )
  slack_api.PostMessage(
    channel,
    fmt.Sprint("Channel does not exist in BoomBox!"),
    slack_message_parameters,
  )
}

// This function is meant to run in a go routine. It will
// poll the Casatunes API every 2 seconds, and report to
// any channels in the boomBox object when the song changes
func (boombox *boomBox) start(slack_api *slack.Client) {
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  casa_api := casatunes.New(CASA_ENDPOINT)

  check_interval, err := time.ParseDuration("2s")
  if err != nil {
    log.Fatal(err)
  }
  tick := time.NewTicker(check_interval).C

  np_check := &casatunes.RESTNowPlayingMediaItem{}
  np_current, err := casa_api.NowPlaying("0")
  if err != nil {
    log.Fatal(err)
  }

  for {
    select {
    case <-tick:
      // on every tick, poll casatunes for the current song.
      // if the current song (np_current) is different than the
      // polled song (np_check), chage np_current & post to
      // configured channels.
      np_check, err = casa_api.NowPlaying("0")
      if err != nil {
        log.Fatal(err)
      }
      if np_check.CurrSong.ID != np_current.CurrSong.ID {
        if len(boombox.Channels) > 0 {
          for _, channel := range(boombox.Channels) {
            casa_NowPlaying(slack_api, channel)
          }
        }
        np_current = np_check
      }
    }
  }
}
