package main

import (
	"wechat/models"
	"wechat/router"
)

// @title        gin+gorm 测试swagger
// @version      v1.0
// @description  gin+gorm 后台API
// @license.name Apache 2.0
// @contact.name go-swagger文档
// @contact.url  https://github.com/swaggo/swag/blob/master/README_zh-CN.md
// @host         222.196.37.70:8080
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @BasePath     /
func main() {
	models.InitMySQL()
	router.InitRouter()
}
