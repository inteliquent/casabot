package main

import (
  "errors"
  "log"
  "time"
  "os"
  "github.com/inteliquent/casatunes"
  "github.com/nlopes/slack"
)

func (boombox *boomBox) addChannel(channel string) (error) {
  log.Println(boombox.Channels)
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

func (boombox *boomBox) removeChannel(channel string) (error) {
  log.Println(boombox.Channels)
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

type boomBox struct {
  Channels []string
}

func (boombox *boomBox) goBoomBox(slack_api *slack.Client) {
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
