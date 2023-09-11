package utils

import (
	"fmt"

	ini "gopkg.in/ini.v1"
)

// 系统变量
var (
	AppMode   string
	HttpPort  string
	JwtKey    string
	TokenTime int

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	Addr        string
	Port        string
	Password    string
	DBNum       string
	PoolSize    string
	MinIdelConn string

	NetaddrA      byte
	NetaddrB      byte
	NetaddrC      byte
	NetaddrD      byte
	NetPortSend   int
	NetPortRecive int

	FILEDIRE string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("variable.go err [26]", err)
	}

	//初始化系统参数
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8080")
	JwtKey = file.Section("server").Key("JwtKey").MustString("123456789@FHT")
	TokenTime = file.Section("server").Key("TokenTime").MustInt(60)

	//初始化数据库参数
	// DB = file.Section("database").Key("DB").MustString("mysql")
	DBHost = file.Section("database").Key("DBHost").MustString("127.0.0.1")
	DBPort = file.Section("database").Key("DBPort").MustString("3306")
	DBUser = file.Section("database").Key("DBUser").MustString("root")
	DBPassword = file.Section("database").Key("DBPassword").MustString("root")
	DBName = file.Section("database").Key("DBName").MustString("ginchat")

	//Redis数据库参数
	Addr = file.Section("redis").Key("Addr").MustString("127.0.0.1")
	Port = file.Section("redis").Key("Port").MustString("6379")
	Password = file.Section("redis").Key("Password").MustString("")
	DBNum = file.Section("redis").Key("DBNum").MustString("0")
	PoolSize = file.Section("redis").Key("PoolSize").MustString("30")
	MinIdelConn = file.Section("redis").Key("MinIdelConn").MustString("300")

	//wbsocket参数设置
	NetaddrA = byte(file.Section("ws").Key("NetaddrA").MustInt(127))
	NetaddrB = byte(file.Section("ws").Key("NetaddrB").MustInt(0))
	NetaddrC = byte(file.Section("ws").Key("NetaddrC").MustInt(0))
	NetaddrD = byte(file.Section("ws").Key("NetaddrD").MustInt(255))
	NetPortSend = file.Section("ws").Key("NetPortSend").MustInt(3000)
	NetPortRecive = file.Section("ws").Key("NetPortRecive").MustInt(3001)

	FILEDIRE = file.Section("filedir").Key("Path").MustString("./var/upload")
}
