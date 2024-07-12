package files

import (
	filesDomain "backend/models/files"
	fileService "backend/services/files_service"

	//"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	var fileReq filesDomain.UploadFile

	if err := c.ShouldBindJSON(&fileReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileService.UploadFile(fileReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, file)
}


func GetFile(c *gin.Context){
	cursoId := c.Param("curso_id")

	files, err := fileService.GetFile(cursoId)
	//log.Printf("files: %v", err);

	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, files)
}