package main

import "regexp"

var (
  regexp_nowplaying *regexp.Regexp = regexp.MustCompile(
    `(?i)what(.?s)? (is |song is )?(this|the|playing)( song| now)? ?\??$`,
  )
  regexp_playsong *regexp.Regexp = regexp.MustCompile(
    `^(?i)(?:@casabot )?listen to (.+)$`,
  )
)
