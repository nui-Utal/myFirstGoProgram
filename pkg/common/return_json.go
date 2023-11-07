package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleError(c *gin.Context, wro int, ms string) {
	c.JSON(wro, gin.H{
		"code": 404,
		"mes":  ms,
		"data": "",
	})
}

func Success(c *gin.Context, ms string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"mes":  ms,
		"data": data,
	})
}
