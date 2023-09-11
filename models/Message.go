package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"wechat/utils"
	"wechat/utils/middleware"
	"wechat/utils/re"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	//发送者
	Fromid int `gorm:"type:bigint(20);not null" json:"fromid"`
	//消息接收者ID
	Targetid int `gorm:"type:bigint(20);not null" json:"targetid"`
	//消息类型
	Msgtype int `gorm:"type:tinyint(5);not null" json:"msgtype"` //消息种类：1,文字 2,图片 3,视频。。。。
	//消息种类
	Msgkind int `gorm:"type:tinyint(5);not null" json:"msgkind"` //消息类型：1,私聊 2,群聊 3,系统。。。
	//消息内容
	Content string `gorm:"type:varchar(250);not null" json:"content"` //消息内容
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var clientMap map[int]*Node = make(map[int]*Node, 0)

var rwLocker sync.RWMutex

var udpsendChan chan []byte = make(chan []byte, 1024)

// 功能初始化，调度函数
func init() {
	go udpSendProc()
	go udpReciveProc()
}

func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(utils.NetaddrA, utils.NetaddrB, utils.NetaddrC, utils.NetaddrD),
		Port: utils.NetPortSend,
	})
	if err != nil {
		fmt.Println("udpsendPro error", err)
		return
	}
	defer con.Close()
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println("udpWriteMsg error", err)
				return
			}
			fmt.Println(data)
			dispatch(data)
		}
	}
}

func udpReciveProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: utils.NetPortRecive,
	})
	if err != nil {
		fmt.Println("udpReciveProc error", err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println("udpreciveProc error", err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 接受消息之后的处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("dispatch error!", err)
		return
	}
	switch msg.Msgtype {
	//私聊消息
	case 1:
		sendPersonMsg(&msg, msg.Targetid)
	case 2:
		sendGroupMsg(&msg)
	case 3:
		broadMsg(&msg)
	}
	db.Create(&msg)
}

func sendPersonMsg(msg *Message, targetid int) {
	rwLocker.RLock()
	node, flag := clientMap[targetid]
	rwLocker.RUnlock()
	if flag {
		node.DataQueue <- []byte(msg.Content)
	}
}

func sendGroupMsg(msg *Message) {
	memeberList := GetAllMemeber(msg.Targetid)
	for _, memeber := range memeberList {
		sendPersonMsg(msg, memeber.Uid)
	}
}

func broadMsg(msg *Message) {
	memeberList := GetUserList()
	for _, user := range memeberList {
		sendPersonMsg(msg, int(user.ID))
	}
}

// 建立websoket连接
func Chat(c *gin.Context) int {
	request := c.Request
	writer := c.Writer
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{request.Header.Get("Sec-WebSocket-Protocol")},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println("Chat error: ", err)
		return re.ERROR_WS_CONN
	}
	fmt.Println(conn.Subprotocol())
	token := conn.Subprotocol()
	middleware.WSTokenVerify(token, c)
	uid := c.GetInt("uid")
	//生成用户连接node
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	rwLocker.Lock()
	clientMap[uid] = node
	rwLocker.Unlock()

	//新建发送、接收消息的协程
	go sendProc(node)
	go reciveProc(node)
	fmt.Println("[userId]--->", uid)
	sendMsg(uid, []byte("welcom to ginchat"))
	return re.SUCCSE
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("sendMsgPro error: ", err)
				return
			}
		}
	}
}

func reciveProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println("reciveProc error: ", err)
			return
		}
		udpsendChan <- data
	}
}

func GetFHistoryMsg(uid, fid int) ([]*Message, int) {
	msgList := make([]*Message, 10)
	err := db.Where("((fromid = ? and targetid = ?) or (fromid = ? and targetid = ?)) and msgkind = 1", uid, fid, fid, uid).Find(&msgList).Error
	if err != nil {
		return msgList, re.ERROR
	}
	return msgList, re.SUCCSE
}

func GetGHistoryMsg(gid int) ([]*Message, int) {
	msgList := make([]*Message, 10)
	err := db.Where("targetid = ? and msdtype = 2", gid).Find(&msgList)
	if err != nil {
		return msgList, re.ERROR
	}
	return msgList, re.SUCCSE

}

func sendMsg(userId int, msg []byte) {
	rwLocker.RLock()

	node, flag := clientMap[userId]
	rwLocker.RUnlock()
	if flag {
		node.DataQueue <- msg
	}
}
