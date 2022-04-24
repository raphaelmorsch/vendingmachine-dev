package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"vendingmachine/config"

	"github.com/Nerzal/gocloak/v8"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/xgfone/cast"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`

	Title string `json:"Title"`

	Description string `json:"Description"`
}

var (
	clientId = os.Getenv("SELLER_CLIENT_ID")

	clientSecret = os.Getenv("SELLER_CLIENT_SECRET")

	realm = os.Getenv("REALM")

	hostname = os.Getenv("CLOAK_HOST")

	realmAdminUser = os.Getenv("REALM_ADMIN_USER")

	realmAdminPassword = os.Getenv("REALM_ADMIN_PASSWORD")
)

var client gocloak.GoCloak

func InitializeOauthServer() {

	client = gocloak.NewClient(hostname)
}

func Protect(next http.Handler, allowedRoles []string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if len(authHeader) < 1 {
			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(config.UnauthorizedError())

			return
		}

		accessToken := strings.Split(authHeader, " ")[1]

		log.Printf("Validating access_token %v \n", accessToken)
		log.Printf("Validating token with clientId=%v / secret=%v / realm=%v \n", clientId, clientSecret, realm)
		rptResult, err := client.RetrospectToken(r.Context(), accessToken, clientId, clientSecret, realm)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(config.BadRequestError(err.Error()))

			return
		}

		if rptResult != nil {
			val, _ := json.MarshalIndent(rptResult, "", " ")
			log.Println(string(val))
		}

		isTokenValid := *rptResult.Active

		if !isTokenValid {

			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(config.InvalidTokenError())

			return

		}

		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return accessToken, nil
		})

		if len(allowedRoles) > 0 {

			if !isRoleAllowed(claims, allowedRoles) {

				w.WriteHeader(http.StatusUnauthorized)

				json.NewEncoder(w).Encode(config.UnauthorizedRoleError())

				return

			}
		}

		userId, _ := cast.ToString(claims["preferred_username"])

		adminToken, _ := client.LoginAdmin(r.Context(), keycloakAdminUser, keycloakAdminPassword, "master")
		users, _ := client.GetUsers(r.Context(), adminToken.AccessToken, realm, gocloak.GetUsersParams{
			Username: gocloak.StringP(userId),
		})
		log.Println("Begining user session validation...user=" + userId)
		if len(users) < 1 {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(config.DataAccessLayerError("Not possible to retrieve User/Session Information"))

			return

		}
		sessions, sessionErr := client.GetUserSessions(r.Context(), adminToken.AccessToken, realm, *users[0].ID)
		if sessionErr != nil {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(config.DataAccessLayerError("Not possible to retrieve Session Information"))

			return

		}

		if len(sessions) > 1 {
			w.WriteHeader(http.StatusConflict)

			json.NewEncoder(w).Encode(config.SessionConflictError("There is an Session already opened for this user. All Sessions for this user will be erased"))

			for _, session := range sessions {
				client.LogoutUserSession(r.Context(), adminToken.AccessToken, realm, *session.ID)

			}

			return

		}

		r.Header.Add("user_id", userId)
		next.ServeHTTP(w, r)
	})
}

func isRoleAllowed(claims jwt.MapClaims, allowedRoles []string) bool {
	for key, val := range claims {
		if key == "resource_access" {
			raMap, _ := cast.ToStringMap(val)
			for _, vv := range raMap {
				rolesMap, _ := cast.ToStringMap(vv)
				for _, vvv := range rolesMap {
					roleSlice, _ := cast.ToSlice(vvv)
					for _, role := range roleSlice {

						for _, allowedRole := range allowedRoles {
							if allowedRole == role {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

// swagger:route POST /user user addNewUser
// Adds new user
//
// security:
// - apiKey: []
// responses:
//  400: BadRequest
//  401: CommonError
//  201: UserCreated
func AddNewUser(w http.ResponseWriter, r *http.Request) {

	token, err := client.LoginAdmin(r.Context(), realmAdminUser, realmAdminPassword, realm)

	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(
			err.Error())
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err.Error())
	}

	var user gocloak.User

	json.Unmarshal(reqBody, &user)

	credtls := []gocloak.CredentialRepresentation{{
		Temporary: gocloak.BoolP(true),
		Type:      gocloak.StringP("password"),
		Value:     gocloak.StringP("12345"),
	}}
	user.Credentials = &credtls
	userCreatedId, err := client.CreateUser(r.Context(), token.AccessToken, realm, user)

	user.ID = &userCreatedId

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			err.Error())
		return
	}

	cloakClients, err := client.GetClients(r.Context(), token.AccessToken, realm, gocloak.GetClientsParams{
		ClientID: gocloak.StringP(clientId),
	})

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			err.Error())
		return
	}

	log.Printf("Looking for Clients on realm=%v / clientID=%v --> Clients found %v", realm, clientId, len(cloakClients))

	roleID := "buyer"
	clientRole, err := client.GetClientRole(r.Context(), token.AccessToken, realm, *cloakClients[0].ID, roleID)

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			err.Error())
		return
	}
	var rolesToAdd []gocloak.Role
	rolesToAdd = append(rolesToAdd, *clientRole)
	err = client.AddClientRoleToUser(r.Context(), token.AccessToken, realm, *cloakClients[0].ID, userCreatedId, rolesToAdd)

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			err.Error())
		return
	}

	json.NewEncoder(w).Encode(user)

}
