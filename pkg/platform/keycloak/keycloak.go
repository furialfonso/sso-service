package keycloak

import (
	"cow_sso/pkg/config"

	"github.com/Nerzal/gocloak"
)

type IKeycloakService interface {
	CreateUser(user gocloak.User) (*string, error)
}

type keycloakService struct {
	client gocloak.GoCloak
	token  *gocloak.JWT
	realm  string
}

func NewKeycloakService() IKeycloakService {
	client := gocloak.NewClient(config.Get().UString("keycloak.host"))
	return &keycloakService{
		client: client,
		token:  conn(client),
		realm:  config.Get().UString("keycloak.realm"),
	}
}

func conn(client gocloak.GoCloak) *gocloak.JWT {
	realm := config.Get().UString("keycloak.realm-admin")
	user := config.Get().UString("keycloak.user")
	password := config.Get().UString("keycloak.pass")
	token, err := client.LoginAdmin(user, password, realm)
	if err != nil {
		panic(err)
	}
	return token
}

func (k *keycloakService) CreateUser(user gocloak.User) (*string, error) {
	var id *string
	id, err := k.client.CreateUser(k.token.AccessToken, k.realm, user)
	if err != nil {
		return id, err
	}
	return id, nil
}
