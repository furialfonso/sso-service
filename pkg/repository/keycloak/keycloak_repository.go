package keycloak

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cow_sso/pkg/config"

	"github.com/Nerzal/gocloak/v13"
)

type IKeycloakRepository interface {
	Login(ctx context.Context, user string, password string) (*gocloak.JWT, error)
	Logout(ctx context.Context, refreshToken string) error
	IsValidToken(ctx context.Context, accessToken string) (bool, error)
	GetUserByID(ctx context.Context, token string, userID string) (*gocloak.User, error)
	GetAllUsers(ctx context.Context, token string) ([]*gocloak.User, error)
	GetUserByNickName(ctx context.Context, token string, nickName string) (*gocloak.User, error)
	GetRoleByID(ctx context.Context, token string, roleID string) (*gocloak.Role, error)
	CreateUser(ctx context.Context, token string, role *gocloak.Role, user gocloak.User) (string, error)
	DeleteUserByID(ctx context.Context, token string, userID string) error
}

type keycloakRepository struct {
	host   *gocloak.GoCloak
	realm  string
	client string
	secret string
}

func NewKeycloakRepository() IKeycloakRepository {
	host := gocloak.NewClient(config.Get().UString("keycloak.host"))
	client := config.Get().UString("keycloak.client")
	secret := os.Getenv("KEYCLOAK_SECRET")
	return &keycloakRepository{
		host:   host,
		realm:  config.Get().UString("keycloak.realm"),
		client: client,
		secret: secret,
	}
}

func (k *keycloakRepository) Login(ctx context.Context, user string, password string) (*gocloak.JWT, error) {
	token, err := k.host.Login(ctx, k.client, k.secret, k.realm, user, password)
	if err != nil {
		return nil, errors.New("user or password incorrect")
	}
	return token, nil
}

func (k *keycloakRepository) Logout(ctx context.Context, refreshToken string) error {
	err := k.host.Logout(ctx, k.client, k.secret, k.realm, refreshToken)
	if err != nil {
		return errors.New("Invalid refresh token")
	}
	return nil
}

func (k *keycloakRepository) IsValidToken(ctx context.Context, accessToken string) (bool, error) {
	ret, err := k.host.RetrospectToken(ctx, accessToken, k.client, k.secret, k.realm)
	if err != nil {
		return false, err
	}

	return *ret.Active, nil
}

func (k *keycloakRepository) GetUserByID(ctx context.Context, token string, userID string) (*gocloak.User, error) {
	return k.host.GetUserByID(ctx, token, k.realm, userID)
}

func (k *keycloakRepository) GetAllUsers(ctx context.Context, token string) ([]*gocloak.User, error) {
	return k.host.GetUsers(ctx, token, k.realm, gocloak.GetUsersParams{})
}

func (k *keycloakRepository) GetUserByNickName(ctx context.Context, token string, nickName string) (*gocloak.User, error) {
	users, err := k.host.GetUsers(ctx, token, k.realm, gocloak.GetUsersParams{
		Username: &nickName,
	})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New(fmt.Sprintf("user %s doesn't exist", nickName))
	}

	return users[0], nil
}

func (k *keycloakRepository) GetRoleByID(ctx context.Context, token string, roleName string) (*gocloak.Role, error) {
	return k.host.GetRealmRole(ctx, token, k.realm, roleName)
}

func (k *keycloakRepository) CreateUser(ctx context.Context, token string, role *gocloak.Role, user gocloak.User) (string, error) {
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

func (k *keycloakRepository) DeleteUserByID(ctx context.Context, token string, userID string) error {
	return k.host.DeleteUser(ctx, token, k.realm, userID)
}
