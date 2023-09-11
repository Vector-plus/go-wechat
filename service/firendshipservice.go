package service

import (
	"net/http"
	"strconv"
	"wechat/models"
	"wechat/utils/re"

	"github.com/gin-gonic/gin"
)

// GetFirends
// @Summary 获取好友列表
// @Tags 好友管理模块
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/getfriends [get]
func GetFirends(c *gin.Context) {
	uid := c.GetInt("uid")
	friends := models.GetFirendship(uid)
	c.JSON(
		http.StatusOK, gin.H{
			"status":      re.SUCCSE,
			"firendsList": friends,
		},
	)
}

// DeleteFirends
// @Summary 删除好友
// @Tags 好友管理模块
// @Success 200 {string} json{"code","message"}
// @Param fid query string true "好友id"
// @Security ApiKeyAuth
// @Router /api/deleteFriends [get]
func DeleteFirends(c *gin.Context) {
	uid := c.GetInt("uid")
	fid, _ := strconv.Atoi(c.Query("fid"))
	code := models.DelteFirend(uid, fid)

	c.JSON(
		http.StatusOK, gin.H{
			"status":      code,
			"firendsList": re.GetErrMsg(code),
		},
	)
}
