package models

import (
	"fmt"
	"wechat/utils/re"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12"`
	Password string `gorm:"type:varchar(250);not null" json:"password" validate:"required,min=6,max=20"`
	Phone    string `gorm:"type:varchar(20);not null" json:"phone" validate:"required"`
	Salt     string `gorm:"type:varchar(40);not null"`
	Roleid   int    `gorm:"type:int;DEFAULT:2;not null" json:"roleId" validate:"required"`
}

func FindByUid(uid int) (User, int) {
	var user User
	err := db.Limit(1).Where("ID = ?", uid).Find(&user).Error
	if err != nil {
		fmt.Println("TestGet error[22]", err)
		return user, re.ERROR
	}
	fmt.Println(user)
	return user, re.SUCCSE
}

func AddUser(data *User) int {
	err := db.Create(&data).Error
	if err != nil {
		fmt.Println("Adduser error", err)
		return re.ERROR
	}
	// fmt.Println(data.ID)
	return re.SUCCSE
}

func FindByPhone(phone string) User {
	user := User{}
	db.Where("phone = ?", phone).First(&user)
	return user
}

func DeleteUser(uid int) int {
	var user User
	err := db.Where("ID = ?", uid).Delete(&user).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}

func GetUserList() []*User {
	userList := make([]*User, 10)
	db.Find(&userList)
	return userList
}

func UpdateUser(data *User) int {
	var user User
	var umap = make(map[string]interface{})
	umap["username"] = data.Username
	umap["phone"] = data.Phone
	umap["password"] = data.Password
	err = db.Model(&user).Where("id = ?", data.ID).Updates(umap).Error
	if err != nil {
		return re.ERROR
	}
	return re.SUCCSE
}
