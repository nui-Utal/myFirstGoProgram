package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPageAndPageSize(c *gin.Context) (page, pageSize int, err error) {
	page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		HandleError(c, http.StatusOK, "请输入正确的页码")
		return 0, 0, err
	}

	pageSize, err = strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return page, 15, nil
	}
	return page, pageSize, nil
}
