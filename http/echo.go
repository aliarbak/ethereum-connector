package http

import (
	"context"
	"fmt"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/datasource"
	"github.com/aliarbak/ethereum-connector/destination"
	chaincontroller "github.com/aliarbak/ethereum-connector/http/controller/chain"
	schemacontroller "github.com/aliarbak/ethereum-connector/http/controller/schema"
	middlewares "github.com/aliarbak/ethereum-connector/http/middleware"
	"github.com/aliarbak/ethereum-connector/service/chain"
	"github.com/aliarbak/ethereum-connector/service/schema"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/ziflex/lecho/v3"
	"net/http"
)

// @title Ethereum Connector
// @version 1.0
// @BasePath /
// @schemes http https
func Run(ctx context.Context, config configs.HttpConfig, dataSourceFactory datasource.Factory, destinationFactory destination.Factory) error {
	e := echo.New()

	// Middleware
	e.Logger = lecho.From(log.Logger)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.HTTPErrorHandler = middlewares.ErrorHandler

	// Services
	chainService, err := chain.NewService(ctx, dataSourceFactory, destinationFactory)
	if err != nil {
		return err
	}
	defer chainService.Close(ctx)

	schemaService, err := schema.NewService(ctx, dataSourceFactory, destinationFactory)
	if err != nil {
		return err
	}
	defer schemaService.Close(ctx)

	// Routes
	e.GET("/healthcheck", HealthCheck)
	if config.SwaggerEnabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	chainsGroup := e.Group("/chains")
	_ = chaincontroller.New(chainService, chainsGroup)

	schemasGroup := e.Group("/schemas")
	_ = schemacontroller.New(schemaService, schemasGroup)

	// Start server
	return e.Start(fmt.Sprintf(":%s", config.Port))
}

// healthCheck godoc
// @Summary Show the status of server
// @Tags HealthCheck
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /healthcheck [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Ethereum connector is up and running",
	})
}
