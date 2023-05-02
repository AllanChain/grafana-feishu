package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type Notification struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Status      string            `json:"status"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
	StartsAt    string            `json:"startsAt"`
}

type FeishuCard struct {
	MsgType string            `json:"msg_type"`
	Card    FeishuCardContent `json:"card"`
}

type FeishuCardContent struct {
	Header   FeishuCardHeader       `json:"header"`
	Elements []FeishuCardDivElement `json:"elements"`
}

type FeishuCardHeader struct {
	Title    FeishuCardTextElement `json:"title"`
	Template string                `json:"template"`
}

type FeishuCardTextElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type FeishuCardDivElement struct {
	Tag  string                `json:"tag"`
	Text FeishuCardTextElement `json:"text"`
}

func main() {
	feishuWebhook := os.Getenv("FEISHU_WEBHOOK")
	if feishuWebhook == "" {
		log.Fatal("Please provide FEISHU_WEBHOOK env var")
		return
	}
	app := fiber.New()
	app.Use(logger.New())

	webhookAuth := os.Getenv("WEBHOOK_AUTH")
	if webhookAuth != "" {
		log.Printf("Enabling basic auth")
		parts := strings.SplitN(webhookAuth, ":", 2)
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				parts[0]: parts[1],
			},
		}))
	}

	app.Post("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		notification := new(Notification)
		if err := c.BodyParser(notification); err != nil {
			return err
		}
		if len(notification.Alerts) == 0 {
			return nil
		}
		for _, alert := range notification.Alerts {
			alertname, ok := alert.Labels["alertname"]
			if !ok {
				alertname = "Unnamed Alert"
			}
			summary, ok := alert.Annotations["summary"]
			if !ok {
				summary = "No summary"
			}
			if alertname == "DatasourceNoData" {
				summary = "No data in datasource!"
			}
			color := "red"
			if alert.Status == "resolved" {
				alertname += " is resolved!"
				color = "green"
				summary = "Looks good!"
			}
			feishuCard := &FeishuCard{
				MsgType: "interactive",
				Card: FeishuCardContent{
					Header: FeishuCardHeader{
						Title: FeishuCardTextElement{
							Tag:     "plain_text",
							Content: alertname,
						},
						Template: color,
					},
					Elements: []FeishuCardDivElement{
						{
							Tag: "div",
							Text: FeishuCardTextElement{
								Tag:     "plain_text",
								Content: summary,
							},
						},
					},
				},
			}
			feishuJson, err := json.Marshal(feishuCard)
			if err != nil {
				return err
			}
			request, err := http.NewRequest("POST", feishuWebhook, bytes.NewBuffer(feishuJson))
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				return err
			}
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			log.Printf("Response body: %s", string(body))
		}
		return c.SendStatus(204)
	})

	app.Listen(":2387")
}
