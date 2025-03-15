package dependencies

import (
	"cow_sso/api/handlers"
	authHandler "cow_sso/api/handlers/auth"
	userHandler "cow_sso/api/handlers/user"
	"cow_sso/api/server"
	"cow_sso/middleware"
	"cow_sso/pkg/integration/keycloak"
	"cow_sso/pkg/integration/restful"
	"cow_sso/pkg/integration/team"
	authService "cow_sso/pkg/service/auth"
	userService "cow_sso/pkg/service/user"

	"go.uber.org/dig"
)

type Dependencies struct{}

func BuildDependencies() *dig.Container {
	Container := dig.New()
	_ = Container.Provide(middleware.NewMetricMiddleWare)
	_ = Container.Provide(server.New)
	_ = Container.Provide(server.NewRouter)
	//handlers
	_ = Container.Provide(handlers.NewHandlerPing)
	_ = Container.Provide(userHandler.NewUserHandler)
	_ = Container.Provide(authHandler.NewAuthHandler)
	//services
	_ = Container.Provide(userService.NewUserService)
	_ = Container.Provide(authService.NewAuthService)
	//repositories
	_ = Container.Provide(keycloak.NewKeycloakClient)
	_ = Container.Provide(team.NewTeamClient)
	//platform
	_ = Container.Provide(restful.NewRestClient)
	return Container
}
