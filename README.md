# FatCat Webhook
Written in Go (~~entirely~~ mostly because I wanted to learn how to properly use Go). Currently supports only Grafana

## How to use
Add the webhook to the Grafana (I tested on Grafana 8.2.6). When configuring an alert, the `tag` tag must be included (as it declares to which queue the message goes to). Each tag you add is added to the message.