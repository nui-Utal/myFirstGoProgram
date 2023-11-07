package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const user_key = "welcome_to_xinye_official_website_kkNqWgR3ULX8ajvg4hY7Z25zMeBdykjn"
const admin_key = "thank_you_for_the_construction_of_xinye_official_website_kkNqWgR3ULX8ajvg4hY7Z25zMeBdykjn"

type Claims struct {
	Userid int
	jwt.StandardClaims
}

func AddToken(c *gin.Context, userid int, end string) {
	var (
		tokenClaims string
		err         error
	)

	// 签署 JWT
	expirationTime := time.Now().Add(time.Hour * 24 * 7).Unix()

	claims := Claims{
		Userid: userid,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expirationTime,
			// 指定token发行人
			Issuer: "xinye_2022",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	switch end {
	case "user":
		tokenClaims, err = token.SignedString([]byte(user_key))
	case "admin":
		tokenClaims, err = token.SignedString([]byte(admin_key))
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 将加密后的令牌发送给用户的浏览器作为一个长期的 Cookie
	c.SetCookie("remember_me_token", "Bearer "+tokenClaims, int(expirationTime), "/", "", false, true)
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token, end string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	t := strings.Split(token, " ")
	tokenClaims, err := jwt.ParseWithClaims(t[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		var signingKey []byte
		switch end {
		case "user":
			signingKey = []byte(user_key)
		case "admin":
			signingKey = []byte(admin_key)
		}
		return signingKey, nil
	})
	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
