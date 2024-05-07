# Grafana-Feishu

Lightweight server to translate Grafana webhook to Feishu card.

## Usage

The program needs two environment variables:

- `FEISHU_WEBHOOK_BASE`: (Optional) The web hook base URL to push Feishu notifications. Default is `https://open.feishu.cn/open-apis/bot/v2/hook`
- `WEBHOOK_AUTH`: (Optional) The username and password. Should be something like `user:password`

Here is an example docker compose file:

```yaml
version: 3
services:
  grafana:
    # ...
  grafana-feishu:
    image: allanchain/grafana-feishu
    container_name: grafana-feishu
    restart: always
    environment:
      - FEISHU_WEBHOOK_BASE=${FEISHU_WEBHOOK_BASE}
      - WEBHOOK_AUTH=${WEBHOOK_AUTH}
```

The exposed port is `2387`.

After setting up the server, go to Grafana > "Alerting" > "Contact points", add a new contact point with integration as "Webhook". Fill in the URL and credentials. 

The URL should be like `http://grafana-feishu:2387/{botUUID}`, where `{botUUID}` is the UUID of the bot you created in Feishu, i.e. the last part of the bot webhook URL.

In the previous version, a whole bot URL (including the UUID) can be set with the environment variable `FEISHU_WEBHOOK`, and the bot UUID is not needed in Grafana. This is still supported.

<img width="737" alt="Grafana config" src="https://user-images.githubusercontent.com/36528777/235901125-181eeb60-df6c-45ff-b550-7756a91c65d1.png">

By default, the color of the card reflects the alert status, and the card title will be the `"summary"` annotation, and the content will be the `"description"` annotation. You can customize the summary and description in alert rules using Go templates. An example:

```
{{ if $values.B }}{{ if eq $values.C.Value 0.0 }}Resolve {{ end }}alert{{ else }}No data{{ end }}
```
