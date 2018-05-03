package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
	"fmt"
)

var (
	helpText []string = []string{
		"Available commands are:",
		"	- `nowplaying | now playing` // Display information on the current song",
		"	- `play (song|album) <search text>`	// Play the first song/album result for _search text_ & add it to the queue",
		"	- `listen to (song|album) <search text>`	// Play the first song/album result for _search text_ & add it to the queue",
		"	- `search (song|album) <search text>` // Print the first song/album result for _search text_",
		"	- `boombox (start|stop)` // Start or stop displaying the current song in this channel",
	}
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
				switch strings.ToLower(casabot_command.Command) {
				case "nowplaying":
					casa_NowPlaying(slack_api, casabot_command)
				case "now":
					switch casabot_command.verb([]string{"playing"}) {
					case "playing":
						casa_NowPlaying(slack_api, casabot_command)
					default:
						var text string
						for _, line := range helpText {
							text += line + "\n"
						}
						slack_api.PostMessage(ev.Channel, text, slack_message_parameters)
					}
				case "play":
					switch casabot_command.verb([]string{"song", "album"}) {
					case "song":
						casa_PlaySong(slack_api, casabot_command)
					case "album":
						casa_PlayAlbum(slack_api, casabot_command)
					default:
						casa_PlaySong(slack_api, casabot_command)
					}
				case "listen":
					switch casabot_command.verb([]string{"to"}) {
					case "to":
						switch casabot_command.verb([]string{"song", "album"}) {
						case "song":
							casa_PlaySong(slack_api, casabot_command)
						case "album":
							casa_PlayAlbum(slack_api, casabot_command)
						default:
							casa_PlaySong(slack_api, casabot_command)
						}
					default:
						var text string
						for _, line := range helpText {
							text += line + "\n"
						}
						slack_api.PostMessage(ev.Channel, text, slack_message_parameters)
					}
				case "search":
					switch casabot_command.verb([]string{"song", "album"}) {
					default:
						casa_SearchSong(slack_api, casabot_command)
					}
				case "boombox":
					switch casabot_command.verb([]string{"start", "stop"}) {
					case "start":
						boombox.addChannel(ev.Channel, slack_api)
					case "stop":
						boombox.removeChannel(ev.Channel, slack_api)
					}
				case "pause", "resume", "next", "previous":
					casa_PlayerAction(slack_api, casabot_command)
				default:
					var text string
					for _, line := range helpText {
						text += line + "\n"
					}
					slack_api.PostMessage(ev.Channel, text, slack_message_parameters)
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
			fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
