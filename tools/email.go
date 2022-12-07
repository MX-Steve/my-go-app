package tools

import (
	"log"
	"net/smtp"
	"strings"
)

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	log.Printf("Send to %s", to)
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/html;charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain;charset=UTF-8"
	}
	body = strings.TrimSpace(body)
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
