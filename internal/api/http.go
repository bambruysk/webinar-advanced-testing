package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"webinar-testing/pkg/models"
)

//go:generate mockery --name Service --with-expecter
type Service interface {
	Add(ctx context.Context, goods models.Order) error
	ListByUser(ctx context.Context, id models.UserID) (models.Order, error)
	Delete(ctx context.Context, goods models.Order) error
	DeleteAllByUser(ctx context.Context, id models.UserID) error
}

type server struct {
	serv    *echo.Echo
	service Service
}

func NewServer(service Service) *server {
	e := echo.New()
	s := &server{
		serv:    e,
		service: service,
	}

	e.POST("/add", s.Add)
	e.GET("/list/:id", s.List)
	e.PUT("/remove", s.Delete)
	e.DELETE("/delete", s.DeleteAll)

	return s
}

func (s *server) Run() error {
	if err := s.serv.Start(":8080"); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *server) Add(c echo.Context) error {
	order := models.Order{}
	if err := c.Bind(&order); err != nil {
		c.Logger().Error(err)
		return err
	}
	ctx := c.Request().Context()
	if err := s.service.Add(ctx, order); err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *server) Delete(c echo.Context) error {
	order := models.Order{}
	if err := c.Bind(&order); err != nil {
		c.Logger().Error(err)
		return err
	}

	err := s.service.Delete(c.Request().Context(), order)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *server) List(c echo.Context) error {
	id := c.Param("id")

	resp, err := s.service.ListByUser(c.Request().Context(), models.UserID(id))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *server) DeleteAll(c echo.Context) error {
	id := c.Param("id")

	err := s.service.DeleteAllByUser(c.Request().Context(), models.UserID(id))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.NoContent(http.StatusOK)
}
