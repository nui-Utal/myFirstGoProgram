package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"xinyeOfficalWebsite/pkg/common"
	"xinyeOfficalWebsite/pkg/models"
	"xinyeOfficalWebsite/pkg/utils"
)

func bindJson(c *gin.Context, cl *models.Carousel) (err error) {
	if err = c.BindJSON(&cl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	return nil
}

// 得到用户列表
func GetUserList(c *gin.Context) {
	page, pageSize, err := common.GetPageAndPageSize(c)
	if err != nil {
		return
	}
	users, err := models.GetUsers(page, pageSize)
	if err != nil {
		common.HandleError(c, http.StatusOK, "查询失败")
		return
	}

	us := common.BatchUserConversion(users)

	common.Success(c, "查询成功", us)
}

func UploadCarousel(c *gin.Context) {
	var cl models.Carousel
	index := c.Request.FormValue("index")
	filename, _ := utils.GenerateSalt(16)
	filename += ".png"
	if err := common.UploadPicture(c, "carousel", filename); err != nil {
		return
	}
	cl.URL = common.GetPath("carousel") + filename
	cl.Order, _ = strconv.Atoi(index)
	cl.Name = filename
	cl.UploadTime = utils.GetLocalTime()
	if models.AddCarousel(cl) {
		common.HandleError(c, http.StatusBadRequest, "图片上传失败")
		return
	}
	common.Success(c, "图片上传成功", "")
}

func GetCarousel(c *gin.Context) {
	page, _, err := common.GetPageAndPageSize(c)
	if err != nil {
		return
	}
	cl, _ := models.GetCarousels(page, 15)
	common.Success(c, "", cl)
}

func DelCarousel(c *gin.Context) {
	image := c.PostForm("name")
	if "" == image {
		common.HandleError(c, http.StatusBadRequest, "未接收到轮播图图片文件名")
		return
	}
	// 从服务器中删除
	err := os.Remove(common.GetPath("carousel") + image)
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "删除失败")
		return
	}
	// 从数据库中删除
	if models.DeleteCarousel(image) {
		common.HandleError(c, http.StatusBadRequest, "未找到您要删除的图片")
		return
	}
	common.Success(c, "删除成功！", "")
}

func UpdateOrder(c *gin.Context) {
	var cl models.Carousel
	if err := bindJson(c, &cl); err != nil {
		return
	}
	if models.UpdateOrderByName(cl.Name, cl.Order) {
		common.HandleError(c, http.StatusBadRequest, "请重新确认您输入的图片名称")
		return
	}
	common.Success(c, "修改成功", "")
}

func GetUserAccount(c *gin.Context) {
	var u userData
	if err := c.BindJSON(&u); err != nil {
		common.HandleError(c, http.StatusBadRequest, "")
		return
	}
	user := models.GetUserByName(u.Username)
	common.Success(c, "获取成功", user)
}

type userInfo struct {
	Page     int    `json:"page"`
	PageSize int    `json:"limit"`
	Time     string `json:"time"`
	UserName string `json:"username"`
}

func GetAllComment(c *gin.Context) {
	var u userInfo
	if err := c.BindJSON(&u); err != nil {
		common.HandleError(c, http.StatusBadRequest, "")
		return
	}
	time := strings.Split(u.Time, "-")
	if "" == u.UserName {
		common.Success(c, "", models.GetCommentLimitTime(time[0], time[1]))
	} else {
		id := models.GetUserByName(u.UserName).ID
		comments := models.GetCommentByUidLimitTime(id, time[0], time[1])
		common.Success(c, "", comments)
	}
}

func AdminLogin(c *gin.Context) {
	var a models.Admin
	if err := c.BindJSON(&a); err != nil {
		common.HandleError(c, http.StatusBadRequest, "无法得到数据")
		return
	}
	admin := models.GetAdminByName(a.Name)
	if admin.Password != utils.Md5Encrypt(a.Password+admin.Salt) {
		common.HandleError(c, http.StatusUnauthorized, "密码错误，请重新输入")
		return
	}
	common.AddToken(c, admin.ID, "admin")
	common.Success(c, "登陆成功", "")
}

func AdminLogout(c *gin.Context) {
	if err := Logout(c, "admin"); err != nil {
		return
	}
	// 重定向到登录页面或其他任意页面
	c.Redirect(http.StatusSeeOther, "/login/backend")
}

func ChangeAdminPwd(c *gin.Context) {
	ori := c.PostForm("ori")
	newPwd := c.PostForm("new")
	if "" == ori || "" == newPwd {
		common.HandleError(c, http.StatusBadRequest, "未接收到完整参数")
		return
	}
	aid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	admin := models.GetAdminById(aid)
	if admin.Password != utils.Md5Encrypt(ori+admin.Salt) {
		common.HandleError(c, http.StatusBadRequest, "您输入的原密码有误")
		return
	}
	if !models.UpdateAdminPwdById(admin.ID, utils.Md5Encrypt(newPwd+admin.Salt)) {
		common.HandleError(c, http.StatusInternalServerError, "更新失败")
		return
	}
	common.Success(c, "密码修改成功", "")

}
func AdminDelComment(c *gin.Context) {
	comid := c.Query("cid")
	cid, err := strconv.Atoi(comid)
	if err != nil {
		common.HandleError(c, http.StatusNotFound, "未接受到评论id")
		return
	}
	if models.DelCommentById(cid) {
		common.Success(c, "评论删除成功", "")
	}
}

func AdminDelText(c *gin.Context) {
	textid := c.Query("textid")
	tid, err := strconv.Atoi(textid)
	if err != nil {
		common.HandleError(c, http.StatusNotFound, "未接受到评论id")
		return
	}
	if !models.DelTextById(tid) {
		common.HandleError(c, http.StatusBadRequest, "文章删除失败")
		return
	}
	common.Success(c, "评论文章成功", "")

}

func GetUserText(c *gin.Context) {
	var u userInfo
	if err := c.BindJSON(&u); err != nil {
		common.HandleError(c, http.StatusBadRequest, "")
		return
	}
	time := strings.Split(u.Time, "-")
	if "" == u.UserName {
		common.Success(c, "", models.GetAllEssayLimitTime(u.Page, u.PageSize, time[0], time[1]))
	} else {
		id := models.GetUserByName(u.UserName).ID
		essays := models.GetEssayByUid(id, time[0], time[1])

		common.Success(c, "", essays)
	}
}

func UpdateAdminInfo(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		common.HandleError(c, http.StatusUnauthorized, "您的登录状态有误，请重新登录")
		return
	}
	uname := c.Request.FormValue("username")

	tx := utils.DB.Begin()
	if !models.UpdateNameById(uid, uname) {
		common.HandleError(c, http.StatusBadRequest, "名称更新失败")
		tx.Rollback()
		return
	}

	filename, _ := utils.GenerateSalt(16)
	filename += ".png"
	if err := common.UploadPicture(c, "head", filename); err != nil {
		tx.Rollback()
		return
	}
	if !models.UpdateAvatorById(uid, filename) {
		common.HandleError(c, http.StatusBadRequest, "头像更新失败")
		tx.Rollback()
		return
	}
	tx.Commit()
	common.Success(c, "更新成功", common.GetPath("head")+filename)
}
