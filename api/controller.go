package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nomionz/ctc-api/models"
	"github.com/nomionz/ctc-api/repositories"
	"gorm.io/gorm"
)

type Controller struct {
	router *gin.Engine
	repo   repositories.Repository
}

func NewController(repo repositories.Repository) *Controller {
	c := new(Controller)
	c.router = c.setupRoutes()
	c.repo = repo
	return c
}

func (c *Controller) Run(port string) error {
	return c.router.Run(port)
}

func (c *Controller) setupRoutes() *gin.Engine {
	router := gin.Default()
	prods := router.Group("/products")
	{
		prods.GET("/", c.list)
		prods.GET("/:id", c.get)
		prods.POST("/", c.create)
		prods.PATCH("/:id", c.update)
		prods.DELETE("/:id", c.delete)
	}
	return router
}

func parseId(ctx *gin.Context) (int, error) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	if id < 0 {
		return 0, fmt.Errorf("id is negative")
	}
	return id, nil
}

func (c *Controller) delete(ctx *gin.Context) {
	id, err := parseId(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err := c.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("product id %d deleted", id)})
}

func (c *Controller) update(ctx *gin.Context) {
	id, err := parseId(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	var prod models.Product
	prod.ID = id
	if err := ctx.ShouldBindJSON(&prod); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err := models.Validate(&prod); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err := c.repo.Update(&prod); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, err)
		}
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, prod)
}

func (c *Controller) get(ctx *gin.Context) {
	id, err := parseId(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	res, err := c.repo.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, err)
		}
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) list(ctx *gin.Context) {
	res, err := c.repo.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) create(ctx *gin.Context) {
	var prod models.Product
	if err := ctx.ShouldBindJSON(&prod); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err := models.Validate(&prod); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err := c.repo.Create(&prod); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, prod)
}
