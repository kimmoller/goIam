IAM Project written in Go and Svelte.

Project is built with docker compose from the root folder.

Core:
    - Handles all the core logic of the application
    - Connects to a postgres db for data storage
    - Supplies other services with apis for fetching data
    - Communicates with handlers through Rabbit MQ

Keycloak-handler:
    - Responsible for managing keycloak accounts
    - Reads Rabbit MQ for new messages

Admin-service:
    - Service between client and core to serve as an layer of authorization

Admin-ui:
    - Admin client for managing identities
