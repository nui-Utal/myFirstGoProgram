package models

import (
	"xinyeOfficalWebsite/pkg/utils"
)

type Collect struct {
	ID     int    `json:"id" gorm:"column:id"`
	UserId int    `json:"user" gorm:"column:user"`
	TextId int    `json:"text" gorm:"column:text"`
	Time   string `json:"time" gorm:"column:time"`
}

func AddCollect(uid, text int) bool {
	col := Collect{
		UserId: uid,
		TextId: text,
		Time:   utils.GetLocalTime(),
	}
	return utils.DB.Save(&col).RowsAffected == 1
}

func IsCollected(uid, tid int) bool {
	var c Collect
	result := utils.DB.Where("user = ? AND text = ?", uid, tid).First(&c)
	return result.RowsAffected != 0
}

func GetCollectByUid(uid, page, pageSize int) []Collect { // query 返回的参数是string
	var c []Collect
	utils.DB.Offset((page-1)*pageSize).Limit(pageSize).Where("user = ?", uid).Find(&c)
	return c
}
