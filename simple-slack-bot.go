package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func run(api *slack.Client) int {
	log.Print("PORT"+os.Getenv("PORT"))


	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("Hello Event")

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world: "+ev.Msg.Text, ev.Channel))

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1

			}
		}
	}
}

func main() {
	slackKey := os.Getenv("SLACK_SECRET_KEY")
	api := slack.New(slackKey)
	os.Exit(run(api))
}
