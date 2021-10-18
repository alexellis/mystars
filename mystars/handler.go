package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/hmac"

	github "github.com/google/go-github/v39/github"
)

type SlackMsg struct {
	Text string `json:"text"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	webhookVal, err := ioutil.ReadFile("/var/openfaas/secrets/slack-stars-webhook-url")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	hmacVal, err := ioutil.ReadFile("/var/openfaas/secrets/slack-stars-hmac-secret")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	hmac := strings.TrimSpace(string(hmacVal))
	webhookURL := strings.TrimSpace(string(webhookVal))

	event, err := parseWebhook(r, hmac)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	emoji := "star"
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "watched" {
		eventType = "stared"
	}
	if r.Header.Get("X-GitHub-Event") == "fork" {
		emoji = "open_file_folder"
	}

	repo := event.GetRepo()
	slackMsg := SlackMsg{
		Text: fmt.Sprintf(":%s: %s/%s by %s",
			emoji, *repo.GetOwner().Login, repo.GetName(),
			*event.GetSender().Login),
	}

	if event.GetSender().HTMLURL != nil {
		slackMsg.Text += fmt.Sprintf(" - %s", *event.GetSender().HTMLURL)
	}

	slackMsgBytes, _ := json.Marshal(slackMsg)
	slackMsgReader := bytes.NewReader(slackMsgBytes)
	req, err := http.NewRequest(http.MethodPost, webhookURL, slackMsgReader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error creating request: %s", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		msg := fmt.Sprintf("Status code unexpected from Slack: %d, body: %s", res.StatusCode, string(body))
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", msg)
		return
	}

	doneMsg := fmt.Sprintf("Sent message %s/%s by %s (%s)", *repo.GetOwner().Login, repo.GetName(), *event.GetSender().Login, eventType)

	fmt.Fprintln(os.Stdout, doneMsg)
	fmt.Fprintln(w, doneMsg)
}

// parseWebhook parses the webhook request from the body as JSON
// and verifies the HMAC signature is valid with the local stored secret.
func parseWebhook(r *http.Request, hmacSecret string) (*github.StarEvent, error) {
	event := &github.StarEvent{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return event, err
	}

	valid := hmac.Validate(body, r.Header.Get("X-Hub-Signature"), hmacSecret)
	if valid != nil {
		return event, fmt.Errorf("invalid signature: %w", valid)
	}

	err = json.Unmarshal(body, &event)
	return event, err
}
