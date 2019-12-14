package mail

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	log "github.com/sirupsen/logrus"
)

const timeout = time.Second * 20

func SendEmail(mgDomain, mgApiKey, from, to, subject, body string) error {
	var err error

	mg := mailgun.NewMailgun(mgDomain, mgApiKey)
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
