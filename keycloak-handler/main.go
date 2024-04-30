package main

import (
	"log"

	"github.com/Nerzal/gocloak/v13"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/guregu/null.v3"
)

type AccountProvision struct {
	IdentityID string      `json:"identityId"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Email      string      `json:"email"`
	AccountID  null.String `json:"accountId"`
	Username   null.String `json:"username"`
	SystemId   null.String `json:"systemId"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func consumeQueueu(channel *amqp.Channel, queueName string) <-chan amqp.Delivery {
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue %s, %s", queueName, err)
	}
	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register consumer for queue %s, %s", queueName, err)
	}
	return messages
}

func main() {
	client := gocloak.NewClient("http://0.0.0.0:8080")

	conn, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	createMessages := consumeQueueu(channel, "accountCreate_keycloak_private")
	enableMessages := consumeQueueu(channel, "accountEnable_keycloak_private")
	disableMessages := consumeQueueu(channel, "accountDisable_keycloak_private")
	deleteMessages := consumeQueueu(channel, "accountDelete_keycloak_private")

	var forever chan struct{}

	go func() {
		for message := range createMessages {
			createAccount(client, channel, &message)
		}
	}()

	go func() {
		for message := range enableMessages {
			enableAccount(client, channel, &message)
		}
	}()

	go func() {
		for message := range disableMessages {
			disableAccount(client, channel, &message)
		}
	}()

	go func() {
		for message := range deleteMessages {
			deleteAccount(client, channel, &message)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
