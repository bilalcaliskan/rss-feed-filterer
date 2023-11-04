package main

import (
	"os"

	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

//func main() {
//	cfg, err := config.ReadConfig("configs/sample_config.yaml")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(cfg.Email)
//
//	// Initialize the AWS SDK's SES client.
//	awsCfg, err := aws.CreateConfig(cfg.Email.AccessKey, "9I9Am88fwZ6k/iNhdRO3jVO/Q1awmternymeS4Yy", "us-east-1")
//	if err != nil {
//		fmt.Println("Failed to load AWS SDK config:", err)
//		return
//	}
//
//	sesClient := ses.NewFromConfig(awsCfg)
//
//	// Initialize the SES sender with the client and sender email.
//	fromEmail := "info@giantrooster.tech"
//	sesSender := internalses.NewSESSender(sesClient, fromEmail)
//
//	// Create an email announcer with SES sender.
//	// Assuming the 'to' address is "receiver@example.com".
//	emailAnnouncer := email.NewEmailAnnouncer(sesSender, "bilalcaliskan@protonmail.com")
//
//	emailPayload := &email.EmailPayload{
//		Subject: "Test Email",
//		Content: "This is the content of the email.",
//		Options: []email.SendOption{
//			email.WithCC("jsagredo@protonmail.com"),
//			email.WithBCC("jsagredo@protonmail.com"),
//			//email.WithBCC("bcc@example.com"),
//		},
//	}
//	err = emailAnnouncer.Notify(emailPayload)
//	if err != nil {
//		fmt.Println("Failed to send announcement:", err)
//	} else {
//		fmt.Println("Announcement sent successfully!")
//	}
//}
