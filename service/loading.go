package service

import (
	"fmt"
	"net/http"
	"wechat/models"
	"wechat/utils/middleware"
	"wechat/utils/re"
	"wechat/utils/secure"

	"github.com/gin-gonic/gin"
)

type userLogin struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

// Login
// @Summary 用户登录
// @Tags 用户模块
// @param user body userLogin true "用户"
// @Success 200 {string} json{"code","message"}
// @Router /api/login [post]
func Login(c *gin.Context) {
	var data userLogin
	var token string
	var err error
	_ = c.ShouldBindJSON(&data)
	fmt.Println(data)
	userdata, code := CheckUser(data)
	fmt.Println(userdata)
	if code == re.SUCCSE {
		token, err = middleware.GenerateToken(int(userdata.ID), userdata.Roleid)
		if err != nil {
			fmt.Println(err)
			code = re.ERROR
		}
		c.JSON(200, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
			"data":    userdata,
			"token":   token,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    userdata,
			"message": re.GetErrMsg(code),
			"token":   token,
		})
	}
}

func CheckUser(data userLogin) (models.User, int) {
	userdata := models.FindByPhone(data.Phone)
	if userdata.Phone == "" {
		return userdata, re.ERROR_USER_EIXT
	}
	password := secure.MakePassword(data.Password, userdata.Salt)
	if password != userdata.Password {
		return userdata, re.ERROR_USER_PASSWORD
	}
	return userdata, re.SUCCSE
}
