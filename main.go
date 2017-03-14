package main

import (
  "github.com/nlopes/slack"
  "os"
  "log"
  "fmt"
  "github.com/inteliquent/casatunes"
  "regexp"
)

func main() {
  SLACK_TOKEN := os.Getenv("SLACK_TOKEN")
  CASA_ENDPOINT := os.Getenv("CASA_ENDPOINT")
  slack_api := slack.New(SLACK_TOKEN)
  casa_api := casatunes.New(CASA_ENDPOINT)

  message_parameters := slack.NewPostMessageParameters()
  message_parameters.AsUser = true

  logger := log.New(
    os.Stdout,
    "slack-bot: ",
    log.Lshortfile|log.LstdFlags,
  )

  slack.SetLogger(logger)

  rtm := slack_api.NewRTM()
  go rtm.ManageConnection()

  for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
      regex := regexp.MustCompile(
        `(?i)what(.?s)? (is |song is )?(this|the|playing)( song| now)? ?\??$`,
      )
      if regex.MatchString(ev.Text) {
        channelID := ev.Channel

        nowplaying, err := casa_api.NowPlaying("0")
        if err != nil {
          log.Fatal(err)
        }

        attachment := slack.Attachment{
          Title: "Now Playing",
          TitleLink: CASA_ENDPOINT,
          ThumbURL: nowplaying.CurrSong.ArtworkURI,
          Fields: []slack.AttachmentField{
            slack.AttachmentField{
              Title: "Title",
              Value: nowplaying.CurrSong.Title,
            },
            slack.AttachmentField{
              Title: "Artist",
              Value: nowplaying.CurrSong.Artists,
              Short: true,
            },
            slack.AttachmentField{
              Title: "Album",
              Value: nowplaying.CurrSong.Album,
              Short: true,
            },
          },
        }

        message_parameters.Attachments = []slack.Attachment{attachment}

        slack_api.PostMessage(
          channelID,
          "",
          message_parameters,
        )
      }

    case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

    case *slack.InvalidAuthEvent:
			log.Fatal("Invalid credentials")

    case *slack.ConnectionErrorEvent:
      log.Fatal("Failed to authenticate with Slack API")

    case *slack.ConnectingEvent:
      log.Printf("Connecting to Slack (attempt %d)", ev.Attempt)

    case *slack.HelloEvent:
      log.Printf("Received hello from Slack API.")

    case *slack.ConnectedEvent:
      log.Printf("Connected to %s as %s.", ev.Info.Team.Name, ev.Info.User.Name)

    case *slack.LatencyReport:
      log.Printf("Latency report: %s", ev.Value)

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
