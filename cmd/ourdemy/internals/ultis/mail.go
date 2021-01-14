package ultis

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func SendMail(toUser string, email string, subject string, content string) error {
	user := viper.Get("SGUSER")
	key := viper.Get("SGKEY")
	from := mail.NewEmail("Support", user.(string))
	to := mail.NewEmail(toUser, email)
	message := mail.NewSingleEmail(from, subject, to, content, "")
	client := sendgrid.NewSendClient(key.(string))
	_, err := client.Send(message)

	return err
}
