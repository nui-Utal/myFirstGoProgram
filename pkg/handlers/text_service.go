package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"xinyeOfficalWebsite/pkg/common"
	"xinyeOfficalWebsite/pkg/models"
	"xinyeOfficalWebsite/pkg/utils"
)

type TextData struct {
	TextId  int    `json:"textid"`
	Title   string `json:"title"`
	Label   string `json:"label"`
	Content string `json:"content"`
	Type    int    `json:"type"`
}

type Page struct {
	//Userid int    `json:"userid"`
	Time     string `json:"time"`
	Page     int    `json:"page"`
	PageSize int    `json:"pagesize"`
}

// 绑定数据
func bindJSON(c *gin.Context, data *TextData) error {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

func CollectText(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	textid := c.PostForm("textid")
	if "" == textid {
		common.HandleError(c, http.StatusBadRequest, "未接受到文章id")
		return
	}
	tid, _ := strconv.Atoi(textid)

	// 判断是否收藏过
	if models.IsCollected(uid, tid) {
		common.HandleError(c, http.StatusBadRequest, "您已收藏过该文章")
		return
	}

	// text表和收藏表
	if err := common.IncrementField(common.TableText, common.TextFieldCollect, tid); err != nil {
		common.HandleError(c, http.StatusNotFound, "收藏失败")
		return
	}
	if !models.AddCollect(uid, tid) {
		common.HandleError(c, http.StatusBadRequest, "收藏失败")
	}
	common.Success(c, "收藏成功", "")
}

// 上传图片
func UploadPicture(c *gin.Context) {
	fileName, err := utils.GenerateSalt(16)
	fileName += ".png"
	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "图片上传失败")
	}
	if err = common.UploadPicture(c, "text", fileName); err != nil {
		return
	}

	//imagePath := template.HTMLEscapeString("<img src=\"" + common.GetPath("essay") + fileName + "\">")
	common.Success(c, "图片上传成功", "<img src=\""+common.GetPath("text")+fileName+"\">")
}

// 上传文章
func UpLoadText(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	var col TextData
	if err = bindJSON(c, &col); err != nil {
		return
	}
	if "" == col.Content {
		common.HandleError(c, http.StatusBadRequest, "帖子内容不可为空")
		return
	}
	col.Content = utils.FilterHTML(col.Content)
	rows := models.AddEssay(uid, col.Type, col.Title, col.Content, col.Label)
	if rows != 1 {
		common.HandleError(c, http.StatusNotFound, "上传失败")
		return
	}
	common.Success(c, "上传成功", "")
}

// 点赞
func LikeText(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	textid := c.PostForm("textid")
	if "" == textid {
		return
	}
	tid, err := strconv.Atoi(textid)
	if err != nil {
		return
	}
	if models.IsLiked(uid, tid) {
		common.HandleError(c, http.StatusBadRequest, "您已经为这篇文章点赞")
		return
	}
	common.IncrementField(common.TableText, common.TextFieldLike, tid)
	err = models.AddLike(uid, tid)
	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "点赞失败")
		return
	}
	common.Success(c, "点赞成功", "")
}

// 取消点赞
func UnLikeText(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	textid := c.PostForm("textid")
	if "" == textid {
		return
	}
	tid, err := strconv.Atoi(textid)
	if err != nil {
		return
	}
	if !models.IsLiked(uid, tid) {
		common.HandleError(c, http.StatusBadRequest, "您已取消点赞")
		return
	}

	if models.DelLike(uid, tid) != nil {
		common.HandleError(c, http.StatusBadRequest, "取消点赞失败")
		return
	}
	common.DecrementField(common.TableText, common.TextFieldLike, tid)
	common.Success(c, "取消点赞成功", "")
}

func GetEssayAll(c *gin.Context) {
	GetAllText(c, 0)
}

func GetPostAll(c *gin.Context) {
	GetAllText(c, 1)
}

func GetAllText(c *gin.Context, textType int) {
	page, pageSize, err := common.GetPageAndPageSize(c)
	if err != nil {
		return
	}
	text, err := models.GetAllText(page, pageSize, textType)
	if err != nil {
		common.HandleError(c, http.StatusNotFound, "获取文章列表失败")
		return
	}
	common.Success(c, "获取所有帖子成功", text)
}

// 查看我的文章
func GetMyEssay(c *gin.Context) {
	uid, _ := GetUidFromContext(c)
	var t Page
	c.BindJSON(&t)
	time := strings.Split(t.Time, "-")
	essay := models.GetEssayByUid(uid, time[0], time[1])
	common.Success(c, "", essay)
}

func CheckText(c *gin.Context) {
	textid := c.PostForm("textid")
	tid, _ := strconv.Atoi(textid)
	var tc common.TextContainer
	var text models.Text
	if !models.GetTextById(tid, &text) {
		common.HandleError(c, http.StatusBadRequest, "您查看的文章可能被删除")
		return
	}
	u, ok := c.Get("user")
	if ok {
		uid, _ := u.(int)
		tc.Collect = models.IsCollected(uid, tid)
		tc.Liked = models.IsLiked(uid, tid)
	}
	common.IncrementField(common.TableText, common.TextFieldView, tid)
	tc.Text = text
	common.Success(c, "", tc)
}

func DelText(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	textid := c.PostForm("textid")
	if "" == textid {
		common.HandleError(c, http.StatusBadRequest, "未接收到文章id")
		return
	}
	tid, err := strconv.Atoi(textid)
	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "文章id有误")
		return
	}

	if uid != models.GetAuthorById(tid) {
		common.HandleError(c, http.StatusForbidden, "无法删除文章")
		return
	}
	if !models.DelTextById(tid) {
		common.HandleError(c, http.StatusBadRequest, "文章删除失败")
		return
	}
	common.Success(c, "删除成功", "")
}

func KeySearch(c *gin.Context) {
	s := c.Query("keywords")
	if "" == s {
		common.HandleError(c, http.StatusBadRequest, "请输入搜索内容")
		return
	}
	keys := strings.Split(utils.FilterHTML(s), " ")
	if keys == nil {
		return
	}
	var res []models.Search
	for _, key := range keys {
		res, _ = models.LookInTitle(key)
		sort.Sort(models.ByWeight(res))
	}
	if len(res) == 0 {
		common.Success(c, "未找到您想要的结果", "")
		return
	}
	common.Success(c, "", res)
}
