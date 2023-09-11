package router

import (
	"wechat/docs"
	"wechat/service"
	"wechat/utils/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() {
	//设置模式
	// gin.SetMode(utils.AppMode)
	gin.SetMode("debug")
	r := gin.New()

	// r.Use(cors.Default())

	r.Use(middleware.Cors())
	//swagger接口文档
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//后台接口请求
	login := r.Group("/api")
	// login.GET("/test", service.ApiTest)
	login.GET("/getwsconn", service.GetConn)
	login.POST("/login", service.Login)
	login.POST("/addUser", service.AddUser)

	auth := r.Group("/api")
	auth.Use(middleware.TokenVerify())
	{
		//用户接口
		auth.GET("/test", service.ApiTest)
		auth.GET("/deleteUser", service.DeleteUser)
		auth.GET("/getAllUser", service.GetAllUser)
		auth.POST("/updateUser", service.UpdateUser)
		//好友申请接口
		auth.POST("/addFriendAppli", service.AddFriendAppli)
		auth.GET("/getAppliMsg", service.GetAppliMsg)
		auth.GET("/userGetAppliMsg", service.UserGetAppliMsg)
		auth.GET("/dealFriendAppli", service.DealFriendAppli)
		//好友管理接口
		auth.GET("/getfriends", service.GetFirends)
		auth.GET("/deleteFriends", service.DeleteFirends)
		//群聊管理接口
		auth.POST("/createGroup", service.CreateGroup)
		auth.GET("/deleteGroup", service.DeleteGroup)
		auth.POST("/updateGroup", service.UpdateGroup)
		auth.GET("/getGroupByName", service.GetGroupByName)
		auth.GET("/getGroupByGid", service.GetGroupByGid)
		auth.POST("/addGroupAppli", service.AddGroupAppli)
		auth.GET("/dealGroupAppli", service.DealGroupAppli)

		//消息管理接口
		// auth.GET("/getwsconn", service.GetConn)
		auth.GET("/getFhistoryMsg", service.GetFHistroryMsg)
		auth.GET("/getGhistoryMsg", service.GetGHistroryMsg)
		auth.POST("/uploadfile", service.UploadFile)
		auth.GET("/downloadfile", service.DownloadFile)
	}

	r.Run(":8080")
}
