package service

import (
	"net/http"
	"strconv"
	"wechat/models"
	"wechat/utils/re"
	"wechat/utils/validator"

	"github.com/gin-gonic/gin"
)

type ApplicationF struct {
	ApplicantId   int    `json:"applicantid" validate:"required"`
	ApplicantName string `json:"applicantname" validate:"required"`
	ReviewerId    int    `json:"reviewerid" validate:"required"`
	ReviewerName  string `json:"reviewername" validate:"required"`
	// CreatedAt       time.Time
	// Status          int //0,未审核；1,通过；2,拒绝
	// ApplicationType int //1,好友申请；2,群聊申请
}

// AddFriendAppli
// @Summary 新增好友申请
// @Tags 好友/群聊申请模块
// @param ApplicationF body ApplicationF true "申请表单"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/addFriendAppli [post]
func AddFriendAppli(c *gin.Context) {
	// var data models.FriendApplication
	var data ApplicationF
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
	friendship := models.GetByUidFid(data.ApplicantId, data.ReviewerId)
	if friendship.Notes != "" {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  re.ERROR_DATA_EIXT,
				"message": "该好友已经存在",
			},
		)
		return
	}

	// data.ApplicationType = 1
	firendAppli := models.FriendApplication{
		ApplicantId:     data.ApplicantId,
		ApplicantName:   data.ApplicantName,
		ReviewerId:      data.ReviewerId,
		ReviewerName:    data.ReviewerName,
		Status:          0,
		ApplicationType: 1,
	}
	code = models.AddAppli(&firendAppli)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		},
	)
}

// GetAppliMsg
// @Summary 申请人查询申请信息
// @Tags 好友/群聊申请模块
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/getAppliMsg [get]
func GetAppliMsg(c *gin.Context) {
	uid := c.GetInt("uid")
	applications := models.GetApplicationList(uid)
	c.JSON(
		http.StatusOK, gin.H{
			"applicationList": applications,
		},
	)
}

// UserGetAppliMsg
// @Summary 好友接受者查询申请信息
// @Tags 好友/群聊申请模块
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/userGetAppliMsg [get]
func UserGetAppliMsg(c *gin.Context) {
	uid := c.GetInt("uid")
	applications := models.GetApplicationListByUser(uid)
	c.JSON(
		http.StatusOK, gin.H{
			"applicationList": applications,
		},
	)
}

// DealFriendAppli
// @Summary 好友处理申请
// @Tags 好友/群聊申请模块
// @Param aid query int true "申请Id"
// @Param flag query int true "是否同意"
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/dealFriendAppli [get]
func DealFriendAppli(c *gin.Context) {
	var code int
	// var msg string
	aid, _ := strconv.Atoi(c.Query("aid"))
	isAc, _ := strconv.Atoi(c.Query("flag"))
	uid := c.GetInt("uid")
	application := models.GetAppliById(aid)
	if application.ApplicantName == "" {
		code = re.ERROR
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	}
	if application.ReviewerId != uid {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  2001,
				"message": "用户没有权限,本人才能通过审核",
			},
		)
		return
	}
	if application.Status != 0 {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  2000,
				"message": "该申请已处理",
			},
		)
		return
	}
	friendship := models.GetByUidFid(application.ApplicantId, application.ReviewerId)
	if friendship.Notes != "" {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  re.ERROR_DATA_EIXT,
				"message": "该好友已经存在",
			},
		)
		return
	}
	if isAc == 0 {
		application.Status = 2
		code = models.UpdateAppli(&application)
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	} else {
		firendship := models.Firendship{
			Uid:   application.ApplicantId,
			Fid:   application.ReviewerId,
			Notes: application.ReviewerName,
		}
		code = models.Createship(&firendship)
		if code != re.SUCCSE {
			c.JSON(
				http.StatusOK, gin.H{
					"status":  code,
					"message": "好友添加失败，请重试",
				},
			)
			return
		}
		firendship1 := models.Firendship{
			Uid:   application.ReviewerId,
			Fid:   application.ApplicantId,
			Notes: application.ApplicantName,
		}
		code = models.Createship(&firendship1)
		if code != re.SUCCSE {
			c.JSON(
				http.StatusOK, gin.H{
					"status":  code,
					"message": "好友添加失败，请重试",
				},
			)
			return
		}

		application.Status = 1
		code = models.UpdateAppli(&application)
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	}
	code = re.ERROR_PARAM
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		},
	)
}

// AddGroupAppli
// @Summary 新增群聊申请
// @Tags 好友/群聊申请模块
// @param ApplicationF body ApplicationF true "申请表单"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/addGroupAppli [post]
func AddGroupAppli(c *gin.Context) {
	var data ApplicationF
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
	groupship := models.GetGMemberByUidGid(data.ReviewerId, data.ApplicantId)
	if groupship.Groupname != "" {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  re.ERROR_DATA_EIXT,
				"message": "你已经是该群成员",
			},
		)
		return
	}
	// data.ApplicationType = 1
	firendAppli := models.FriendApplication{
		ApplicantId:     data.ApplicantId,
		ApplicantName:   data.ApplicantName,
		ReviewerId:      data.ReviewerId,
		ReviewerName:    data.ReviewerName,
		Status:          0,
		ApplicationType: 2,
	}
	code = models.AddAppli(&firendAppli)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		},
	)
}

// DealGroupAppli
// @Summary 群主处理申请
// @Tags 好友/群聊申请模块
// @Param aid query int true "申请Id"
// @Param flag query int true "是否同意"
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/dealGroupAppli [get]
func DealGroupAppli(c *gin.Context) {
	var code int
	// var msg string
	aid, _ := strconv.Atoi(c.Query("aid"))
	isAc, _ := strconv.Atoi(c.Query("flag"))
	uid := c.GetInt("uid")
	application := models.GetAppliById(aid)
	if application.ApplicantName == "" {
		code = re.ERROR
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	}
	group := models.FindByGid(application.ReviewerId)
	if group.Uid != uid {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  2001,
				"message": "用户没有权限,群主才能通过审核",
			},
		)
		return
	}
	if application.Status != 0 {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  2000,
				"message": "该申请已处理",
			},
		)
		return
	}
	if isAc == 0 {
		application.Status = 2
		code = models.UpdateAppli(&application)
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	} else if isAc == 1 {
		groupship := models.Groupship{
			Gid:       application.ReviewerId,
			Groupname: application.ReviewerName,
			Uid:       application.ApplicantId,
			Username:  application.ApplicantName,
			Roleid:    1,
		}
		code = models.CreateGroupShip(&groupship)
		if code != re.SUCCSE {
			c.JSON(
				http.StatusOK, gin.H{
					"status":  code,
					"message": re.GetErrMsg(code),
				},
			)
			return
		}
		application.Status = 1
		code = models.UpdateAppli(&application)
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	}
	code = re.ERROR_PARAM
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		},
	)
}
