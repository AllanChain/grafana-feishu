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
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

var feishuWebhookBase string = "https://open.feishu.cn/open-apis/bot/v2/hook"

func main() {
	feishuWebhook := os.Getenv("FEISHU_WEBHOOK")
	if feishuWebhook == "" {
		log.Println("FEISHU_WEBHOOK not provided, expecting bot uuid in path")
	}
	if os.Getenv("FEISHU_WEBHOOK_BASE") != "" {
		feishuWebhookBase = os.Getenv("FEISHU_WEBHOOK_BASE")
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

	app.Post("/:botUUID?", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		notification := new(Notification)
		if err := c.BodyParser(notification); err != nil {
			return err
		}
		if len(notification.Alerts) == 0 {
			return nil
		}
		for _, alert := range notification.Alerts {
			title, ok := alert.Annotations["summary"]
			if !ok {
				title, ok = alert.Labels["alertname"]
				if !ok {
					title = "[No Title]"
				}
			}
			description, ok := alert.Annotations["description"]
			if !ok {
				description = "[No description]"
			}
			color := "red"
			if alert.Status == "resolved" {
				color = "green"
			}
			feishuCard := &FeishuCard{
				MsgType: "interactive",
				Card: FeishuCardContent{
					Header: FeishuCardHeader{
						Title: FeishuCardTextElement{
							Tag:     "plain_text",
							Content: title,
						},
						Template: color,
					},
					Elements: []FeishuCardDivElement{
						{
							Tag: "div",
							Text: FeishuCardTextElement{
								Tag:     "plain_text",
								Content: description,
							},
						},
					},
				},
			}
			feishuJson, err := json.Marshal(feishuCard)
			if err != nil {
				return err
			}
			feishuWebhook := feishuWebhook
			if botUUID := c.Params("botUUID"); botUUID != "" {
				feishuWebhook = feishuWebhookBase + "/" + botUUID
			}
			request, err := http.NewRequest("POST", feishuWebhook, bytes.NewBuffer(feishuJson))
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			response, err := http.DefaultClient.Do(request)
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
