package main

import (
	"log"
	"os"
	"fmt"
	"golang.org/x/net/websocket"
)

type responseSelf struct {
    Id string `json:"id"`
}

type responseRtmStart struct {
    Ok    bool         `json:"ok"`
    Error string       `json:"error"`
    Url   string       `json:"url"`
    Self  responseSelf `json:"self"`
}

type Message struct {
  Id      uint64 `json:"id"`
  Type    string `json:"type"`
  Channel string `json:"channel"`
  Text    string `json:"text"`
}

var (
  respObj responseRtmStart
  m Message
  counter uint64
)
/*
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
*/

func main() {
  slackKey := os.Getenv("SLACK_SECRET_KEY")
  url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", slackKey)
  resp, err := http.Get(url)
  if err != nil {
    log.Fatalln("Fail to get ws url")
    return 1
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatalln("Fail to read get ws url response body")
    return 1
  }
  err = json.Unmarshal(body, &respObj)
  if err != nil {
    log.Fatalln("Fail to parse json get ws url response")
    return 1
  }
  if !respObj.Ok {
		err = log.Fatalln("Slack error: %s", respObj.Error)
		return
	}
  
  wsurl := respObj.Url
	id := respObj.Self.Id
  log.Printf("wsurl: %s, id: %s", wsurl, id)
  
  ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
  if err != nil {
    log.Fatalln("Fail to dial websocket")
    return 1
  }
  
  
  for {
    err := websocket.JSON.Receive(ws, &m)
    m.Id = atomic.AddUint64(&counter, 1)
    err := websocket.JSON.Send(ws, m)
  }
  
	//api := slack.New(slackKey)
	//os.Exit(run(api))
}
