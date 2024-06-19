package eureka

import (
	"cow_sso/pkg/config"
	"errors"
	"fmt"
	"os"
	"strconv"

	eureka "github.com/xuanbo/eureka-client"
)

type IEurekaClient interface {
}

type eurekaClient struct {
}

func NewEurekaClient() IEurekaClient {
	err := register()
	if err != nil {
		panic(fmt.Errorf("[eureka] [NewEurekaService] error registering with eureka: %v", err))
	}
	return &eurekaClient{}
}

func register() error {
	portStr := os.Getenv("API_PORT")
	if portStr == "" {
		return errors.New("port is empty")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	name := fmt.Sprintf(os.Getenv("APPLICATION_NAME"))
	if name == "" {
		return errors.New("application name is empty")
	}

	url := config.Get().UString("eureka.url")
	if url == "" {
		return errors.New("eureka url is empty")
	}

	eurekaConfig := eureka.Config{
		App:                   name,
		Port:                  port,
		DefaultZone:           url,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
	}

	client := eureka.NewClient(&eurekaConfig)
	client.Start()

	return nil
}
