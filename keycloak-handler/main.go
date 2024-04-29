package main

import (
	"context"
	"encoding/json"
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

func main() {
	client := gocloak.NewClient("http://0.0.0.0:8080")
	ctx := context.Background()

	conn, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"accountProvision_keycloak_private",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)

			var account AccountProvision
			err := json.Unmarshal(message.Body, &account)

			if err != nil {
				log.Println(err)
			} else {
				requiredActions := []string{"UPDATE_PASSWORD", "VERIFY_EMAIL"}
				user := gocloak.User{
					Enabled:         gocloak.BoolP(true),
					FirstName:       gocloak.StringP(account.FirstName),
					LastName:        gocloak.StringP(account.LastName),
					Email:           gocloak.StringP(account.Email),
					Username:        gocloak.StringP(account.Username.String),
					RequiredActions: &requiredActions,
				}

				token, err := client.LoginAdmin(ctx, "keycloak-handler", "keycloak", "private")
				if err != nil {
					panic("Something wrong with the credentials or url")
				}
				_, err = client.CreateUser(ctx, token.AccessToken, "private", user)

				if err != nil {
					log.Printf("Error while creating user, %s", err)
				} else {
					err = channel.PublishWithContext(context.Background(),
						"",              // exchange
						"accountCommit", // routing key
						false,           // mandatory
						false,           // immediate
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(message.Body),
						})
					if err != nil {
						log.Printf("Error while sending message: %s", err)
					}
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
