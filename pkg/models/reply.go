package models

import (
	"fmt"
	"xinyeOfficalWebsite/pkg/utils"
)

type Reply struct {
	ID         int    `json:"rid" gorm:"column:id;primaryKey;autoIncrement"`
	Comment    int    `json:"comid"`
	Send       int    `gorm:"column:send"`
	Receive    int    `json:"parentid" gorm:"column:receive"`
	Content    string `json:"content" gorm:"column:content;not null"`
	Reply_time string
}

// 设置表名
func (Reply) TableName() string {
	return "replies"
}

func AddReply(r Reply) bool {
	return utils.DB.Create(&r).RowsAffected == 1
}

func DelReplyById(id int) error {
	result := utils.DB.Delete(&Reply{}, id)
	if result.Error != nil {
		// 删除过程中发生了错误 return result.Error
		return fmt.Errorf("删除失败")
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到对应的记录")
	}
	return nil
}

func GetReplyById(id int) Reply {
	var rep Reply
	utils.DB.Where("id = ?", id).
		Find(&rep)
	return rep
}

func GetReplyByCid(cid int) []Reply {
	var rep []Reply
	utils.DB.Where("comment = ?", cid).
		Order("reply_time ASC").
		Find(&rep)
	return rep
}

func GetReplyByUidLimitTime(sender int, after, before string) []Reply {
	var rep []Reply
	utils.DB.Where("sender = ?", sender).
		Where("reply_time > ? AND reply_time < ?", after, before).
		Order("comment_time DESC").
		Find(&rep)
	return rep
}
