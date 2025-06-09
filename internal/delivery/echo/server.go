// @title My Echo API
// @version 1.0
// @description API documentation
// @BasePath /
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/bright-pentium/go-client-practice/docs/swaggo"
	"github.com/bright-pentium/go-client-practice/internal/configs"
	"github.com/bright-pentium/go-client-practice/internal/delivery/echo/controller"
	"github.com/bright-pentium/go-client-practice/internal/domain"
	clientRepo "github.com/bright-pentium/go-client-practice/internal/repository/client"
	userRepo "github.com/bright-pentium/go-client-practice/internal/repository/user"

	"github.com/bright-pentium/go-client-practice/internal/usecase"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type EchoServer struct {
	echo   *echo.Echo
	config *configs.AppConfig
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err == nil {
		return nil
	}

	// Only handle validation errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldErr := range validationErrors {
			switch fieldErr.Tag() {
			case "perm":
				messages = append(messages, fmt.Sprintf("Invalid permission: '%v'", fieldErr.Value()))
			default:
				messages = append(messages, fmt.Sprintf("Field '%s' failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag()))
			}
		}
		return fmt.Errorf(strings.Join(messages, ", "))
	}

	// return as-is for non-validation errors
	return err
}

func validPermission(fl validator.FieldLevel) bool {
	perm, ok := fl.Field().Interface().(domain.Permission)
	if !ok {
		return false
	}
	_, exists := domain.ValidPermissions[perm]
	return exists
}

func NewServer(config *configs.AppConfig) *EchoServer {
	e := echo.New()
	v := validator.New()
	v.RegisterValidation("perm", validPermission)
	e.Validator = &CustomValidator{validator: v}
	return &EchoServer{echo: e, config: config}
}

func (s *EchoServer) Serving(ctx context.Context) error {
	// TODO(bright) refactor config and Serving logic
	// suppport multiple db configs, some other repo may not need pxpool
	config, err := pgxpool.ParseConfig(s.config.DbURL)
	if err != nil {
		return err
	}

	config.MaxConns = int32(s.config.MaxConn)
	config.MinConns = int32(s.config.MinConn)
	pgxpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}
	defer pgxpool.Close()
	err = pgxpool.Ping(ctx)
	if err != nil {
		return err
	}

	userRepo := userRepo.NewPgxUserRepository(pgxpool)
	clientRepo := clientRepo.NewPgxClientRepository(pgxpool)

	SysUserUseCase := usecase.NewSysUserUseCase(userRepo)
	userUsecase := usecase.NewUserUseCase(userRepo)
	clientUsecase := usecase.NewClientUseCase(clientRepo)
	resourcetUsecase := usecase.NewResourceUseCase()

	sysUserControler := controller.NewSysUserControler(SysUserUseCase, s.config)
	sysUserControler.RegisterRoutes(s.echo)

	userControler := controller.NewUserControler(userUsecase, s.config)
	userControler.RegisterRoutes(s.echo)

	clientControler := controller.NewClientController(clientUsecase, s.config)
	clientControler.RegisterRoutes(s.echo)

	resourceControler := controller.NewResourceControler(resourcetUsecase, s.config)
	resourceControler.RegisterRoutes(s.echo)

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// Channel to capture server start errors
	errChan := make(chan error, 1)

	// Start Echo server in a goroutine
	go func() {
		errChan <- s.echo.Start(fmt.Sprintf(":%d", s.config.Port))
	}()

	select {
	case <-ctx.Done():
		// Context cancelled: initiate graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.echo.Shutdown(shutdownCtx)
	case err := <-errChan:
		// Server stopped with error
		return err
	}
}
