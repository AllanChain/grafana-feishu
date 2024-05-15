# Grafana-Feishu

Lightweight server to translate Grafana webhook to Feishu card.

## Usage

The program needs two environment variables:

- `FEISHU_WEBHOOK`: (Optional) The web hook URL to push Feishu notifications. Should be something like `https://open.feishu.cn/open-apis/bot/v2/hook/a21843-123-123-abc987`
- `FEISHU_WEBHOOK_BASE`: (Optional, has no effect if `FEISHU_WEBHOOK` is set) Something like `https://open.feishu.cn/open-apis/bot/v2/hook` (default).
- `FEISHU_WEBHOOK_UUID`: (Optional, has no effect if `FEISHU_WEBHOOK` is set) The bot UUID, something like `a21843-123-123-abc987`.
If the bot UUID is not specified, it should be provided in the path. See below for detail.
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
      - WEBHOOK_AUTH=${WEBHOOK_AUTH}
      - FEISHU_WEBHOOK=${FEISHU_WEBHOOK}
      # Or, instead of FEISHU_WEBHOOK:
      # - FEISHU_WEBHOOK_BASE=${FEISHU_WEBHOOK_BASE}
      # - FEISHU_WEBHOOK_UUID=${FEISHU_WEBHOOK_UUID}
```

The exposed port is `2387`.

After setting up the server, go to Grafana > "Alerting" > "Contact points", add a new contact point with integration as "Webhook". Fill in the URL and credentials. 

If the bot UUID is provided in `FEISHU_WEBHOOK` or `FEISHU_WEBHOOK_UUID` is set, the URL should be `http://grafana-feishu:2387`. Otherwise, the URL should be `http://grafana-feishu:2387/{botUUID}`, where `{botUUID}` is the UUID of the bot you created in Feishu, i.e. the last part of the bot webhook URL.

<img width="737" alt="Grafana config" src="https://user-images.githubusercontent.com/36528777/235901125-181eeb60-df6c-45ff-b550-7756a91c65d1.png">

By default, the color of the card reflects the alert status, and the card title will be the `"summary"` annotation, and the content will be the `"description"` annotation. You can customize the summary and description in alert rules using Go templates. An example:

```
{{ if $values.B }}{{ if eq $values.C.Value 0.0 }}Resolve {{ end }}alert{{ else }}No data{{ end }}
```
