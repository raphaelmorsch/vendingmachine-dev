package services

import (
	"context"
	"os"

	"log"

	"github.com/Nerzal/gocloak/v8"
)

var (
	keycloakAdminUser     = os.Getenv("KEYCLOAK_ADMIN_USER")
	keycloakAdminPassword = os.Getenv("KEYCLOAK_ADMIN_PASSWORD")
)

func AddRealmAdminUser() {

	token, err := client.LoginAdmin(context.Background(), keycloakAdminUser, keycloakAdminPassword, "master")

	if err != nil {
		log.Fatal("Not possible to log as admin on keycloak")
	}

	cloakClients, err := client.GetClients(context.Background(), token.AccessToken, realm, gocloak.GetClientsParams{
		ClientID: gocloak.StringP("realm-management"),
	})

	if err != nil {
		log.Printf("Not possible Identify Realm Client - Reason %v", err)
	}

	if len(cloakClients) < 1 {
		log.Println("No realm-management client found. Skipping Realm Admin User creation...")
	} else {
		credtls := []gocloak.CredentialRepresentation{{
			Temporary: gocloak.BoolP(false),
			Type:      gocloak.StringP("password"),
			Value:     gocloak.StringP(realmAdminPassword),
		}}

		user := gocloak.User{
			Credentials: &credtls,
			Username:    &realmAdminUser,
			Enabled:     gocloak.BoolP(true),
		}

		user.Credentials = &credtls
		userCreatedId, err := client.CreateUser(context.Background(), token.AccessToken, realm, user)

		if err != nil {
			log.Printf("Not possible Create Realm Admin User - Reason %v", err)
		}

		roleIDs := []string{"create-client", "manage-clients", "manage-realm",
			"manage-users", "query-clients", "query-realms", "query-users", "view-clients", "view-realm", "view-users"}
		var rolesToAdd []gocloak.Role
		for _, roleID := range roleIDs {
			clientRole, err := client.GetClientRole(context.Background(), token.AccessToken, realm, *cloakClients[0].ID, roleID)
			if err != nil {
				log.Fatalf("Not able to find role %v", roleID)
			}
			rolesToAdd = append(rolesToAdd, *clientRole)
		}

		errRoles := client.AddClientRoleToUser(context.Background(), token.AccessToken, realm, *cloakClients[0].ID, userCreatedId, rolesToAdd)

		if errRoles != nil {
			log.Printf("Not possible Add Roles to User - Reason %v", errRoles)
		}
		log.Printf("User created %v\n", userCreatedId)

	}

}
func UpdateAPIClientSercret() {
	token, err := client.LoginAdmin(context.Background(), keycloakAdminUser, keycloakAdminPassword, "master")

	if err != nil {
		log.Fatal("Not possible to log as admin on keycloak")
	}

	cloakClients, err := client.GetClients(context.Background(), token.AccessToken, realm, gocloak.GetClientsParams{
		ClientID: gocloak.StringP(clientId),
	})

	if err != nil {
		log.Printf("Not possible Identify Realm Client - Reason %v", err)
	}

	if len(cloakClients) < 1 {
		log.Println("No realm-management client found. Skipping Realm Admin User creation...")
	} else {
		credential, err := client.RegenerateClientSecret(context.Background(), token.AccessToken, realm, *cloakClients[0].ID)
		if err != nil {
			log.Printf("Not possible Regenerate Secret - Reason %v", err)
		}

		clientSecret = *credential.Value
		log.Printf("Refreshed ClientSecret -> %v", clientSecret)
	}

}
