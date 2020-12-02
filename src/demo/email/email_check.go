package email

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"strings"
)

func email_check_single(email string) (bool, string) {
	emailPrefixHost := strings.Split(email, "@")
	domain := emailPrefixHost[1]

	MXs, err := net.LookupMX(domain)
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}

	ipList := []string{}
	for _, MX := range MXs {
		//fmt.Println(MX.Host)
		ipList = append(ipList, MX.Host)
	}

	ip := ipList[rand.Intn(len(ipList))]

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(fmt.Sprintf("%s:25", ip))
	if err != nil {
		log.Fatal(err)
	}

	helloServer := "joananne.info"
	c.Hello(helloServer)
	// Set the sender and recipient first
	if err := c.Mail(fmt.Sprintf("support@%s", helloServer)); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(email); err != nil {
		log.Fatal(err)
	}
	//c.Verify("postmaster11@joananne.info")

	// Send the email body.
	//wc, err := c.Data()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = fmt.Fprintf(wc, "This is the email body")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = wc.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		//log.Fatal(err)
		return false, err.Error()
	}
	return true, "OK"
}

func Email_check() {
	var emailList []string
	emailList = append(emailList, "2548252929@qq.com")
	emailList = append(emailList, "qmxgame@yahoo.com")
	emailList = append(emailList, "china.martin@hotmail.com")
	emailList = append(emailList, "postmaster@joananne.info")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	for _, email := range emailList {
		status, msg := email_check_single(email)
		log.Printf("email:%s, status %t , msg : %s", email, status, msg)
	}

}
