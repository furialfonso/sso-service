package eureka

import (
	"cow_sso/pkg/config"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	eureka "github.com/xuanbo/eureka-client"
)

func NewEurekaClient() {
	app := os.Getenv("APPLICATION_NAME")
	host := "localhost"
	appPort := os.Getenv("PORT")

	if os.Getenv("SCOPE") != "local" {
		host = app
		appPort = os.Getenv("SSO_INTERNAL_PORT")
	}
	port, err := strconv.Atoi(appPort)
	if err != nil {
		fmt.Printf("Error convirtiendo el puerto a entero: %v\n", err)
	}
	client := eureka.NewClient(&eureka.Config{
		DefaultZone:           fmt.Sprintf("%s/", config.Get().UString("eureka.url")),
		App:                   app,
		Port:                  port,
		InstanceID:            fmt.Sprintf("%s:%s", app, uuid.New().String()),
		HostName:              host,
		RenewalIntervalInSecs: config.Get().UInt("eureka.interval"),
		DurationInSecs:        config.Get().UInt("eureka.duration"),
	})

	client.Start()
}
