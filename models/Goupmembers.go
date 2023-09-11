package models

import (
	"wechat/utils/re"

	"gorm.io/gorm"
)

type Groupship struct {
	gorm.Model
	Gid       int    `gorm:"type:bigint(20);not null" json:"gid"`
	Groupname string `gorm:"type:varchar(250);not null" json:"groupname"`
	Uid       int    `gorm:"type:bigint(20);not null" json:"uid"`
	Username  string `gorm:"type:varchar(250);not null" json:"username"`
	Roleid    int    `gorm:"type:tinyint(10);not null" json:"roleid"` //成员角色，0为群主，1为普通成员
}

func CreateGroupShip(groupship *Groupship) int {
	err := db.Create(&groupship).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func GetGMemberByUidGid(gid, uid int) Groupship {
	var groupship Groupship
	db.Where("gid = ? and uid = ?", gid, uid).First(&groupship)
	return groupship
}

func GetAllMemeber(gid int) []*Groupship {
	userList := make([]*Groupship, 10)
	db.Where("gid = ?").Find(&userList)
	return userList
}
