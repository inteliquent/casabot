package main

import "regexp"

var (
  regexp_nowplaying *regexp.Regexp = regexp.MustCompile(
    `(?i)what(.?s)? (is |song is )?(this|the|playing)( song| now)? ?\??$`,
  )
  regexp_playsong *regexp.Regexp = regexp.MustCompile(
    `(?i)listen to (.+)$`,
  )
  regexp_playeraction *regexp.Regexp = regexp.MustCompile(
    `(?i)(play|pause|next|previous) (?:song|track|music)$`,
  )
  regexp_boombox *regexp.Regexp = regexp.MustCompile(
    `^(?i)(start|stop) boombox$`,
  )
)
