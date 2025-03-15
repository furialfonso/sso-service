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
    
    subgraph "SSO Service"
        subgraph "API Layer"
            GinServer["API Server<br>Gin Framework"]
            
            subgraph "Handlers"
                AuthHandler["Auth Handler<br>Go"]
                UserHandler["User Handler<br>Go"]
                PingHandler["Ping Handler<br>Go"]
            end
            
            subgraph "Middleware"
                MetricsMiddleware["Metrics Middleware<br>Go"]
            end
            
            Router["URL Router<br>Gin"]
        end
        
        subgraph "Service Layer"
            AuthService["Auth Service<br>Go"]
            UserService["User Service<br>Go"]
        end
        
        subgraph "Repository Layer"
            KeycloakRepo["Keycloak Repository<br>Go"]
            TeamRepo["Team Repository<br>Go"]
        end
        
        subgraph "Client Layer"
            RestClient["REST Client<br>Go HTTP"]
        end
    end
    
    subgraph "External Systems"
        KeycloakAuth["Keycloak Auth Server<br>Keycloak"]
        CowAPI["COW API<br>REST API"]
    end

    %% User interactions
    User -->|"Authenticates"| GinServer
    
    %% API Layer connections
    GinServer -->|"Routes requests"| Router
    Router -->|"Auth requests"| AuthHandler
    Router -->|"User requests"| UserHandler
    Router -->|"Health checks"| PingHandler
    GinServer -->|"Applies"| MetricsMiddleware
    
    %% Handler to Service connections
    AuthHandler -->|"Uses"| AuthService
    UserHandler -->|"Uses"| UserService
    
    %% Service to Repository connections
    AuthService -->|"Uses"| KeycloakRepo
    UserService -->|"Uses"| KeycloakRepo
    UserService -->|"Uses"| TeamRepo
    
    %% Repository to External Systems connections
    KeycloakRepo -->|"Authenticates/Manages users"| KeycloakAuth
    TeamRepo -->|"Fetches team data"| RestClient
    RestClient -->|"Makes HTTP requests"| CowAPI
    
    %% Main API endpoints
    Router -->|"POST /auth/login"| AuthHandler
    Router -->|"POST /auth/logout"| AuthHandler
    Router -->|"POST /auth/valid-token"| AuthHandler
    Router -->|"GET /users"| UserHandler
    Router -->|"GET /users/:code"| UserHandler
    Router -->|"POST /users"| UserHandler
    Router -->|"DELETE /users/:code"| UserHandler
    Router -->|"GET /ping"| PingHandler
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