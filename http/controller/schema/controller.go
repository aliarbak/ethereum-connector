package schema

import (
	httpconfig "github.com/aliarbak/ethereum-connector/http/config"
	"github.com/aliarbak/ethereum-connector/service/schema"
	schemamodel "github.com/aliarbak/ethereum-connector/service/schema/model"
	"github.com/labstack/echo/v4"
)

type Controller interface {
}

type controller struct {
	schemaService schema.Service
}

func New(schemaService schema.Service, g *echo.Group) Controller {
	controller := controller{
		schemaService: schemaService,
	}

	g.POST("", controller.CreateSchema)
	g.GET("", controller.GetSchemas)
	g.GET("/:id", controller.GetSchemaById)
	g.POST("/:id/filters", controller.CreateSchemaFilter)
	g.GET("/:id/filters", controller.GetSchemaFilters)
	g.GET("/:id/filters/:filterId", controller.GetSchemaFilterById)
	return &controller
}

// CreateSchema godoc
// @Summary Create a new schema
// @Tags schemas
// @produce application/json
// @Param request body schemamodel.CreateSchemaInput true "Request"
// @Success 201 {object} schemamodel.SchemaOutput
// @Router /schemas [post]
func (c *controller) CreateSchema(ctx echo.Context) error {
	var req schemamodel.CreateSchemaInput
	err := httpconfig.ReadValidatedBody(ctx, &req)
	if err != nil {
		return err
	}

	schema, err := c.schemaService.CreateSchema(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(201, schema)
}

// GetSchemas godoc
// @Summary Retrieve all schemas
// @Tags schemas
// @produce application/json
// @Success 200 {object} []schemamodel.SchemaOutput
// @Router /schemas [get]
func (c *controller) GetSchemas(ctx echo.Context) error {
	schemas, err := c.schemaService.GetSchemas(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(200, schemas)
}

// GetSchemaById godoc
// @Summary Retrieve schema by id
// @Tags schemas
// @produce application/json
// @Param id path string true "id"
// @Success 200 {object} schemamodel.SchemaOutput
// @Router /schemas/{id} [get]
func (c *controller) GetSchemaById(ctx echo.Context) error {
	id := ctx.Param("id")
	schema, err := c.schemaService.GetSchemaById(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(200, schema)
}

// CreateSchemaFilter godoc
// @Summary Create a new schema filter
// @Tags schemas
// @produce application/json
// @Param id path string true "id"
// @Param request body schemamodel.CreateSchemaFilterInput true "Request"
// @Success 201 {object} model.SchemaFilter
// @Router /schemas/{id}/filters [post]
func (c *controller) CreateSchemaFilter(ctx echo.Context) error {
	var req schemamodel.CreateSchemaFilterInput
	err := httpconfig.ReadValidatedBody(ctx, &req)
	if err != nil {
		return err
	}

	id := ctx.Param("id")
	filter, err := c.schemaService.CreateSchemaFilter(ctx.Request().Context(), id, req)
	if err != nil {
		return err
	}

	return ctx.JSON(201, filter)
}

// GetSchemaFilters godoc
// @Summary Retrieve schema filters by schema id
// @Tags schemas
// @produce application/json
// @Param id path string true "id"
// @Success 200 {object} []model.SchemaFilter
// @Router /schemas/{id}/filters [get]
func (c *controller) GetSchemaFilters(ctx echo.Context) error {
	id := ctx.Param("id")
	filter, err := c.schemaService.GetSchemaFilters(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(200, filter)
}

// GetSchemaFilterById godoc
// @Summary Retrieve schema filters by schema filter id
// @Tags schemas
// @produce application/json
// @Param id path string true "id"
// @Param filterId path string true "filterId"
// @Success 200 {object} model.SchemaFilter
// @Router /schemas/{id}/filters/{filterId} [get]
func (c *controller) GetSchemaFilterById(ctx echo.Context) error {
	id := ctx.Param("id")
	filterId := ctx.Param("filterId")
	filter, err := c.schemaService.GetSchemaFilterById(ctx.Request().Context(), id, filterId)
	if err != nil {
		return err
	}

	return ctx.JSON(200, filter)
}
