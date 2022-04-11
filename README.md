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
        - On Keycloak, for this client, create 2 new Roles: **buyer** and **seller**
        - On Keycloak, create an User called **vending-machine-admin**
            - On the **Credentials** tab add the password **12345** and change the **Temporary** flag to **off**
            - On the **Role Mappings** tab, pick the **realm-management** as **Client Roles** and add the following Roles:
                **create-client**
                **manage-clients**
                **manage-realm**
                **manage-users**
            - Create a new User to act as a seller. Name it as you want. On the **Role Mappings** tab, pick **vendingmachine-app** as **Client Role** and assign the **seller** Role to it
            - Create a new User to act as a buyer. Name it as you want. On the **Role Mappings** tab, pick **vendingmachine-app** as **Client Role** and assign the **buyer** Role to it
        - Create the env variable REALM with the value vendingmachine ```export REALM=vendingmachine```
        - Create the env variable CLOAK_HOST with the value http://localhost:8086 ```export CLOAK_HOST=http://localhost:8086```
        - Create the env variable REALM_ADMIN_USER with the value vending-machine-admin ```export REALM_ADMIN_USER=vending-machine-admin```
        - Create the env variable REALM_ADMIN_PASSWORD with the value "12345" ```export REALM_ADMIN_PASSWORD=12345```
    - This previous steps should all be removed and the process would be all automated, once the Dockerfile it's done for the Docker Compose
- **Missing Features**
    - OpenAPI docs
    - DockerCompose
    - Not all API features are covered by Tests -> There is an ```purchaseService_test.go``` as an example of how Unit Tests will behave
