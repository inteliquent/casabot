package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

func main() {
	SLACK_TOKEN := os.Getenv("SLACK_TOKEN")
	slack_api := slack.New(SLACK_TOKEN)

	slack_message_parameters := slack.NewPostMessageParameters()
	slack_message_parameters.AsUser = true

	// Start the boombox thread
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
			casabot_command := parseCommand(ev)
			if casabot_command != nil && casabot_command.UserID == rtm.GetInfo().User.ID {
				log.Printf("Command: %+v", casabot_command)
				switch strings.ToLower(casabot_command.Command) {
				case "nowplaying":
					casa_NowPlaying(slack_api, casabot_command)
				case "play":
					switch casabot_command.verb([]string{"song", "album"}) {
					case "song":
						log.Printf("playing song %s", casabot_command.Args)
						//casa_PlaySong(slack_api, casabot_command)
					case "album":
						log.Printf("playing album %s", casabot_command.Args)
						//casa_PlayAlbum(slack_api, casabot_command)
					default:
						log.Printf("default play song: %s", casabot_command.Args)
						//casa_PlaySong(slack_api, casabot_command)
					}
				case "search":
					switch casabot_command.verb([]string{"song", "album"}) {
					default:
						casa_SearchSong(slack_api, casabot_command)
					}
				}
			}

			if regexp_playeraction.MatchString(ev.Text) {
				casa_PlayerAction(slack_api, ev)
			}

			if regexp_boombox.MatchString(ev.Text) {
				user_input := regexp_boombox.FindStringSubmatch(ev.Text)[1]
				switch user_input {
				case "start":
					boombox.addChannel(ev.Channel, slack_api)
				case "stop":
					boombox.removeChannel(ev.Channel, slack_api)
				}
			}

		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())

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
