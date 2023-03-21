package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	srv  *gin.Engine
	repo *ProductRepository
}

func NewController(repo *ProductRepository) *Controller {
	c := new(Controller)
	c.srv = gin.Default()
	c.repo = repo
	c.setupRoutes()
	return c
}

func (c *Controller) Run(port string) error {
	return c.srv.Run(port)
}

func (c *Controller) setupRoutes() {
	prods := c.srv.Group("/products")
	{
		prods.GET("/list", c.list)
		prods.GET("/:id", c.get)
	}
}

func (c *Controller) get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	res, err := c.repo.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) list(ctx *gin.Context) {
	res, err := c.repo.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
	}
	ctx.JSON(http.StatusOK, res)
}
