package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"xinyeOfficalWebsite/pkg/common"
	"xinyeOfficalWebsite/pkg/models"
	"xinyeOfficalWebsite/pkg/utils"
)

type comment struct {
	UserId int    `json:"userid"`
	Page   int    `json:"page"`
	Time   string `json:"time"`
}

func bindComJSON(c *gin.Context, data *models.Comment) error {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

// 查看我的评论
func CheckMyComment(c *gin.Context) {
	// 获取数据
	var com comment
	if err := c.ShouldBindJSON(&com); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	time := strings.Split(com.Time, "-")

	var mycom []models.Comment
	mycom = models.GetCommentByUidLimitTime(com.UserId, time[0], time[1])
	var myrep []models.Reply
	myrep = models.GetReplyByUidLimitTime(com.UserId, time[0], time[1])

	common.Success(c, "", mycom) // 响应给浏览器的评论数据
	common.Success(c, "", myrep) // 响应给浏览器的回复数据
}

// 评论文章
func CommentText(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	var com models.Comment
	err = bindComJSON(c, &com)
	if err != nil {
		return
	}

	com.Commentator_id = uid
	com.Commentator_name = models.GetUserById(uid).Username
	com.CommentTime = utils.GetLocalTime()
	if err := models.AddComment(com); err != nil {
		common.HandleError(c, http.StatusBadRequest, "评论失败")
		return
	}
	common.IncrementField(common.TableText, common.TextFieldComment, com.Text)
	common.Success(c, "评论成功", "")
}

func DelComment(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	comid := c.PostForm("comid")
	textid := c.PostForm("textid")
	if "" == comid || "" == textid {
		common.HandleError(c, http.StatusBadRequest, "未接收到评论id")
		return
	}
	cid, err := strconv.Atoi(comid)
	if err != nil {
		return
	}
	tid, err := strconv.Atoi(textid)
	if err != nil {
		return
	}
	// 身份验证
	commented, comuserid := models.IsCommented(cid)
	if !commented {
		common.HandleError(c, http.StatusUnauthorized, "该评论不存在")
		return
	}
	if uid != comuserid {
		common.HandleError(c, http.StatusUnauthorized, "无法删除该评论")
		return
	}
	common.DecrementField(common.TableText, common.TextFieldComment, tid)
	if models.DelCommentById(cid) {
		common.Success(c, "删除成功", "")
	}
}

func GetTextComment(c *gin.Context) {
	textid := c.Query("textid")
	coms := models.GetCommentByTid(textid)
	common.Success(c, "", coms)
	var rep []models.Reply
	for _, com := range coms {
		if com.Replied != 0 {
			rep = append(rep, models.GetReplyByCid(com.ID)...)
		}
	}
	if rep != nil {
		common.Success(c, "", rep)
	}
}
