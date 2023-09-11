package service

import (
	"fmt"
	"net/http"
	"strconv"
	"wechat/models"
	"wechat/utils/re"
	"wechat/utils/secure"
	"wechat/utils/validator"

	"github.com/gin-gonic/gin"
)

type userParam struct {
	Username    string `json:"username" validate:"required,min=2,max=12"`
	OldPassword string `json:"oldpassword" validate:"min=6,max=20"`
	Password    string `json:"password" validate:"required,min=6,max=20"`
	RePassword  string `json:"repassword" validate:"required,min=6,max=20"`
	Phone       string `json:"phone" validate:"required,len=11"`
	Roleid      int    `json:"roleId"`
}

// GetUser
// @Summary 查询用户
// @Tags 用户模块
// @param uid query string false "用户id"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/test [get]
func ApiTest(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("uid"))
	fmt.Println(uid)
	data, code := models.FindByUid(uid)
	role := c.GetInt("role")
	fmt.Println(role)
	c.JSON(200, gin.H{
		"status":  code,
		"data":    data,
		"total":   1,
		"message": re.GetErrMsg(code),
	})
}

// AddUser
// @Summary 新增用户
// @Tags 用户模块
// @param user body userParam true "用户"
// @Success 200 {string} json{"code","message"}
// @Router /api/addUser [post]
func AddUser(c *gin.Context) {
	fmt.Println("okoko")
	var data userParam
	_ = c.ShouldBindJSON(&data)
	msg, validCode := validator.Validate(&data)
	if validCode != re.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		return
	}
	user := models.FindByPhone(data.Phone)
	if user.Phone != "" {
		c.JSON(400, gin.H{
			"message": "该手机号已被注册！",
		})
		return
	}
	password := data.Password
	repassword := data.RePassword
	if password != repassword {
		c.JSON(400, gin.H{
			"message": "两次密码不一致!",
		})
		return
	}
	user.Username = data.Username
	user.Phone = data.Phone
	user.Salt = secure.RandomStr(10)
	user.Password = secure.MakePassword(password, user.Salt)
	code := models.AddUser(&user)
	c.JSON(200, gin.H{
		"status":  code,
		"message": re.GetErrMsg(code),
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param uid query int true "用户id"
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/deleteUser [get]
func DeleteUser(c *gin.Context) {
	uid := c.GetInt("uid")
	user, _ := models.FindByUid(uid)
	if user.Phone == "" {
		c.JSON(400, gin.H{
			"message": "该用户不存在!",
		})
		return
	}
	code := models.DeleteUser(uid)
	c.JSON(200, gin.H{
		"status":  code,
		"message": re.GetErrMsg(code),
	})
}

// SearchUsers
// @Summary 查询所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/getAllUser [get]
func GetAllUser(c *gin.Context) {
	var code int
	role := c.GetInt("role")
	fmt.Println(role)
	if role != 1 {
		code = re.ERROR_USER_ROLE
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	}
	userList := models.GetUserList()
	code = re.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":   code,
			"message":  re.GetErrMsg(code),
			"userList": userList,
		},
	)
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param user body userParam true "用户"
// @Success 200 {string} json{"code","message"}
// @Security ApiKeyAuth
// @Router /api/updateUser [post]
func UpdateUser(c *gin.Context) {
	var data userParam
	_ = c.ShouldBindJSON(&data)
	msg, code := validator.Validate(&data)
	if code != re.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": msg,
			},
		)
		return
	}
	//用户校验
	var userdata models.User
	userPhone := models.FindByPhone(data.Phone)
	userdata, _ = models.FindByUid(c.GetInt("uid"))
	password := secure.MakePassword(data.OldPassword, userdata.Salt)
	if userPhone.Phone != "" && userPhone.Phone != data.Phone {
		code = re.ERROR_USER_PHONE
	} else if data.Password != data.RePassword {
		code = re.ERROR_USER_PASSWORD_RE
	} else if userdata.Phone == "" {
		code = re.ERROR_USER_EIXT
	} else if password != userdata.Password {
		code = re.ERROR_USER_PASSWORD
	}
	if code != re.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			},
		)
		return
	}
	//更新数据
	userdata.Username = data.Username
	userdata.Phone = data.Phone
	userdata.Password = password
	code = models.UpdateUser(&userdata)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		},
	)
}
