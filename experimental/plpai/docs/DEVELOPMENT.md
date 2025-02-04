# REST API Project Structure

```
├── cmd/
│   └── api/
│       └── main.go           # API server entry point
├── internal/
│   ├── api/                  # API-specific code
│   │   ├── handlers/        # HTTP request handlers
│   │   │   ├── users.go
│   │   │   ├── auth.go
│   │   │   └── health.go
│   │   ├── middleware/      # HTTP middleware components
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── cors.go
│   │   ├── router/         # Router setup and configuration
│   │   │   └── router.go
│   │   └── responses/      # Standard API responses
│   │       └── errors.go
│   ├── domain/             # Business logic and entities
│   │   ├── models/        # Domain models/entities
│   │   └── services/      # Business logic services
│   ├── repository/         # Data access layer
│   │   ├── postgres/      # PostgreSQL implementations
│   │   └── interfaces.go  # Repository interfaces
│   └── config/            # Application configuration
├── pkg/
│   ├── validator/         # Request validation helpers
│   └── httputils/         # Common HTTP utilities
├── api/
│   ├── swagger/           # Swagger/OpenAPI documentation
│   │   └── swagger.yaml
│   └── proto/            # If using gRPC alongside REST
├── configs/
│   ├── app.env           # Application configuration
│   └── app.yaml
└── [rest of standard structure...]
```

## Key REST API Components

1. **internal/api/**
   - **handlers/**: HTTP request handlers implementing your API endpoints
     ```go
     // handlers/users.go
     package handlers

     type UserHandler struct {
         userService services.UserService
     }

     func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
         // Handle user creation
     }

     func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
         // Handle get user
     }
     ```

   - **middleware/**: Common middleware functions
     ```go
     // middleware/auth.go
     package middleware

     func AuthMiddleware(next http.Handler) http.Handler {
         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
             // Authenticate request
         })
     }
     ```

   - **router/**: API route definitions
     ```go
     // router/router.go
     package router

     func Setup(handlers *handlers.Handlers) *chi.Router {
         r := chi.NewRouter()
         
         r.Use(middleware.Logger)
         r.Use(middleware.Recoverer)
         
         r.Route("/api/v1", func(r chi.Router) {
             r.Post("/users", handlers.Users.Create)
             r.Get("/users/{id}", handlers.Users.Get)
         })
         
         return r
     }
     ```

2. **internal/domain/services/**
   ```go
   // services/user_service.go
   package services

   type UserService interface {
       Create(ctx context.Context, user *models.User) error
       Get(ctx context.Context, id string) (*models.User, error)
   }
   ```

3. **pkg/validator/**
   ```go
   // validator/validator.go
   package validator

   type RequestValidator struct {
       validator *validator.Validate
   }

   func (v *RequestValidator) ValidateUser(user *models.User) error {
       return v.validator.Struct(user)
   }
   ```

## Main Application Setup

```go
// cmd/api/main.go
package main

func main() {
    // Initialize configuration
    cfg := config.Load()

    // Setup database connection
    db := postgres.NewConnection(cfg.DatabaseURL)

    // Initialize repositories
    userRepo := postgres.NewUserRepository(db)

    // Initialize services
    userService := services.NewUserService(userRepo)

    // Setup handlers
    handlers := handlers.New(userService)

    // Setup router
    router := router.Setup(handlers)

    // Start server
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## API Documentation

1. **OpenAPI/Swagger Specification**
   - Place your API specifications in `api/swagger/`
   - Use tools like `swag` to generate documentation from code comments
   ```go
   // handlers/users.go
   // @Summary Create user
   // @Description Create a new user
   // @Tags users
   // @Accept json
   // @Produce json
   // @Param user body models.User true "User object"
   // @Success 201 {object} models.User
   // @Router /users [post]
   func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```


## Best Practices for REST APIs

1. **Versioning**
   - Version your API in the URL path (e.g., `/api/v1/users`)
   - Consider using content negotiation for versioning

2. **Error Handling**
   - Use consistent error responses
   - Include error codes and messages
   ```go
   // responses/errors.go
   type ErrorResponse struct {
       Code    string `json:"code"`
       Message string `json:"message"`
       Details any    `json:"details,omitempty"`
   }
   ```

3. **Middleware**
   - Use middleware for cross-cutting concerns
   - Common middleware: logging, authentication, CORS, request ID

4. **Testing**
   - Write handler tests using `httptest`
   - Use table-driven tests for different scenarios
   - Test both success and error cases


### Best Practices for Development in Go

1. **Package Organization**
   - Keep packages small and focused
   - Use meaningful package names
   - Avoid package name collisions
   - Follow the "Clean Architecture" principles

2. **Dependency Management**
   - Use Go modules for dependency management
   - Consider vendoring dependencies for production builds
   - Keep dependencies up to date and secure

3. **Configuration**
   - Use environment variables for configuration
   - Keep secrets separate from code
   - Use configuration files for default values

4. **Testing**
   - Write unit tests in the same package as the code
   - Place integration tests in the `test/` directory
   - Aim for high test coverage
   - Use table-driven tests when appropriate

5. **Documentation**
   - Document all exported types and functions
   - Include examples in documentation
   - Keep API documentation up to date
   - Document configuration options


