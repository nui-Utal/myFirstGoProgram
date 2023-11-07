package models

import (
	"xinyeOfficalWebsite/pkg/utils"
)

type Comment struct {
	ID               int `json:"comid" gorm:"primary_key;autoIncrement"`
	Text             int `json:"textid" gorm:"index;"`
	Commentator_id   int
	Commentator_name string
	Comment          string `json:"content"`
	CommentTime      string `gorm:"comment:'评论发表时间'"`
	Replied          int
}

func GetCommentById(comid int) Comment {
	var comments Comment
	utils.DB.Where("id = ?", comid).
		Find(&comments)
	return comments
}

func GetCommentLimitTime(after, before string) (c []Comment) {
	utils.DB.
		Where("comment_time > ? AND comment_time < ?", after, before).
		Order("comment_time DESC").
		Find(&c)
	return
}

func GetCommentByUidLimitTime(uid int, after, before string) []Comment {
	var comments []Comment
	utils.DB.Where("commentator_id = ?", uid).
		Where("comment_time > ? AND comment_time < ?", after, before).
		Order("comment_time DESC").
		Find(&comments)
	return comments
}

func AddComment(c Comment) error {
	err := utils.DB.Save(&c).Error
	if err != nil {
		return err
	}
	return nil
}

func DelCommentById(id int) bool {
	return utils.DB.Delete(&Comment{}, id).RowsAffected == 1
}

func GetCommentByTid(tid string) []Comment {
	var c []Comment
	utils.DB.Where("text = ?", tid).Find(&c)
	return c
}

func IsCommented(cid int) (bool, int) {
	var c Comment
	res := utils.DB.Where("id = ?", cid).Find(&c)
	return res.RowsAffected == 1, c.Commentator_id
}
