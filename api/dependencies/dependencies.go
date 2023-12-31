package dependencies

import (
	"cow_sso/api/handlers"
	"cow_sso/api/server"
	"cow_sso/pkg/platform/keycloak"
	"cow_sso/pkg/services"

	"go.uber.org/dig"
)

type Dependencies struct{}

func BuildDependencies() *dig.Container {
	Container := dig.New()
	_ = Container.Provide(server.New)
	_ = Container.Provide(server.NewRouter)
	_ = Container.Provide(handlers.NewHandlerPing)
	_ = Container.Provide(handlers.NewUserHandler)
	_ = Container.Provide(services.NewUserService)
	_ = Container.Provide(keycloak.NewKeycloakService)

	return Container
}
