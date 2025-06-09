package controller

import (
	"net/http"

	"github.com/bright-pentium/go-client-practice/internal/configs"
	"github.com/bright-pentium/go-client-practice/internal/delivery/echo/middleware"
	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

type ResourceControler struct {
	usecase *usecase.ResourceUseCase
	config  *configs.AppConfig
}

func NewResourceControler(usecase *usecase.ResourceUseCase, config *configs.AppConfig) *ResourceControler {
	return &ResourceControler{usecase: usecase, config: config}
}

func (r *ResourceControler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/resources", echojwt.WithConfig(echojwt.Config{
		// ...
		SigningKey: []byte(r.config.SecretKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(domain.JwtClaims)
		},
		// ...
	}))
	api.Use(middleware.JWTMiddleware)
	api.Use(middleware.RequirePermission(domain.PermCreateResource))
	api.POST("", r.CreateResource)
}

// @Summary Create Resource
// @Description Create Resource
// @Tags resourece
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.Resource "Success"
// @Failure 400 {object} echo.HTTPError "Bad Requests"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Security Bearer
// @Router /resources  [post]
func (r *ResourceControler) CreateResource(ctx echo.Context) error {
	resource, _ := r.usecase.CreateResource(ctx.Request().Context())
	return ctx.JSON(http.StatusOK, resource)
}
