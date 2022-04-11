# vendingmachine-dev

(Work still in-progress)

**Description**: This API is developed in Go and provides endpoints to simulate basic Backend behavior of a Vending Machine.

- **External Dependencies**: 
    - User's Authentication / Authorization handled by keycloak (gocloak library is used to make API communication with keycloak instance)
    - MySQL as a database to persist necessary information properly
- **Containerazation**
    - Temporarily using docker.io/library/mysql:latest image to run MySQL instance (listening on port 3306)
    - Temporarily using docker.io/jboss/keycloak:latest image to run Keycloak instance (listening on port 8086)
- **Environment customization**
    - So far, due to lack of time, it wasn't possible to provide a docker-compose to automate all the environment preparation automatically, so a few steps are necessary to make it running:
    - Create manually env variables: 
        - ```$ export SELLER_CLIENT_ID=vendingmachine-app```
        - Login to Keycloak as admin:admin and create a realm called **vendingmachine**
        - Create on Keycloak an Client called ```vendingmachine-app```
            - Access Type **confidential**
            - Standard Flow **enabled**
            - Direct Access Granted **enabled**
            - Service Account **enabled**
            - Authorization **enabled**
            - Valid Redirect URIs: *
            - Go to the Tab **Credentials** and copy the **Secret** value
            - Create the env variable **SELLER_CLIENT_SECRET** ```export SELLER_CLIENT_SECRET=<<the-secret-you-copied-on-previous-step>>
            - Create the env variable **SELLER_CLIENT_ID** ```export SELLER_CLIENT_ID=vendingmachine-app```
        - On Keycloak, for this client, create 2 new Roles: **buyer** and **seller**
        - On Keycloak, create an User called **vending-machine-admin**
