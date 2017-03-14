package main

import (
  "github.com/nlopes/slack"
  "os"
  "log"
  "fmt"
  "regexp"
)

func main() {
  SLACK_TOKEN := os.Getenv("SLACK_TOKEN")
  slack_api := slack.New(SLACK_TOKEN)

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
        nowPlaying(slack_api, ev)
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
