package controllers

import (
	"admin/models"
	"admin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListarInstancias(c *gin.Context) {
    instancias, err := services.ObtenerInstancias()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las instancias"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "instances": instancias,
    })
}


func CrearInstancia(c *gin.Context) {
    var nuevaInstancia models.Instancia
    if err := c.ShouldBindJSON(&nuevaInstancia); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    err := services.CrearInstancia(nuevaInstancia)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la instancia"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Instancia creada correctamente"})
}

func EliminarInstancia(c *gin.Context) {
    var instancia models.Instancia
    if err := c.ShouldBindJSON(&instancia); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    err := services.EliminarInstancia(instancia)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la instancia"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Instancia eliminada correctamente"})
}
