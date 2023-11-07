package models

import (
	"time"
	"xinyeOfficalWebsite/pkg/utils"
)

type Admin struct {
	ID        int    `gorm:"column:id;primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
	Password  string `json:"password" gorm:"column:pwd"`
	Salt      string
	Avatar    string
	LoginTime *time.Time
	Mark      string `gorm:"column:mark"`
}

func GetAdminByName(name string) (a Admin) {
	//return utils.DB.Where("name = ?", name).Find(&a).Name() return "mysql"
	utils.DB.Where("name = ?", name).Find(&a)
	return
}

func GetAdminById(id int) (a Admin) {
	utils.DB.Find(&a, id)
	return
}

func UpdateNameById(id int, name string) bool {
	res := utils.DB.Model(&Admin{}).Where("id = ?", id).Update("name", name)
	return res.RowsAffected == 1
}

func UpdateAvatorById(id int, name string) bool {
	return utils.DB.Model(&Admin{}).Where("id = ?", id).Update("avatar", name).RowsAffected == 1
}

func UpdateAdminPwdById(id int, pwd string) bool {
	return utils.DB.Model(&Admin{}).Where("id = ?", id).Update("pwd", pwd).RowsAffected == 1
}
