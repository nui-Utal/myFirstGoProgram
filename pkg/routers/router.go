package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinyeOfficalWebsite/pkg/handlers"
	"xinyeOfficalWebsite/pkg/middlewares"
)

func LoadRoutes(r *gin.Engine) {

	//设置没有路由时访问的页面（以免出现404）
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "该页面不存在，请返回首页"})
		//c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	visit := r.Group("/", middlewares.CheckNotAsk)
	{
		visit.POST("/post/look", handlers.CheckText)
		visit.GET("/user/:uid", handlers.GetUserHomepage)
	}

	r.POST("/user/register", handlers.UserRegistration)
	r.POST("/user/captchaChack", handlers.PictureVerification)
	r.POST("/user/login", handlers.Login)
	r.POST("/user/search", handlers.SearchUser)

	r.GET("/text/search", handlers.KeySearch)
	r.GET("/post/getcomment", handlers.GetTextComment)
	r.GET("/user/logout", handlers.UserLogout)
	r.GET("/user/getCaptchaData", handlers.GetCaptchaData)
	r.GET("/post/getall", handlers.GetPostAll)
	r.GET("/essay/getall", handlers.GetEssayAll)
	//r.GET("/user")

	// 登录检查
	auth := r.Group("/", middlewares.CheckUserLogin)
	{
		auth.POST("/user/head", handlers.UploadAvatar)
		auth.POST("/user/pwd", handlers.EditPwd)
		auth.POST("/user/phone", handlers.EditPhone)
		auth.POST("/user/sno", handlers.EditSno)
		auth.POST("/user/follow", handlers.Follow)
		auth.POST("/user/unfollow", handlers.UnFollow)
		auth.POST("/user/checkmyfans", handlers.CheckMyFans)
		auth.POST("/user/checkmyfollow", handlers.CheckMyFollow)
		auth.POST("/user/comment", handlers.CheckMyComment)
		//auth.POST("/user/article", handlers.GetMyEssay)

		auth.POST("/essay/collect", handlers.CollectText)
		auth.POST("/essay/uploadPicture", handlers.UploadPicture)
		auth.POST("/essay/upload", handlers.UpLoadText)
		auth.POST("/essay/like", handlers.LikeText)
		auth.POST("/essay/unlike", handlers.UnLikeText)
		auth.POST("/essay/reply", handlers.Reply)
		auth.POST("/essay/delReply", handlers.DelReply)
		auth.POST("/essay/comment", handlers.CommentText)
		auth.POST("/essay/delcom", handlers.DelComment)
		auth.POST("/essay/delEssay", handlers.DelText)

		auth.POST("/post/collect", handlers.CollectText)
		auth.POST("/post/upload", handlers.UpLoadText)
		auth.POST("/post/delpost", handlers.DelText)
		auth.POST("/post/like", handlers.LikeText)
		auth.POST("/post/notlike", handlers.UnLikeText)
		auth.POST("/post/comment", handlers.CommentText)
		auth.POST("/post/delcom", handlers.DelComment)
		auth.GET("/collect/getall", handlers.GetCollectAll)
		auth.GET("/like/getall", handlers.GetLikeAll)
	}

	r.POST("/login/backend", handlers.AdminLogin)
	r.POST("/logout/backend", handlers.AdminLogout)
	adm := r.Group("/", middlewares.CheckAdminLogin)
	{
		adm.GET("/users/getAllUser", handlers.GetUserList)
		adm.POST("/upload/carousel", handlers.UploadCarousel)
		adm.GET("/get/carousel", handlers.GetCarousel)
		adm.POST("/delete/carousel", handlers.DelCarousel)
		adm.POST("/update/carousel", handlers.UpdateOrder)
		adm.POST("/get/useraccount", handlers.GetUserAccount)
		adm.POST("/get/usercomment", handlers.GetAllComment)
		adm.POST("/update/userpwd", handlers.ChangeAdminPwd)
		adm.GET("/delete/usercomment", handlers.AdminDelComment)
		adm.GET("/get/userarticle", handlers.GetUserText)
		adm.GET("/delete/userarticle", handlers.AdminDelText)
		adm.POST("/update/admininfo", handlers.UpdateAdminInfo)
	}
}
