package models

import (
	"xinyeOfficalWebsite/pkg/utils"
)

type Like struct {
	User     int
	Text     int
	LikeTime string
}

func AddLike(uid, tid int) error {
	l := Like{
		User:     uid,
		Text:     tid,
		LikeTime: utils.GetLocalTime(),
	}
	return utils.DB.Create(&l).Error
}

func GetLikeByUid(uid, page, pagesize int) (l []Like) {
	utils.DB.Offset((page-1)*pagesize).Limit(pagesize).Where("user = ?", uid).Find(&l)
	return
}

func DelLike(uid, tid int) error {
	return utils.DB.Where("user = ? AND text = ?", uid, tid).Delete(&Like{}).Error
}

func IsLiked(uid, tid int) bool {
	var c Like
	result := utils.DB.Where("user = ? AND text = ?", uid, tid).Find(&c)
	return result.RowsAffected != 0
}
