package mail

import (
	"context"
	"time"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/mailgun/mailgun-go/v3"
	log "github.com/sirupsen/logrus"
)

const timeout = time.Second * 20

type EmailSender struct {
	Conf *config.Config
}

func NewEmailSender(conf *config.Config) EmailSender {
	return EmailSender{Conf: conf}
}

func (e EmailSender) Send(from, to, subject, body string) error {
	var err error

	mg := mailgun.NewMailgun(e.Conf.MgDomain, e.Conf.MgAPIKEY)
	mg.SetAPIBase(mailgun.APIBaseEU)

	m := mg.NewMessage(from, subject, "", to)
	m.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	msg, id, err := mg.Send(ctx, m)

	if err != nil {
		log.Errorf("Can't send email, got err %s", err)
	} else {
		log.Debugf("msg %s, id %s", msg, id)
	}

	return err
}
