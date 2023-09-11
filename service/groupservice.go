package service

import (
	"net/http"
	"strconv"
	"wechat/models"
	"wechat/utils/re"
	"wechat/utils/validator"

	"github.com/gin-gonic/gin"
)

type Mgroup struct {
	Gid         int    `json:"gid"`
	Uid         int    `json:"uid" validate:"required"`
	GroupName   string `json:"groupname" validate:"required,min=2,max=12"`
	Description string `json:"description" validate:"required,min=10,max=250"`
}

// CreateGroup
// @Summary 新增群聊
// @Tags 群聊管理模块
// @param Mgroup body Mgroup true "群聊信息"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/createGroup [post]
func CreateGroup(c *gin.Context) {
	var data Mgroup
	_ = c.ShouldBindJSON(&data)
	msg, code := validator.Validate(data)
	if code != re.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": msg,
			},
		)
		return
	}
	msgGroup := models.Msggroup{
		Uid:         data.Uid,
		Groupname:   data.GroupName,
		Description: data.Description,
	}
	code = models.CreateGroup(&msgGroup)
	if code == re.SUCCSE {
		user, _ := models.FindByUid(data.Uid)
		groupship := models.Groupship{
			Gid:       int(msgGroup.ID),
			Groupname: msgGroup.Groupname,
			Uid:       msgGroup.Uid,
			Username:  user.Username,
			Roleid:    0,
		}
		code = models.CreateGroupShip(&groupship)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"code":    code,
			"message": re.GetErrMsg(code),
		},
	)
}

// DeleteGroup
// @Summary 删除群聊
// @Tags 群聊管理模块
// @param gid query string true "群聊id"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/deleteGroup [get]
func DeleteGroup(c *gin.Context) {
	uid := c.GetInt("uid")
	gid, _ := strconv.Atoi(c.Query("gid"))
	msgGooup := models.FindByGid(gid)
	if msgGooup.Groupname == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    re.ERROR,
				"message": "该群不存在",
			},
		)
		return
	}
	if msgGooup.Uid != uid {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    re.ERROR,
				"message": "群主才能删除群聊",
			},
		)
	}
	code := models.DeleteGroup(gid)
	c.JSON(
		http.StatusOK, gin.H{
			"code":    code,
			"message": re.GetErrMsg(code),
		},
	)
}

// UpdateGroup
// @Summary 更新群聊
// @Tags 群聊管理模块
// @param Mgroup body Mgroup true "群聊信息"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/updateGroup [post]
func UpdateGroup(c *gin.Context) {
	var data Mgroup
	_ = c.ShouldBindJSON(&data)
	msg, code := validator.Validate(data)
	if code != re.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": msg,
			},
		)
		return
	}
	uid := c.GetInt("uid")
	gid := data.Gid
	msgGooup := models.FindByGid(gid)
	if msgGooup.Groupname == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    re.ERROR,
				"message": "该群不存在",
			},
		)
		return
	}
	if msgGooup.Uid != uid {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    re.ERROR,
				"message": "群主才能更新群聊",
			},
		)
	}
	msgGooup.Description = data.Description
	msgGooup.Groupname = data.GroupName
	msgGooup.Uid = data.Uid
	code = models.UpdateGroup(&msgGooup)
	c.JSON(
		http.StatusOK, gin.H{
			"code":    code,
			"message": re.GetErrMsg(code),
		},
	)
}

// GetGroupByName
// @Summary 群名模糊查询群聊
// @Tags 群聊管理模块
// @param gname query string true "群名"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/getGroupByName [get]
func GetGroupByName(c *gin.Context) {
	gname := c.Query("gname")
	if gname == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    re.ERROR,
				"message": "群名为空，参数不全",
			},
		)
		return
	}
	msgGoupsList := models.GetGroupsByName(gname)
	c.JSON(
		http.StatusOK, gin.H{
			"code":          re.SUCCSE,
			"msgGroupsList": msgGoupsList,
		},
	)
}

// GetGroupByGid
// @Summary 使用Gid查询群聊
// @Tags 群聊管理模块
// @param gid query string true "群id"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/getGroupByGid [get]
func GetGroupByGid(c *gin.Context) {
	gid := c.Query("gid")
	if gid == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    re.ERROR,
				"message": "群id为空,参数不全",
			},
		)
		return
	}
	id, _ := strconv.Atoi(gid)
	msgGroup := models.FindByGid(id)

	if msgGroup.Groupname == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"code": re.SUCCSE,
				"msg":  "该群不存在",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"code":          re.SUCCSE,
			"msgGroupsList": msgGroup,
		},
	)
}
