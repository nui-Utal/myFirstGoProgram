package models

import (
	"xinyeOfficalWebsite/pkg/utils"
)

type Follow struct {
	Followed int    `json:"userid" gorm:"column:followed"`
	Fan      int    `json:"curuserid" gorm:"column:fan"`
	Time     string `gorm:"type:datetime;column:time"`
}

// 添加关注
func AddFollow(followId, fanId int) {
	f := Follow{
		Followed: followId,
		Fan:      fanId,
		Time:     utils.GetLocalTime(),
	}
	utils.DB.Create(&f)
}

// 取消关注
func DelFollow(followId, fanId int) {
	utils.DB.Exec("DELETE FROM follows WHERE followed = ? AND fan = ?", followId, fanId)

}

// 查看粉丝
func ShowFans(follow, page, pageSize int) []User {
	var users []User
	var f []Follow
	utils.DB.Offset((page-1)*pageSize).Limit(pageSize).Where("followed = ?", follow).Find(&f)
	for _, fan := range f {
		users = append(users, GetUserById(fan.Fan))
	}
	return users
}

// 查看关注
func ShowFollow(fan, page, pageSize int) []User {
	var users []User
	var f []Follow
	utils.DB.Offset((page-1)*pageSize).Limit(pageSize).Where("fan = ?", fan).Find(&f)
	for _, fan := range f {
		users = append(users, GetUserById(fan.Followed))
	}
	return users
}

func IsFollowed(followed, fan int) int64 {
	var f Follow
	res := utils.DB.Where("followed = ? AND fan = ?", followed, fan).Find(&f)
	return res.RowsAffected
}
