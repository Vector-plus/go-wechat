package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"wechat/models"
	"wechat/utils/re"
	"wechat/utils/upload"

	"github.com/gin-gonic/gin"
)

type MessageParam struct {
	Fromid   int    `json:"fromid"`
	Targetid int    `json:"targetid"`
	Msgtype  int    `json:"msgtype"`  //消息种类：1,文字 2,图片 3,视频。。。。
	Msgkind  int    `json:"msgkind"`  //消息类型：1,私聊 2,群聊 3,系统。。。
	Content  string `json:"content"`  //消息内容
	UserInfo string `json:"userinfo"` //用户信息
}

// GetConn
// @Summary 建立websocket连接
// @Tags 消息模块
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/getwsconn [get]
func GetConn(c *gin.Context) {
	code := models.Chat(c)
	fmt.Println(c.GetInt("uid"))
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
		},
	)
}

// Uploadfile
// @Summary 文件上传
// @Tags 消息模块
// @param file formData file true "file"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/uploadfile [post]
func UploadFile(c *gin.Context) {
	code, path := upload.UPloadfile(c.Writer, c.Request)
	c.JSON(
		http.StatusOK, gin.H{
			"status":   code,
			"message":  re.GetErrMsg(code),
			"filepath": path,
		},
	)
}

// DownloadFile
// @Summary 文件下载
// @Tags 消息模块
// @param filepath query string true "file路径"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/downloadfile [get]
func DownloadFile(c *gin.Context) {
	filepath := c.Query("filepath")
	filename := filepath
	tem := strings.Split(filename, "/")
	if len(tem) > 1 {
		filename = tem[len(tem)-1]
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filepath)
}

// GetFHistroryMsg
// @Summary 获取与的好友历史消息
// @Tags 消息模块
// @param fid query int true "好友id"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/getFhistoryMsg [get]
func GetFHistroryMsg(c *gin.Context) {
	uid := c.GetInt("uid")
	fmt.Println(uid)
	fid, _ := strconv.Atoi(c.Query("fid"))
	msgList, _ := models.GetFHistoryMsg(uid, fid)
	user, _ := models.FindByUid(uid)
	fuser, code := models.FindByUid(fid)
	c.JSON(
		http.StatusOK, gin.H{
			"status":     code,
			"message":    re.GetErrMsg(code),
			"username":   user.Username,
			"friendname": fuser.Username,
			"msgList":    msgList,
		},
	)
}

// GetGHistroryMsg
// @Summary 获取群聊历史消息
// @Tags 消息模块
// @param gid query int true "群聊id"
// @Security ApiKeyAuth
// @Success 200 {string} json{"code","message"}
// @Router /api/getGhistoryMsg [get]
func GetGHistroryMsg(c *gin.Context) {
	uid := c.GetInt("uid")
	gid, _ := strconv.Atoi(c.Query("gid"))
	group := models.GetGMemberByUidGid(gid, uid)
	if group.Groupname == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  re.ERROR_USER_EIXT,
				"message": "该用户没在群聊中",
			},
		)
	}
	msgList, code := models.GetGHistoryMsg(gid)
	memeberList := models.GetAllMemeber(gid)
	var memebermap = make(map[int]string)
	for _, v := range memeberList {
		memebermap[v.Uid] = v.Username
	}
	msgPramList := make([]*MessageParam, len(msgList))
	for i, v := range msgList {
		msgPramList[i].Fromid = v.Fromid
		msgPramList[i].Targetid = v.Targetid
		msgPramList[i].Msgtype = v.Msgtype
		msgPramList[i].Msgkind = v.Msgkind
		msgPramList[i].Content = v.Content
		msgPramList[i].UserInfo = memebermap[v.Fromid]
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": re.GetErrMsg(code),
			"msgList": msgList,
		},
	)
}
