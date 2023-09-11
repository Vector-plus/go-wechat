package middleware

import (
	"net/http"
	"wechat/utils/re"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 0
		tokenHeader := c.Request.Header.Get("Authorization")
		myJwtCustomClaims, err := ParseToken(tokenHeader)
		if tokenHeader == "" {
			code = re.ERROR_TOKEN_EIXT
		} else if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					code = re.ERROR_TOKEN_WRONG
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					code = re.ERROR_TOKEN_EXPIRED
				} else {
					code = re.ERROR_TOKEN_FAIL
				}
			}
		}
		if code != 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": re.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		c.Set("role", myJwtCustomClaims.Role)
		c.Set("uid", myJwtCustomClaims.ID)
		c.Next()
	}
}

func WSTokenVerify(tokenHeader string, c *gin.Context) {
	code := 0
	// tokenHeader := c.Request.Header.Get("Authorization")
	myJwtCustomClaims, err := ParseToken(tokenHeader)
	if tokenHeader == "" {
		code = re.ERROR_TOKEN_EIXT
	} else if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				code = re.ERROR_TOKEN_WRONG
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				code = re.ERROR_TOKEN_EXPIRED
			} else {
				code = re.ERROR_TOKEN_FAIL
			}
		}
	}
	// fmt.Println("tokenId: ", myJwtCustomClaims.ID)
	if code != 0 {
		// fmt.Println(re.GetErrMsg(code))
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		})
		c.Abort()
		return
	}
	c.Set("role", myJwtCustomClaims.Role)
	c.Set("uid", myJwtCustomClaims.ID)
	// fmt.Println("tokenId: ", myJwtCustomClaims.ID)
	c.Next()
}
