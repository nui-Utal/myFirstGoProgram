package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xinyeOfficalWebsite/pkg/common"
	"xinyeOfficalWebsite/pkg/models"
)

func GetCollectAll(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	page, pagesize, err := common.GetPageAndPageSize(c)
	if err != nil {
		return
	}
	cols := models.GetCollectByUid(uid, page, pagesize)
	if len(cols) == 0 {
		common.Success(c, "", "您还未收藏")
		return
	}
	var t []models.Text
	for _, col := range cols {
		var text models.Text
		if !models.GetTextById(col.TextId, &text) {
			err = fmt.Errorf("您收藏的部分文章已被删除")
		} else {
			t = append(t, text)
		}
	}
	if err != nil {
		common.HandleError(c, http.StatusNotFound, err.Error())
	}
	common.Success(c, "", t)
}

func GetLikeAll(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	page, pagesize, err := common.GetPageAndPageSize(c)
	if err != nil {
		return
	}
	likes := models.GetLikeByUid(uid, page, pagesize)
	if len(likes) == 0 {
		common.Success(c, "您还未喜欢", "")
	}
	var t []models.Text
	for _, like := range likes {
		var text models.Text
		if !models.GetTextById(like.Text, &text) {
			err = fmt.Errorf("您喜欢的部分文章已被删除")
		} else {
			t = append(t, text)
		}
	}
	if err != nil {
		common.HandleError(c, http.StatusNotFound, err.Error())
	}
	common.Success(c, "", t)
}
