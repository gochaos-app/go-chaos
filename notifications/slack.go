package notifications

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func SlackNotification(slackChannel []string, body string) {
	fmt.Println("testing")
	token := os.Getenv("SLACK_AUTH_TOKEN")
	if token == "" {
		log.Println("Cannot find variable SLACK_AUTH_TOKEN")
		return
	}
	slackClient := slack.New(token)

	attachment := slack.Attachment{
		Pretext: "go-chaos",
		Text:    body,
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}

	for i := 0; i < len(slackChannel); i++ {

		_, timestamp, err := slackClient.PostMessage(
			slackChannel[i],

			slack.MsgOptionAttachments(attachment),
		)
		if err != nil {
			log.Println("error posting slack message via slack bot", err)
			return
		}
		log.Println("message sent to slack channel", slackChannel[i], "at", timestamp)

	}

}
