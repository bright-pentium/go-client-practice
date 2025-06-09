package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/bright-pentium/go-client-practice/internal/configs"
	"github.com/bright-pentium/go-client-practice/internal/delivery/echo/middleware"
	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

type ClientController struct {
	usecase *usecase.ClientUseCase
	config  *configs.AppConfig
}

func NewClientController(usecase *usecase.ClientUseCase, config *configs.AppConfig) *ClientController {
	return &ClientController{usecase: usecase, config: config}
}

func (c *ClientController) RegisterRoutes(e *echo.Echo) {

	api := e.Group("/clients", echojwt.WithConfig(echojwt.Config{
		// ...
		SigningKey: []byte(c.config.SecretKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(domain.JwtClaims)
		},
		// ...
	}))
	api.Use(middleware.JWTMiddleware)
	api.GET("", c.ListClientsByUser)
	api.POST("", c.CreateClient)

	e.POST("/clients/login", c.ClientLogin)

	// api.GET("/:client-id", c.GetClientByID)
	// api.PATCH("/:client-id", c.UpdateClientByIDandUser)
	// api.DELETE("/:client-id", c.DeleteClientByIDandUser)

}

type ClientResponse struct {
	Client []domain.Client `json:"client"`
}

// @Summary List Clients by User ID
// @Description Retrieves all clients associated with user ID.
// @Tags client
// @Accept  json
// @Produce  json
// @Success 200 {array} ClientResponse "Success"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Security Bearer
// @Router /clients [get]
func (c *ClientController) ListClientsByUser(ctx echo.Context) error {
	// Extract JWT token
	userID, _ := ctx.Get("userID").(uuid.UUID)
	// Fetch clients from usecase
	clients, err := c.usecase.ListClientsByUser(ctx.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, clients)
}

type CreateClientRequest struct {
	Scope []domain.Permission `json:"scope" example:"resource:create" validate:"required,min=1,dive,required,perm"`
}

type CreateClientReponse struct {
	domain.Client
	Secret string `json:"secret" example:"o44z4KUzru7uW4jtzxVt84Ma8f76Mnwj"`
}

// @Summary Create Client
// @Description Creates a new client associated with the authenticated user.
// @Tags client
// @Accept  json
// @Produce  json
// @Param request body CreateClientRequest true "Create Client Request"
// @Success 201 {object} CreateClientReponse "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Security Bearer
// @Router /clients [post]
func (c *ClientController) CreateClient(ctx echo.Context) error {
	userID, _ := ctx.Get("userID").(uuid.UUID)
	req := new(CreateClientRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, secret, err := c.usecase.CreateClient(ctx.Request().Context(), uuid.New(), userID, req.Scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(
		http.StatusCreated,
		CreateClientReponse{
			Client: *client,
			Secret: secret,
		},
	)
}

type ClientLoinRequest struct {
	Secret string    `json:"secret" example:"o44z4KUzru7uW4jtzxVt84Ma8f76Mnwj"`
	ID     uuid.UUID `json:"id" example:"6cc2b688-1246-4a62-a293-dae7e67d6097"`
}

type ClientLoginReponse struct {
	AccessToken string `json:"access_token"`
}

// @Summary Create Client
// @Description Creates a new client associated with the authenticated user.
// @Tags client
// @Accept  json
// @Produce  json
// @Param request body ClientLoinRequest true "Create Client Request"
// @Success 201 {object} ClientLoginReponse "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /clients/login [post]
func (c *ClientController) ClientLogin(ctx echo.Context) error {
	req := new(ClientLoinRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := c.usecase.ClientLogin(ctx.Request().Context(), req.ID, req.Secret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	scope := strings.Join(Map(client.Scope, func(p domain.Permission) string {
		return string(p)
	}), " ")

	// Create token with claims
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.JwtClaims{
		Name:  "",
		Scope: scope,
		Type:  domain.ClientType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(c.config.SecretExpiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    c.config.Issuer,
			Subject:   client.ID.String(),
		},
	})

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(c.config.SecretKey))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, UserLoginResponse{AccessToken: tokenString})
}

type UpdateClientRequest struct {
	Scope []domain.Permission `json:"scope" example:"resource:create" validate:"required,min=1,dive,required,perm"`
}

// @Summary Update Client
// @Description Update a client with scopes.
// @Tags client
// @Accept  json
// @Produce  json
// @Param request body UpdateClientRequest true "Update Client Request"
// @Success 200 {object} domain.Client "Success"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /clients/{client-id} [patch]
func (c *ClientController) UpdateClientByIDandUser(ctx echo.Context) error {
	userID, _ := ctx.Get("userID").(uuid.UUID)
	clientID, err := uuid.Parse(ctx.Param("client-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := new(UpdateClientRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := c.usecase.UpdateClientScope(ctx.Request().Context(), clientID, userID, req.Scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, client)
}

func Map[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}
