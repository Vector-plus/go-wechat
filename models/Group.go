package models

import (
	"wechat/utils/re"

	"gorm.io/gorm"
)

type Msggroup struct {
	gorm.Model
	Uid         int    `gorm:"type:bigint(20);not null" json:"uid"`
	Groupname   string `gorm:"type:varchar(250);not null" json:"groupname"`
	Description string `gorm:"type:varchar(250);not null" json:"description"`
}

func CreateGroup(msgGroup *Msggroup) int {
	err := db.Create(msgGroup).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func FindByGid(gid int) Msggroup {
	msgGooup := Msggroup{}
	db.Where("id = ?", gid).First(&msgGooup)
	return msgGooup
}

func DeleteGroup(gid int) int {
	msgGoup := Msggroup{}
	err := db.Where("id = ?", gid).Delete(&msgGoup).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func UpdateGroup(msgGooup *Msggroup) int {
	msgGroup := Msggroup{}
	gmap := make(map[string]interface{})
	gmap["uid"] = msgGooup.Uid
	gmap["groupname"] = msgGooup.Groupname
	gmap["description"] = msgGooup.Description
	err := db.Model(&msgGroup).Where("id = ?", msgGooup.ID).Updates(gmap).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func GetGroupsByName(gname string) []*Msggroup {
	groupList := make([]*Msggroup, 10)
	db.Where("Groupname like ?", gname+"%").Find(&groupList)
	return groupList
}
