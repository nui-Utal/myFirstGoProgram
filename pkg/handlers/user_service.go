package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/wenlng/go-captcha/captcha"
	_ "image/png"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"xinyeOfficalWebsite/pkg/common"
	"xinyeOfficalWebsite/pkg/models"
	"xinyeOfficalWebsite/pkg/utils"
)

var ctx = context.Background()

var capt = captcha.GetCaptcha()

type captchaData struct {
	Dots string `json:"dots"`
	Key  string `json:"key"`
}

type userData struct {
	UserId   int    `json:"userid"`
	Username string `json:"username,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Sno      string `json:"sno,omitempty"`
	Special  string `json:"special,omitempty"`
	Pwd      string `json:"password,omitempty" mysql:"default:123456"`
	Remember bool   `json:"remember_me,omitempty"`
}

func GetUidFromContext(c *gin.Context) (i int, err error) {
	u, ok := c.Get("user")
	i, ok1 := u.(int)
	if !ok || !ok1 {
		common.HandleError(c, http.StatusInternalServerError, "请登录后尝试")
	}
	return
}

// 验证手机号
func PhoneNumVerification(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	r, _ := regexp.Compile(pattern)
	return r.MatchString(phone)
}

// 用户注册
func UserRegistration(c *gin.Context) {
	var u userData
	// 获取数据
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if "" != u.Phone {
		if !PhoneNumVerification(u.Phone) {
			common.HandleError(c, http.StatusOK, "请使用有效的手机号码")
			return
		}
	}
	if u.Remember {
		common.AddToken(c, u.UserId, "user")
	}
	err, user := models.CreateUser(u.Username, u.Phone, u.Sno, u.Special, u.Pwd)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			common.HandleError(c, http.StatusBadRequest, "该用户名或学号已被注册过")
			return
		}
	}
	common.Success(c, "注册成功", common.UserConversion(user))
}

func Logout(c *gin.Context, end string) (er error) {
	token, err := c.Cookie("remember_me_token")
	if err != nil {
		common.HandleError(c, http.StatusUnauthorized, "您未处于登录状态")
		return
	}
	// 解析并验证 token
	switch end {
	case "user":
		_, err = common.ParseToken(token, "user")
	case "admin":
		_, err = common.ParseToken(token, "admin")
	}
	// 删除 token 对应的 Cookie
	c.SetCookie("remember_me_token", "", -1, "/", "", false, true)

	common.Success(c, "成功退出登录！", "")
	return nil
}

func UserLogout(c *gin.Context) {
	if err := Logout(c, "user"); err != nil {
		return
	}

	// 重定向到登录页面或其他任意页面
	c.Redirect(http.StatusSeeOther, "/user/login")
}

// 发送图片验证相关数据
func GetCaptchaData(c *gin.Context) {
	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		common.HandleError(c, http.StatusInternalServerError, "无法获取验证图片")
		return
	}
	writeCache(dots, key)
	common.Success(c, "", map[string]interface{}{
		"image_base64": b64,
		"thumb_base64": tb64,
		"captcha_key":  key,
	})
}

// 图片验证
func PictureVerification(c *gin.Context) {
	var data captchaData
	if err := c.BindJSON(&data); err != nil {
		common.HandleError(c, http.StatusBadRequest, "未接收到数据")
		return
	}
	if data.Dots == "" || data.Key == "" {
		common.HandleError(c, http.StatusBadRequest, "请点击后确认")
		return
	}

	cacheData := readCache(data.Key)
	if cacheData == "" {
		common.HandleError(c, http.StatusNotFound, "请刷新页面后重试")
		return
	}
	src := strings.Split(data.Dots, ",")

	var dct map[int]captcha.CharDot
	if err := json.Unmarshal([]byte(cacheData), &dct); err != nil {
		common.HandleError(c, http.StatusNotFound, "请刷新页面后重试")
		return
	}

	chkRet := false
	if (len(dct) * 2) == len(src) {
		for i, dot := range dct {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[k]), 64)

			chkRet = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5)
			if !chkRet {
				break
			}
		}
	}

	if chkRet {
		common.Success(c, "通过验证", "")
		return
	}
}

// 将数据写入redis缓存
func writeCache(v interface{}, key string) {
	bt, _ := json.Marshal(v)
	expiration := 7 * 24 * time.Hour
	// set(context.background, key, value, expiration
	err := utils.Rdb.Set(ctx, key, bt, expiration).Err()
	if err != nil {
		panic(err)
	}
}

// 从Redis缓存中读取数据
func readCache(key string) string {
	val, err := utils.Rdb.Get(ctx, key).Result()
	if err != nil {
		// 处理缓存不存在的情况
		if err == redis.Nil {
			return ""
		}
		panic(err)
	}
	return val
}

// 用户登录
func Login(c *gin.Context) {
	var u userData
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 得到salt和加密后的密码
	var salt, pwd string
	if "" == u.Username {
		salt, pwd, _ = models.GetSaltAndPwd("phone", u.Phone)
	} else {
		salt, pwd, _ = models.GetSaltAndPwd("username", u.Username)
	}

	curPwd := utils.Md5Encrypt(u.Pwd + salt)
	// 判断当前密码加密后是否和查询得到的密码相同
	if curPwd != pwd {
		common.HandleError(c, http.StatusBadRequest, "登录失败，请重新输入")
		return
	}

	// 返回对应的user数据
	user := models.LoginWithPassword(pwd)
	if u.Remember {
		common.AddToken(c, user.ID, "user")
	}

	common.Success(c, "登录成功", common.UserConversion(user))
}

// 修改密码
func EditPwd(c *gin.Context) {
	//获取数据
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	password := c.PostForm("password")
	if "" == password {
		common.HandleError(c, http.StatusBadRequest, "未接收到新密码")
		return
	}

	salt := models.GetSaltById(uid)
	pwd := utils.Md5Encrypt(password + salt)
	user := models.UpdatePwdById(uid, pwd)

	common.Success(c, "修改密码成功", common.UserConversion(user))
}

// 修改手机号
func EditPhone(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	phone := c.PostForm("phone")
	if "" == phone {
		common.HandleError(c, http.StatusBadRequest, "未接收到手机号")
		return
	}

	if !PhoneNumVerification(phone) {
		common.HandleError(c, http.StatusOK, "请使用有效的手机号码")
		return
	}

	user := models.UpdatePhoneById(uid, phone)

	common.Success(c, "手机号修改成功", common.UserConversion(user))
}

// 修改学号
func EditSno(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	sno := c.PostForm("sno")
	if "" == sno {
		common.HandleError(c, http.StatusBadRequest, "未接收到学号")
		return
	}

	user := models.UpdateSnoById(uid, sno)

	common.Success(c, "修改学号成功", common.UserConversion(user))
}

// 搜索用户
func SearchUser(c *gin.Context) {
	// 获取数据
	var u userData
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if "" == u.Sno {
		// 用户名搜索
		users := models.SearchUserByName(u.Username)
		if len(users) == 0 {
			common.HandleError(c, http.StatusBadRequest, "未搜索到有关的用户")
			return
		}
		common.Success(c, "", common.BatchUserConversion(users))
	} else {
		// 学号搜索
		user := models.GetUserBySno(u.Sno)
		if "" == user.Username {
			common.HandleError(c, http.StatusBadRequest, "未搜索到有关的用户")
			return
		}
		common.Success(c, "", common.UserConversion(user))
	}
}

// 关注
func Follow(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	followed := c.PostForm("userid")
	if "" == followed {
		common.HandleError(c, http.StatusBadRequest, "未接收到关注人信息")
		return
	}
	f, _ := strconv.Atoi(followed)
	models.AddFollow(f, uid)

	if err = common.IncrementField(common.TableUsers, common.UserFieldFollowNum, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = common.IncrementField(common.TableUsers, common.UserFieldFanNum, f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	common.Success(c, "关注成功", "")
}

// 取关
func UnFollow(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	followed := c.PostForm("userid")
	if "" == followed {
		common.HandleError(c, http.StatusBadRequest, "未接收到关注人信息")
		return
	}
	f, _ := strconv.Atoi(followed)

	models.DelFollow(f, uid)
	if err := common.DecrementField(common.TableUsers, common.UserFieldFollowNum, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // f.Fan
		return
	}
	if err := common.DecrementField(common.TableUsers, common.UserFieldFanNum, f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	common.Success(c, "取消关注成功", "")
}

// 上传头像
func UploadAvatar(c *gin.Context) {
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}
	fileName, _ := utils.GenerateSalt(16)
	fileName += ".png"
	err = common.UploadPicture(c, "head", fileName)
	if err != nil {
		return
	}

	// 删除原头像
	ori := models.GetAvatar(uid)
	if "" != ori {
		imagePath := common.GetPath("head") + ori // 图片文件路径
		err = os.Remove(imagePath)
		if err != nil {
			// 处理删除文件时发生的错误
			panic(err)
		}
	}
	// 更新数据库表
	models.UploadAvatar(fileName, uid)
	avator := common.GetPath("head") + fileName

	common.Success(c, "头像上传成功", avator)
}

// 查看粉丝
func CheckMyFans(c *gin.Context) {
	var p Page
	if err := c.BindJSON(&p); err != nil {
		common.HandleError(c, http.StatusBadRequest, "")
	}
	// 获取数据
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}

	users := models.ShowFans(uid, p.Page, p.PageSize)
	if len(users) != 0 {
		common.Success(c, "", common.BatchUserConversion(users))
	}
}

// 查看关注的人
func CheckMyFollow(c *gin.Context) {
	var p Page
	if err := c.BindJSON(&p); err != nil {
		common.HandleError(c, http.StatusBadRequest, "")
	}
	// 获取数据
	uid, err := GetUidFromContext(c)
	if err != nil {
		return
	}

	users := models.ShowFollow(uid, p.Page, p.PageSize)
	if len(users) != 0 {
		common.Success(c, "", common.BatchUserConversion(users))
	} else {
		common.Success(c, "您还未进行关注", "")
	}
}

func GetUserHomepage(c *gin.Context) {
	uid := c.Param("uid")
	id, err := strconv.Atoi(uid)
	if err != nil {
		common.HandleError(c, http.StatusBadRequest, "您输入的用户id有误")
		return
	}

	user := models.GetUserById(id)

	curid, ok := c.Get("user")
	curID, ok1 := curid.(int)
	if !ok || !ok1 {
		err = fmt.Errorf("未登录")
	}

	// 未登录状态
	if err != nil {
		common.Success(c, "", map[string]interface{}{
			"followed": -1,
			"userinfo": common.UserConversion(user),
		})
		return
	}

	// 登录
	if curID == id {
		common.Success(c, "", map[string]interface{}{
			"followed": 0,
			"userinfo": common.UserConversion(user),
		})
		return
	}
	common.Success(c, "", map[string]interface{}{
		"followed": models.IsFollowed(id, curID),
		"userinfo": common.UserConversion(user),
	})

}
