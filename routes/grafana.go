package routes

import (
	"context"
	"encoding/json"
	"fatcat_webhook/m/v2/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

type dictionary = map[string]any
type GrafanaAlert struct {
	Title       string     `json:"title"`
	RuleId      int        `json:"ruleId"`
	RuleName    string     `json:"ruleName"`
	State       string     `json:"state"`
	EvalMatches []any      `json:"evalMatches"`
	OrgId       int        `json:"orgId"`
	DashboardId int        `json:"dashboardId"`
	PanelId     int        `json:"panelId"`
	Tags        dictionary `json:"tags"`
	RuleUrl     string     `json:"ruleUrl"`
	Message     string     `json:"message"`
}

func GrafanaHandler(c *gin.Context) {
	var body GrafanaAlert
	if err := c.BindJSON(&body); err != nil {
		log.Println("[-] Got a bad body")
		c.IndentedJSON(422, dictionary{
			"message": "bad body",
		})
		return
	}

	if body.State != "alerting" || body.RuleName == "Test notification" {
		log.Println("[*] Got either a test message or an alert that is not of state `alerting`")
		c.IndentedJSON(400, dictionary{
			"message": "Either a test message or it's not alerting",
		})
		return
	}

	message := dictionary{
		"message": body.Message,
	}
	for k, v := range body.Tags {
		if k == "tag" {
			message["FATCAT_QUEUE_TAG"] = v
			continue
		}
		message[k] = v
	}

	actualMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("[-] Couldn't create the message object")
		c.IndentedJSON(500, dictionary{
			"message": "Failed to create a message object",
		})
		return
	}

	connection, err := amqp091.Dial(
		"amqp://" + utils.Getenv("FATCAT_RABBITMQ_USERNAME", "guest") + ":" +
			utils.Getenv("FATCAT_RABBITMQ_PASSWORD", "guest") + "@" +
			utils.Getenv("FATCAT_RABBITMQ_HOST", "localhost") + ":" + utils.Getenv("FATCAT_RABBITMQ_PORT", "5672"),
	)

	if err != nil {
		log.Println("[-] Couldn't connect to the RabbitMQ")
		c.IndentedJSON(500, dictionary{
			"message": "Failed to connect to the RabbitMQ",
		})
		return
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Println("[-] Couldn't get the RabbitMQ channel")
		defer connection.Close()
		c.IndentedJSON(500, dictionary{
			"message": "Failed to get the RabbitMQ channel",
		})
		return
	}

	_, err = channel.QueueDeclare(utils.Getenv("FATCAT_RABBITMQ_QUEUE", "incoming"), false, false, false, false, nil)
	if err != nil {
		log.Println("[-] Couldn't declare the queue " + utils.Getenv("FATCAT_RABBITMQ_QUEUE", "incoming"))
		defer connection.Close()
		c.IndentedJSON(500, dictionary{
			"message": "Filaed to declare the queue",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx, "", utils.Getenv("FATCAT_RABBITMQ_QUEUE", "incoming"), false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        actualMessage,
	})
	if err != nil {
		log.Println("[-] Failed to send message to queue " + utils.Getenv("FATCAT_RABBITMQ_QUEUE", "incoming"))
		defer connection.Close()
		c.IndentedJSON(500, dictionary{
			"message": "Failed to send the message to the queue",
		})
		return
	}

	defer connection.Close()

	log.Println("[*] Message successfully sent")

	c.IndentedJSON(202, "")
}
