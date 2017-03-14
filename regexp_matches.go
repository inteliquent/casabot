package main

import "regexp"

var (
  regex_nowplaying *regexp.Regexp = regexp.MustCompile(
    `(?i)what(.?s)? (is |song is )?(this|the|playing)( song| now)? ?\??$`,
  )
)
