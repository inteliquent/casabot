package main

import "regexp"

var (
  regexp_nowplaying *regexp.Regexp = regexp.MustCompile(
    `(?i)what(.?s)? (is |song is )?(this|the|playing)( song| now)? ?\??$`,
  )
  regexp_playsong *regexp.Regexp = regexp.MustCompile(
    `(?i)listen to (.+)$`,
  )
  regexp_playalbum *regexp.Regexp = regexp.MustCompile(
    `(?i)put on (?:album|playlist) (.+)$`,
  )
  regexp_searchsong *regexp.Regexp = regexp.MustCompile(
    `(?i)search(?: for)? (song|title|playlist|album) (.+)$`,
  )
  regexp_playeraction *regexp.Regexp = regexp.MustCompile(
    `(?i)(play|pause|next|previous|stop) (?:song|track|music)$`,
  )
  regexp_boombox *regexp.Regexp = regexp.MustCompile(
    `^(?i)(start|stop) boombox$`,
  )
)
