package email

import (
	"fmt"
	"github.com/badoux/checkmail"
)

func MailVerify() {
	var (
		serverHostName    = "smtp.gmail.com"           // set your SMTP server here
		serverMailAddress = "china.codehome@gmail.com" // set your valid mail address here
	)
	err := checkmail.ValidateHostAndUser(serverHostName, serverMailAddress, "129083726@gmail.com")
	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		fmt.Printf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
	}
}
