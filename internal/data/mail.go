package data

import (
	"context"
	"fmt"
	"kratos_single/internal/conf"
	"net/smtp"

	"github.com/go-kratos/kratos/v2/log"
)

type MailRepo struct {
	log  *log.Helper
	conf *conf.Data
}

func NewMailRepo() *MailRepo {
	return &MailRepo{}
}

func (r *MailRepo) Send(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {

	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	host := "smtp.gmail.com"        // smtp.gmail.com
	port := "587"        // 587
	user := ""        // your@gmail.com
	pass := ""    // app password
	from := ""        // your@gmail.com

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body,
	)

	auth := smtp.PlainAuth("", user, pass, host)

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", host, port),
		auth,
		from,
		[]string{to},
		msg,
	)
}