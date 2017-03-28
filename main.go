package main

import (
  "github.com/nlopes/slack"
  "os"
  "log"
  "fmt"
)

func main() {
  SLACK_TOKEN := os.Getenv("SLACK_TOKEN")
  slack_api := slack.New(SLACK_TOKEN)

  slack_message_parameters := slack.NewPostMessageParameters()
  slack_message_parameters.AsUser = true

  boombox := boomBox{}

  go boombox.start(slack_api)

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
      if regexp_nowplaying.MatchString(ev.Text) {
        casa_NowPlaying(slack_api, ev.Channel)
      }

      if regexp_playsong.MatchString(ev.Text) {
        casa_PlaySong(slack_api, ev)
      }

      if regexp_playeraction.MatchString(ev.Text) {
        casa_PlayerAction(slack_api, ev)
      }

      if regexp_boombox.MatchString(ev.Text) {
        user_input := regexp_boombox.FindStringSubmatch(ev.Text)[1]
        switch user_input {
        case "start":
          err := boombox.addChannel(ev.Channel)
          if err != nil {
            log.Printf(
              "Failed to start BoomBox in channel [%s]: %s",
              ev.Channel,
              err,
            )
            slack_api.PostMessage(
              ev.Channel,
              fmt.Sprint(err),
              slack_message_parameters,
            )
          } else {
            log.Printf("BoomBox started in channel [%s]", ev.Channel)
            slack_api.PostMessage(
              ev.Channel,
              "BoomBox started! I'll report future song changes to this channel.",
              slack_message_parameters,
            )
          }
        case "stop":
          err := boombox.removeChannel(ev.Channel)
          if err != nil {
            log.Printf(
              "Failed to stop BoomBox in channel [%s]: %s",
              ev.Channel,
              err,
            )
            slack_api.PostMessage(
              ev.Channel,
              fmt.Sprint(err),
              slack_message_parameters,
            )
          } else {
            log.Printf("BoomBox stopped in channel [%s]", ev.Channel)
            slack_api.PostMessage(
              ev.Channel,
              "BoomBox stopped! I'll stop reporting song changes to this channel.",
              slack_message_parameters,
            )
          }
        }
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
