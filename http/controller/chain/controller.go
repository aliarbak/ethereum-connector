package chain

import (
	"github.com/aliarbak/ethereum-connector/errors"
	httpconfig "github.com/aliarbak/ethereum-connector/http/config"
	"github.com/aliarbak/ethereum-connector/service/chain"
	chainmodel "github.com/aliarbak/ethereum-connector/service/chain/model"
	"github.com/labstack/echo/v4"
	"strconv"
)

type Controller interface {
}

type controller struct {
	chainService chain.Service
}

func New(chainService chain.Service, g *echo.Group) Controller {
	controller := controller{
		chainService: chainService,
	}

	g.POST("", controller.CreateChain)
	g.GET("", controller.GetChains)
	g.GET("/:id", controller.GetChainById)
	g.PATCH("/:id", controller.PatchChain)
	return &controller
}

// CreateChain godoc
// @Summary Register a new chain
// @Tags chains
// @produce application/json
// @Param request body chainmodel.CreateChainInput true "Request"
// @Success 201 {object} model.Chain
// @Router /chains [post]
func (c *controller) CreateChain(ctx echo.Context) error {
	var req chainmodel.CreateChainInput
	err := httpconfig.ReadValidatedBody(ctx, &req)
	if err != nil {
		return err
	}

	chain, err := c.chainService.CreateChain(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(201, chain)
}

// GetChains godoc
// @Summary Retrieve all chains
// @Tags chains
// @produce application/json
// @Success 200 {object} []model.Chain
// @Router /chains [get]
func (c *controller) GetChains(ctx echo.Context) error {
	chains, err := c.chainService.GetChains(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(200, chains)
}

// GetChainById godoc
// @Summary Retrieve chain by id
// @Tags chains
// @produce application/json
// @Param id path int64 true "id"
// @Success 200 {object} model.Chain
// @Router /chains/{id} [get]
func (c *controller) GetChainById(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return errors.InvalidInput("invalid id: %s", ctx.Param("id"))
	}

	chain, err := c.chainService.GetChainById(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(200, chain)
}

// PatchChain godoc
// @Summary Update a chain
// @Tags chains
// @produce application/json
// @Param id path int64 true "id"
// @Param request body chainmodel.PatchChainInput true "Request"
// @Success 200 {object} model.Chain
// @Router /chains/{id} [patch]
func (c *controller) PatchChain(ctx echo.Context) error {
	var req chainmodel.PatchChainInput
	err := httpconfig.ReadValidatedBody(ctx, &req)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return errors.InvalidInput("invalid id: %s", ctx.Param("id"))
	}

	chain, err := c.chainService.PatchChain(ctx.Request().Context(), id, req)
	if err != nil {
		return err
	}

	return ctx.JSON(200, chain)
}
