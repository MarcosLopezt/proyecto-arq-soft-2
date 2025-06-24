package controllers

import (
	"context"
	"fmt"
	"net/http"
	coursesDomain "search_cursos/domain"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Search(ctx context.Context, query string, offset int, limit int, availableOnly bool) ([]coursesDomain.Curso, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) Search(c *gin.Context) {
    query := c.DefaultQuery("q", "*:*")
    offsetStr := c.DefaultQuery("offset", "0")
    offset, err := strconv.Atoi(strings.TrimSpace(offsetStr))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid offset: %s", err)})
        return
    }

    limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(strings.TrimSpace(limitStr))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid limit: %s", err)})
        return
    }

    availableOnlyStr := c.DefaultQuery("availableOnly", "false")
    availableOnly := false
    if strings.ToLower(availableOnlyStr) == "true" {
        availableOnly = true
    }

    courses, err := controller.service.Search(c.Request.Context(), query, offset, limit, availableOnly)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error searching courses: %s", err.Error())})
        return
    }

    c.JSON(http.StatusOK, gin.H{"courses": courses})
}
