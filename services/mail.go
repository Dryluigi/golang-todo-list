package services

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail() {
	mailer := gomail.NewMessage()

	mailer.SetHeader("From", "mario@gmail.com")
	mailer.SetHeader("To", "dmowqodm@gmail.com")
	mailer.SetAddressHeader("Cc", "qwqdw@gmail.com", "Joko")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", "Hello, <b>have a nice day</b>")

	dialer := gomail.NewDialer(
		"smtp.mailtrap.io",
		587,
		"318b88c293edce",
		"a346f05b8b93f0",
	)

	err := dialer.DialAndSend(mailer)

	if err != nil {
		log.Println(err.Error())
	}
}
