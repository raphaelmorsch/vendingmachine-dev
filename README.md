# vendingmachine-dev

(Work still in-progress)

**Description**: This API is developed in Go and provides endpoints to simulate basic Backend behavior of a Vending Machine.

- **External Dependencies**: 
    - User's Authentication / Authorization handled by keycloak (gocloak library is used to make API communication with keycloak instance)
    - MySQL as a database to persist necessary information properly
- **Containerazation**
    - Using docker.io/library/mysql:latest image to run MySQL instance (listening on port 3306)
    - Using docker.io/jboss/keycloak:latest image to run Keycloak instance (listening on port 8086)
    - Using podman to create the pod with all dependant containers
    - vending-machine:latest container is available at quay.io repository (https://quay.io/repository/ramoreir/vending-machine)
    - Check Dockerfile, pods.sh and restart.sh for more information
- **Environment customization**
    - Create a new User to act as a seller. Name it as you want. On the **Role Mappings** tab, pick **vendingmachine-app** as **Client Role** and assign the **seller** Role to it
    - Create a new User to act as a buyer. Name it as you want. On the **Role Mappings** tab, pick **vendingmachine-app** as **Client Role** and assign the **buyer** Role to it
    - This steps can be done on keycloak UI or with the endpoint POST /user
- **Missing Features** (Unfortunately, due to a lack of time it was not possible to finish)
    - OpenAPI docs
    - Deep Refactoring
