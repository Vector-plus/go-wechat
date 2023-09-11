package models

import (
	"wechat/utils/re"

	"gorm.io/gorm"
)

type Firendship struct {
	gorm.Model
	Uid   int    `gorm:"type:bigint(20);not null" json:"uid" validate:"required"`
	Fid   int    `gorm:"type:bigint(20);not null" json:"fid" validate:"required"`
	Notes string `gorm:"type:varchar(250);not null" json:"notes"`
}

func Createship(firendship *Firendship) int {
	err := db.Create(&firendship).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func GetFirendship(uid int) []*Firendship {
	firends := make([]*Firendship, 10)
	db.Where("uid = ?", uid).Find(&firends)
	return firends
}

func DelteFirend(uid, fid int) int {
	var data Firendship
	err := db.Where("uid = ? and fid = ?", uid, fid).Delete(&data).Error
	errt := db.Where("uid = ? and fid = ?", fid, uid).Delete(&data).Error
	if err != nil || errt != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func GetByUidFid(uid, fid int) Firendship {
	var data Firendship
	db.Where("uid = ? and fid = ?", uid, fid).First(&data)
	return data
}
