package mailer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/davesaah/fatch/internal/config"
)

type Mailer struct {
	config config.MailConfig
}

func New(cfg config.MailConfig) *Mailer {
	return &Mailer{config: cfg}
}

func (m *Mailer) Send(to, subject, body string) error {
	payload := strings.NewReader(
		fmt.Sprintf(
			"{\"from\":{\"email\":\"%s\",\"name\":\"Fatch\"},\"to\":[{\"email\":\"%s\"}],\"subject\":\"%s\",\"text\":\"%s\",\"category\":\"OTP\"}",
			m.config.From, to, subject, body,
		),
	)

	req, err := http.NewRequest("POST", m.config.Host, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+m.config.Key)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
