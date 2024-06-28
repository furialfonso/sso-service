package dependencies

import (
	"cow_sso/api/handlers"
	"cow_sso/api/server"
	"cow_sso/middleware"
	"cow_sso/pkg/auth"
	"cow_sso/pkg/platform/eureka"
	"cow_sso/pkg/platform/keycloak"
	"cow_sso/pkg/platform/restful"
	"cow_sso/pkg/services"

	"go.uber.org/dig"
)

type Dependencies struct{}

func BuildDependencies() *dig.Container {
	Container := dig.New()
	_ = Container.Invoke(eureka.NewEurekaClient)
	_ = Container.Provide(middleware.NewCorsConfig)
	_ = Container.Provide(server.New)
	_ = Container.Provide(server.NewRouter)
	_ = Container.Provide(handlers.NewHandlerPing)
	_ = Container.Provide(handlers.NewUserHandler)
	_ = Container.Provide(services.NewUserService)
	_ = Container.Provide(keycloak.NewKeycloakService)
	_ = Container.Provide(restful.NewRestfulService)
	_ = Container.Provide(handlers.NewAuthHandler)
	_ = Container.Provide(auth.NewAuthService)
	return Container
}
