# SSO-SERVICE

It is a 3-tier based architecture with dependency injection.

## Author
  - *Andres Felipe Alfonso Ortiz*

### Technologies
  - *Golang*: programming language.
  - *Gin*: framework for rest applications.
  - *Mokery*: automatic mocks for unit tests.
  - *Dig*: automatic dependency injection.
  - *Docker*: application's contenerization.

### Architecture
```mermaid
graph TB
    User((External User))

    subgraph "SSO Service System"
        subgraph "API Gateway Container"
            GinServer["API Server<br>(Gin)"]
            Router["URL Router<br>(Gin Router)"]
            MetricsMiddleware["Metrics Middleware<br>(Prometheus)"]
        end

        subgraph "Handler Container"
            AuthHandler["Auth Handler<br>(Go)"]
            UserHandler["User Handler<br>(Go)"]
            PingHandler["Ping Handler<br>(Go)"]
        end

        subgraph "Service Container"
            AuthService["Auth Service<br>(Go)"]
            UserService["User Service<br>(Go)"]
        end

        subgraph "Client Container"
            KeycloakClient["Keycloak Client<br>(GoCloak)"]
            RestClient["REST Client<br>(Go HTTP)"]
            TeamClient["Team Client<br>(Go)"]
        end

        subgraph "External Systems"
            Keycloak["Keycloak Server<br>(Keycloak)"]
            TeamAPI["Team API<br>(REST)"]
            PrometheusMetrics["Prometheus<br>(Metrics)"]
        end
    end

    %% User interactions
    User -->|"HTTP Requests"| GinServer

    %% API Gateway relationships
    GinServer -->|"Routes requests"| Router
    Router -->|"Applies"| MetricsMiddleware
    MetricsMiddleware -->|"Exposes metrics"| PrometheusMetrics

    %% Router to Handler relationships
    Router -->|"/auth/*"| AuthHandler
    Router -->|"/users/*"| UserHandler
    Router -->|"/ping"| PingHandler

    %% Handler to Service relationships
    AuthHandler -->|"Uses"| AuthService
    UserHandler -->|"Uses"| UserService

    %% Service to Client relationships
    AuthService -->|"Authenticates via"| KeycloakClient
    UserService -->|"Manages users via"| KeycloakClient
    UserService -->|"Gets team info via"| TeamClient
    TeamClient -->|"Makes HTTP calls via"| RestClient

    %% External system connections
    KeycloakClient -->|"Authenticates/Manages users"| Keycloak
    TeamClient -->|"Fetches team data"| TeamAPI

    %% Styling
    classDef container fill:#e6e6e6,stroke:#666,stroke-width:2px
    classDef component fill:#fff,stroke:#000,stroke-width:1px
    classDef external fill:#ccf,stroke:#66f,stroke-width:2px
    
    class GinServer,Handler,Service,Client container
    class Router,AuthHandler,UserHandler,PingHandler,AuthService,UserService,KeycloakClient,RestClient,TeamClient component
    class Keycloak,TeamAPI,PrometheusMetrics external
```

### Run unit tests
  ```
    export CONFIG_DIR=$(pwd)/pkg/config && export SCOPE=local && go test -v ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./... -count=1
  ```
  #### Look result in html
  ```
    go tool cover -html=coverage.out
  ```

**Gin**
  - Documentation
    - https://gin-gonic.com/docs/quickstart/

**Mokery**
  - Documentacion
    - https://vektra.github.io/mockery/latest/
  - Instalacion 
    - mac
    ```
      brew install mockery
    ```
    - windows
    ```
    docker pull vektra/mockery
    ```
  - Crear mocks
    - Mac:
    ```
      mockery --all --disable-version-string
    ```
    - Windows:
    ```
      docker run -v $PWD:/src -w /src vektra/mockery --all
    ```
  - Sort app
    ```
      fieldalignment -fix ./...
    ```
    ```
      gofumpt -l -w .
    ```
  
**Dig**
  - Documentation
    - https://ruslan.rocks/posts/golang-dig
    - https://www.golanglearn.com/golang-tutorials/golang-dig-a-better-way-to-manage-dependency/

**Start Aplication**
  Execute the next command for start/stop the application.
  - API
   ```
    docker-compose up/down 
    docker push furialfonso/sso-service:test
  ```
**Config project**
  - For unit test
  ```
    "go.testEnvVars": {
          "CONFIG_DIR": "${workspaceRoot}/pkg/config",
          "SCOPE":"local"
      },
  ```
  - Environment vs-code
  ```
    "APPLICATION_NAME": "sso-service-local",
    "SCOPE": "local",
    "PORT": "8080",
    "CONFIG_DIR": "${workspaceRoot}/pkg/config",
    "GIN_MODE": "release",
    "KEYCLOAK_SECRET": ""
  ```

**Utils**
- docker-golang
```
https://www.youtube.com/watch?v=Ms5RKs8TNU4&t=1504s
```