package notifications

import (
	"log"
	"net/smtp"
	"os"
)

func GmailNotification(email []string, body string, from string) {
	pass := os.Getenv("GMAIL_APP_TOKEN")

	for i := 0; i < len(email); i++ {
		subject := "Subject: Go-chaos execution notification \n"
		message := subject + body
		err := smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from, []string{email[i]}, []byte(message))
		if err != nil {
			log.Printf("smtp error: %s", err)
			return
		}
		log.Println("Sending notification to:", email[i])
	}
}
