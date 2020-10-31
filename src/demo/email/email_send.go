package email

import (
	"gopkg.in/gomail.v2"
)

func sendMail(email string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "alex@mail.brandu.vip")
	m.SetHeader("To", "china.codehome@gmail.com")
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "user", "123456")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
