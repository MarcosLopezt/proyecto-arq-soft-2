package controllers

import (
	"admin/models"
	"admin/services"
	"net/http"
	"strconv"

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

func ProbarBalanceoCarga(c *gin.Context) {
    // Obtener el número de peticiones desde query parameter, por defecto 10
    numPeticionesStr := c.DefaultQuery("num_peticiones", "10")
    numPeticiones, err := strconv.Atoi(numPeticionesStr)
    if err != nil || numPeticiones <= 0 {
        numPeticiones = 10
    }
    
    // Limitar el número máximo de peticiones para evitar sobrecarga
    if numPeticiones > 50 {
        numPeticiones = 50
    }
    
    resultados, err := services.RealizarPruebaBalanceo(numPeticiones)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al realizar la prueba de balanceo: " + err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "resultados": resultados,
        "total_peticiones": len(resultados),
        "resumen": generarResumen(resultados),
    })
}

// generarResumen crea un resumen de los resultados de la prueba
func generarResumen(resultados []services.ResultadoPruebaBalanceo) map[string]interface{} {
    contadorPuertos := make(map[string]int)
    
    for _, resultado := range resultados {
        contadorPuertos[resultado.Puerto]++
    }
    
    return map[string]interface{}{
        "distribucion_puertos": contadorPuertos,
        "total_instancias": len(contadorPuertos),
    }
}
