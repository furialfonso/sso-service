package eureka

import (
	"cow_sso/pkg/config"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/google/uuid"
)

func NewEurekaClient() {

	url := config.Get().UString("eureka.url")
	client := eureka.NewClient([]string{url})

	app := os.Getenv("APPLICATION_NAME")
	host := "localhost"
	appPort := os.Getenv("PORT")
	if os.Getenv("SCOPE") != "local" {
		host = app
		appPort = os.Getenv("SSO_EXTERNAL_PORT")
	}

	fmt.Println("appPort", appPort)
	port, err := strconv.Atoi(appPort)
	if err != nil {
		fmt.Printf("Error convirtiendo el puerto a entero: %v\n", err)
	}

	instance := eureka.NewInstanceInfo(
		host,
		app,
		getIp(),
		port,
		30,
		false)
	instance.InstanceID = fmt.Sprintf("%s:%s", app, uuid.New().String())

	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}

	err = client.RegisterInstance(strings.ToUpper(app), instance)
	if err != nil {
		fmt.Printf("Error registrando la instancia en Eureka: %v\n", err)
	}

	applications, err := client.GetApplications()
	if err != nil {
		fmt.Printf("error buscando aplicaciones: %v\n", err)
	}

	fmt.Println(applications)

	_, err = client.GetInstance(instance.App, instance.InstanceID)
	err = client.SendHeartbeat(instance.App, instance.InstanceID)
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error obteniendo direcciones de red:", err)
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
