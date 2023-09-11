package models

import (
	"wechat/utils/re"

	"gorm.io/gorm"
)

type FriendApplication struct {
	gorm.Model
	ApplicantId     int    `gorm:"type:bigint(20);not null" json:"applicantid" validate:"required"`
	ApplicantName   string `gorm:"type:varchar(250);not null" json:"applicantname" validate:"required"`
	ReviewerId      int    `gorm:"type:bigint(20);not null" json:"reviewerid" validate:"required"`
	ReviewerName    string `gorm:"type:varchar(250);not null" json:"reviewername" validate:"required"`
	Status          int    `gorm:"type:tinyint(5);not null;DEFAULT:0" json:"status"`          //0,未审核；1,通过；2,拒绝
	ApplicationType int    `gorm:"type:tinyint(5);not null;DEFAULT:1" json:"applicationType"` //1,好友申请；2,群聊申请
}

func GetApplicationList(uid int) []*FriendApplication {
	applicationList := make([]*FriendApplication, 10)
	db.Where("applicant_id = ?", uid).Find(&applicationList)
	return applicationList
}

func AddAppli(firendAppli *FriendApplication) int {
	err := db.Create(&firendAppli).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func GetApplicationListByUser(uid int) []*FriendApplication {
	applicationList := make([]*FriendApplication, 10)
	db.Where("reviewer_id = ? and application_type = ?", uid, 1).Find(&applicationList)
	return applicationList
}

func GetAppliById(aid int) FriendApplication {
	application := FriendApplication{}
	db.Where("id = ?", aid).First(&application)
	return application
}

func UpdateAppli(application *FriendApplication) int {
	var firendA FriendApplication
	var amap = make(map[string]interface{})
	amap["status"] = application.Status
	err := db.Model(&firendA).Where("id = ?", application.ID).Updates(amap).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}
