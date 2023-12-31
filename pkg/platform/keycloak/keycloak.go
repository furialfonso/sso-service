package keycloak

import (
	"context"
	"cow_sso/pkg/config"

	"github.com/Nerzal/gocloak/v13"
)

type IKeycloakService interface {
	CreateToken(ctx context.Context) string
	GetUserByID(ctx context.Context, token string, userID string) (*gocloak.User, error)
	GetAllUsers(ctx context.Context, token string) ([]*gocloak.User, error)
	GetUserByNickName(ctx context.Context, token string, nickName string) ([]*gocloak.User, error)
	GetRoleByID(ctx context.Context, token string, roleID string) (*gocloak.Role, error)
	CreateUser(ctx context.Context, token string, role *gocloak.Role, user gocloak.User) (string, error)
	DeleteUserByID(ctx context.Context, token string, nickName string) error
}

type keycloakService struct {
	host   *gocloak.GoCloak
	realm  string
	client string
}

func NewKeycloakService() IKeycloakService {
	host := gocloak.NewClient(config.Get().UString("keycloak.host"))
	client := config.Get().UString("keycloak.client")
	return &keycloakService{
		host:   host,
		realm:  config.Get().UString("keycloak.realm"),
		client: client,
	}
}

func (k *keycloakService) CreateToken(ctx context.Context) string {
	realm := config.Get().UString("keycloak.realm-admin")
	user := config.Get().UString("keycloak.user")
	password := config.Get().UString("keycloak.pass")
	token, err := k.host.LoginAdmin(ctx, user, password, realm)
	if err != nil {
		panic("Something wrong with the credentials or url")
	}
	return token.AccessToken
}

func (k *keycloakService) GetUserByID(ctx context.Context, token string, userID string) (*gocloak.User, error) {
	return k.host.GetUserByID(ctx, token, k.realm, userID)
}

func (k *keycloakService) GetAllUsers(ctx context.Context, token string) ([]*gocloak.User, error) {
	return k.host.GetUsers(ctx, token, k.realm, gocloak.GetUsersParams{})
}

func (k *keycloakService) GetUserByNickName(ctx context.Context, token string, nickName string) ([]*gocloak.User, error) {
	users, err := k.host.GetUsers(ctx, token, k.realm, gocloak.GetUsersParams{
		Username: &nickName,
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (k *keycloakService) GetRoleByID(ctx context.Context, token string, roleName string) (*gocloak.Role, error) {
	return k.host.GetRealmRole(ctx, token, k.realm, roleName)
}

func (k *keycloakService) CreateUser(ctx context.Context, token string, role *gocloak.Role, user gocloak.User) (string, error) {
	id, err := k.host.CreateUser(ctx, token, k.realm, user)
	if err != nil {
		return "", err
	}
	var roles []gocloak.Role
	roles = append(roles, *role)
	err = k.host.AddRealmRoleToUser(ctx, token, k.realm, id, roles)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (k *keycloakService) DeleteUserByID(ctx context.Context, token string, userID string) error {
	return k.host.DeleteUser(ctx, token, k.realm, userID)
}
