package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/bright-pentium/go-client-practice/internal/configs"
	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type UserControler struct {
	usecase *usecase.UserUseCase
	config  *configs.AppConfig
}

func NewUserControler(usecase *usecase.UserUseCase, config *configs.AppConfig) *UserControler {
	return &UserControler{usecase: usecase, config: config}
}

func (u *UserControler) RegisterRoutes(e *echo.Echo) {
	e.POST("/auth/users/login", u.Login)
}

type UserLoginRequest struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token" validate:"required"`
}

// @Summary User login
// @Description User Login
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body UserLoginRequest true "User Login Request"
// @Success 200 {object} UserLoginResponse "Success"
// @Failure 400 {object} echo.HTTPError "Bad Requests"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /users/login [post]
func (u *UserControler) Login(ctx echo.Context) error {
	req := new(UserLoginRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := u.usecase.LoginUser(ctx.Request().Context(), req.Account, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserLoginFail) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// Create token with claims
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.JwtClaims{
		Name:  user.Name,
		Scope: string(domain.PermAll),
		Type:  domain.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(u.config.SecretExpiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    u.config.Issuer,
			Subject:   user.ID.String(),
		},
	})

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(u.config.SecretKey))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, UserLoginResponse{AccessToken: tokenString})
}
