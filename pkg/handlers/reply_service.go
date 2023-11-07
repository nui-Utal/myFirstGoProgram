package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xinyeOfficalWebsite/pkg/common"
	"xinyeOfficalWebsite/pkg/models"
	"xinyeOfficalWebsite/pkg/utils"
)

func Reply(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	var r models.Reply
	if err := c.BindJSON(&r); err != nil {
		return
	}
	if "" == r.Content {
		common.HandleError(c, http.StatusBadRequest, "回复内容不能为空")
		return
	}
	r.Send = uid
	r.Reply_time = utils.GetLocalTime()
	if !models.AddReply(r) {
		common.HandleError(c, http.StatusBadRequest, "回复失败")
		return
	}
	// 通过评论得到文章id，文章的评论数++
	com := models.GetCommentById(r.Comment)
	common.IncrementField(common.TableText, common.TextFieldComment, com.Text)
	common.IncrementField(common.TableComment, common.CommentFieldReplied, com.ID)
	common.Success(c, "回复成功", "")
}

func DelReply(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}

	replyid := c.PostForm("replyid")
	comid := c.PostForm("comid")
	if "" == comid || "" == replyid {
		common.HandleError(c, http.StatusBadRequest, "未接收到回复id")
		return
	}
	cid, err := strconv.Atoi(comid)
	if err != nil {
		return
	}
	rid, err := strconv.Atoi(replyid)
	if err != nil {
		return
	}

	// 身份判断
	if uid != models.GetReplyById(rid).Send {
		common.HandleError(c, http.StatusUnauthorized, "无法删除该评论")
		return
	}
	// 删除记录
	if err := models.DelReplyById(rid); err != nil {
		common.HandleError(c, http.StatusNotFound, err.Error())
		return
	}
	// 减少关联文章评论数
	com := models.GetCommentById(cid)
	common.DecrementField(common.TableText, common.TextFieldComment, com.Text)
	common.DecrementField(common.TableComment, common.CommentFieldReplied, com.Text)
	common.Success(c, "删除成功", "")
}
