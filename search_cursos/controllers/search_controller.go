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
	Search(ctx context.Context, query string, offset int, limit int) ([]coursesDomain.Curso, error)
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
    offsetStr = strings.TrimSpace(offsetStr)  
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("invalid request: %s", err),
        })
        return
    }

    limitStr := c.DefaultQuery("limit", "10")
    limitStr = strings.TrimSpace(limitStr)  
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("invalid request: %s", err),
        })
        return
    }

    courses, err := controller.service.Search(c.Request.Context(), query, offset, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("error searching courses: %s", err.Error()),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "courses": courses,
    })
}
