package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Nerzal/gocloak/v13"
	amqp "github.com/rabbitmq/amqp091-go"
)

var realm = "private"

func createAccount(client *gocloak.GoCloak, channel *amqp.Channel, message *amqp.Delivery) {
	log.Printf("Received create message: %s", message.Body)
	ctx := context.Background()

	var account AccountProvision
	err := json.Unmarshal(message.Body, &account)

	if err != nil {
		log.Println(err)
	} else {
		requiredActions := []string{"UPDATE_PASSWORD", "VERIFY_EMAIL"}
		user := gocloak.User{
			Enabled:         gocloak.BoolP(false),
			FirstName:       gocloak.StringP(account.FirstName),
			LastName:        gocloak.StringP(account.LastName),
			Email:           gocloak.StringP(account.Email),
			Username:        gocloak.StringP(account.Username.String),
			RequiredActions: &requiredActions,
		}

		token, err := client.LoginAdmin(ctx, "keycloak-handler", "keycloak", realm)
		if err != nil {
			panic("Something wrong with the credentials or url")
		}
		_, err = client.CreateUser(ctx, token.AccessToken, realm, user)

		if err != nil {
			log.Printf("Error while creating user, %s", err)
		} else {
			err = channel.PublishWithContext(context.Background(),
				"",                    // exchange
				"accountCreateCommit", // routing key
				false,                 // mandatory
				false,                 // immediate
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

func enableAccount(client *gocloak.GoCloak, channel *amqp.Channel, message *amqp.Delivery) {
	log.Printf("Received enable message: %s", message.Body)
	ctx := context.Background()

	var account AccountProvision
	err := json.Unmarshal(message.Body, &account)

	if err != nil {
		log.Println(err)
	} else {
		token, err := client.LoginAdmin(ctx, "keycloak-handler", "keycloak", realm)
		if err != nil {
			panic("Something wrong with the credentials or url")
		}
		users, err := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{Username: &account.Username.String})
		if err != nil {
			log.Printf("Error while getting user %s, %s", account.Username.String, err)
		} else {
			user := users[0]
			user.Enabled = gocloak.BoolP(true)
			err = client.UpdateUser(ctx, token.AccessToken, realm, *user)
			if err != nil {
				log.Printf("Error while enabling user %s, %s", account.Username.String, err)
			} else {
				log.Printf("Commit user %s", account.IdentityID)
				err = channel.PublishWithContext(context.Background(),
					"",                    // exchange
					"accountEnableCommit", // routing key
					false,                 // mandatory
					false,                 // immediate
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

}

func disableAccount(client *gocloak.GoCloak, channel *amqp.Channel, message *amqp.Delivery) {
	log.Printf("Received disable message: %s", message.Body)
	ctx := context.Background()

	var account AccountProvision
	err := json.Unmarshal(message.Body, &account)

	if err != nil {
		log.Println(err)
	} else {
		token, err := client.LoginAdmin(ctx, "keycloak-handler", "keycloak", realm)
		if err != nil {
			panic("Something wrong with the credentials or url")
		}
		users, err := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{Username: &account.Username.String})
		if err != nil {
			log.Printf("Error while getting user %s, %s", account.Username.String, err)
		} else {
			user := users[0]
			user.Enabled = gocloak.BoolP(false)
			err = client.UpdateUser(ctx, token.AccessToken, realm, *user)
			if err != nil {
				log.Printf("Error while disabling user %s, %s", account.Username.String, err)
			} else {
				log.Printf("Commit user %s", account.IdentityID)
				err = channel.PublishWithContext(context.Background(),
					"",                     // exchange
					"accountDisableCommit", // routing key
					false,                  // mandatory
					false,                  // immediate
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
}

func deleteAccount(client *gocloak.GoCloak, channel *amqp.Channel, message *amqp.Delivery) {
	log.Printf("Received delete message: %s", message.Body)
	ctx := context.Background()

	var account AccountProvision
	err := json.Unmarshal(message.Body, &account)

	if err != nil {
		log.Println(err)
	} else {
		token, err := client.LoginAdmin(ctx, "keycloak-handler", "keycloak", realm)
		if err != nil {
			panic("Something wrong with the credentials or url")
		}
		users, err := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{Username: &account.Username.String})
		if err != nil {
			log.Printf("Error while getting user %s, %s", account.Username.String, err)
		} else {
			user := users[0]
			err = client.DeleteUser(ctx, token.AccessToken, realm, *user.ID)
			if err != nil {
				log.Printf("Error while deleting user %s, %s", account.Username.String, err)
			} else {
				log.Printf("Commit user %s", account.IdentityID)
				err = channel.PublishWithContext(context.Background(),
					"",                    // exchange
					"accountDeleteCommit", // routing key
					false,                 // mandatory
					false,                 // immediate
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
}
