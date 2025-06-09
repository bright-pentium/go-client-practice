package controller

import (
	"errors"
	"net/http"

	"github.com/bright-pentium/go-client-practice/internal/configs"
	"github.com/bright-pentium/go-client-practice/internal/domain"
	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SysUserControler struct {
	usecase *usecase.SysUserUseCase
	config  *configs.AppConfig
}

func NewSysUserControler(usecase *usecase.SysUserUseCase, config *configs.AppConfig) *SysUserControler {
	return &SysUserControler{usecase: usecase, config: config}
}

func (u *SysUserControler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/admin/users")
	api.GET("/:user-id", u.GetUserByID)
	api.PATCH("/:user-id", u.UpdateUserByID)
	api.DELETE("/:user-id", u.DeleteUserByID)
	api.POST("", u.CreateUser)
}

// CreateUserRequest defines the request body for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary Create a User
// @Description Creates a new user.
// @Tags admin
// @Accept  json
// @Produce  json
// @Param request body CreateUserRequest true "User creation request"
// @Success 200 {object} domain.User "Success"
// @Success 400 {object} echo.HTTPError "Bad Requests"
// @Success 500 {object} echo.HTTPError "Internal Error"
// @Router /admin/users [post]
func (u *SysUserControler) CreateUser(ctx echo.Context) error {
	req := new(CreateUserRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := u.usecase.CreateUser(ctx.Request().Context(), req.Name, req.Account, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) || errors.Is(err, domain.ErrInvalidUserData) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusCreated, user)
}

// @Summary Get User by ID
// @Description Retrieves a user by their ID.
// @Tags admin
// @Accept  json
// @Produce  json
// @Param user-id path string true "User ID"
// @Success 200 {object} domain.User "Success"
// @Failure 400 {object} echo.HTTPError "Bad Requests"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /admin/users/{user-id} [get]
func (u *SysUserControler) GetUserByID(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := u.usecase.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, user)
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// @Summary Update User by ID
// @Description Update a user by their ID.
// @Tags admin
// @Accept json
// @Produce json
// @Param user-id path string true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} domain.User "Success"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /admin/users/{user-id} [patch]
func (u *SysUserControler) UpdateUserByID(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req := new(UpdateUserRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := u.usecase.UpdateUserByID(ctx.Request().Context(), userID, req.Name, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, user)
}

// @Summary Delete User by ID
// @Description Delete a user by their ID.
// @Tags admin
// @Accept  json
// @Produce  json
// @Param user-id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /admin/users/{user-id} [delete]
func (u *SysUserControler) DeleteUserByID(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = u.usecase.DeleteUserByID(ctx.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.NoContent(http.StatusNoContent)
}
