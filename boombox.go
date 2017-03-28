package main

import (
  "errors"
  "log"
  "time"
  "os"
  "github.com/inteliquent/casatunes"
  "github.com/nlopes/slack"
)

type boomBox struct {
  Channels []string
}

// Add a Slack channel to the boomBox object. Return an error
// if the boomBox object already contains that channel.
func (boombox *boomBox) addChannel(channel string) (error) {
  if len(boombox.Channels) > 0 {
    for _, c := range(boombox.Channels) {
      if c == channel {
        return errors.New("Channel already in BoomBox!")
      }
    }
  }
  boombox.Channels = append(boombox.Channels, channel)
  return nil
}

// Remove a Slack channel from the boomBox object. Return an error
// if the channel does not exist within the boomBox object.
func (boombox *boomBox) removeChannel(channel string) (error) {
  if len(boombox.Channels) > 0 {
    for i, c := range(boombox.Channels) {
      if c == channel {
        boombox.Channels = append(boombox.Channels[:i], boombox.Channels[i+1:]...)
        return nil
      }
    }
  }
  return errors.New("Channel does not exist in BoomBox!")
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
