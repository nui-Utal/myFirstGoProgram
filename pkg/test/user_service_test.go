package test

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"xinyeOfficalWebsite/pkg/common"
	. "xinyeOfficalWebsite/pkg/handlers"
)

var key string = "welcome_to_xinye_official_website_kkNqWgR3ULX8ajvg4hY7Z25zMeBdykjn"

// 从token中获取id
func TestUserIDFromToken(t *testing.T) {
	var userID int
	Convey("Given a valid token", t, func() {
		claims := common.Claims{
			Userid: 10,
			StandardClaims: jwt.StandardClaims{
				// 过期时间
				ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
				// 指定token发行人
				Issuer: "xinye_2022",
			},
		}
		token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))

		router := gin.Default()
		router.POST("/user/head", func(c *gin.Context) {
			// 从请求头中提取令牌
			token := c.Request.Header.Get("Authorization")
			// 解析令牌，获取用户ID
			userID = extractUserIDFromToken(token)
		})

		Convey("When calling the API with the token", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/user/head", nil)
			// 添加测试令牌到请求头
			req.Header.Set("Authorization", token)

			router.ServeHTTP(w, req)

			Convey("The user ID should be extracted correctly", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				So(userID, ShouldEqual, 10)
			})
		})
	})
}

func extractUserIDFromToken(token string) int {
	claims, _ := common.ParseToken(token, "user")
	return claims.Userid
}

func TestUploadAvator(t *testing.T) {
	// 创建一个路由引擎
	router := gin.Default()

	// 使用自定义中间件进行登录认证
	router.Use(func(c *gin.Context) {
		// 假设这里进行登录认证，并将用户ID设置到上下文中
		c.Set("user", 10)
		c.Next()
	})

	// 定义需要测试的路由
	router.POST("/your-endpoint", UploadAvatar)

	// 创建虚拟请求和响应对象
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/your-endpoint", nil)
	req.Header.Set("Authorization", "Bearer your-token")

	// 调用处理函数
	router.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expectedResponseBody := `{"message":"Success","user_id":10}`
	if w.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body %s, but got %s", expectedResponseBody, w.Body.String())
	}
}
