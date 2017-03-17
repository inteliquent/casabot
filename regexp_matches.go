package main

import "regexp"

var (
  regexp_nowplaying *regexp.Regexp = regexp.MustCompile(
    `(?i)what(.?s)? (is |song is )?(this|the|playing)( song| now)? ?\??$`,
  )
  regexp_playsong *regexp.Regexp = regexp.MustCompile(
    `^(?i)(?:@casabot )?listen to (.+)$`,
  )
  regexp_playeraction *regexp.Regexp = regexp.MustCompile(
    `^(?i)(?:@casabot )?(play|pause|next|previous) (?:song|track|music)$`,
  )
)
