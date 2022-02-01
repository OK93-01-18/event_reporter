package mattermost

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Mattermost struct {
	username   string
	webhookUrl string
}

func New(username string, hookUrl string) *Mattermost {
	return &Mattermost{
		username:   username,
		webhookUrl: hookUrl,
	}
}

func (m *Mattermost) Send(ctx context.Context, subject string, msg string) error {

	bodymsg := "{\"username\" : \"" + m.username + "\", \"text\" : \"**Subject:** " + subject + "\n" +
		"**Message:**\n ```log" + msg + "```\"}"

	req, err := http.NewRequest("POST", m.webhookUrl, bytes.NewBuffer([]byte(bodymsg)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s | %v\n", resp.Status, string(b))
	}

	return nil

}
