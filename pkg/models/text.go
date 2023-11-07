package models

import (
	"xinyeOfficalWebsite/pkg/utils"
)

type Text struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" db:"title"`
	Label       string `json:"label" db:"label"`
	Author      int    `json:"author" db:"author"`
	Content     string `json:"content" db:"content"`
	PublishTime string `json:"publish_time" db:"publish_time"`
	View        int    `json:"view" db:"view"`
	Likes       int    `json:"like" db:"likes"`
	Collect     int    `json:"collect" db:"collect"`
	Comment     int    `json:"comment" db:"comment"`
	Type        int    `json:"type" db:"type"`
}

type Search struct {
	Textid int
	Weight int
	Text   Text
}

type ByWeight []Search

func (w ByWeight) Len() int {
	return len(w)
}

func (w ByWeight) Less(i, j int) bool {
	return w[i].Weight > w[j].Weight
}

func (w ByWeight) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func AddEssay(author, t int, title, content, lable string) int64 {
	text := Text{
		Title:       title,
		Label:       lable,
		Author:      author,
		Content:     content,
		PublishTime: utils.GetLocalTime(),
		Type:        t,
	}
	result := utils.DB.Create(&text)

	return result.RowsAffected
}

func DelTextById(tid int) bool {
	return utils.DB.Delete(&Text{}, tid).RowsAffected == 1
}

func GetAuthorById(id int) int {
	var t Text
	utils.DB.Where("id = ?", id).Find(&t)
	return t.Author
}

func GetTextById(id int, t *Text) bool {
	utils.DB.First(&t, id)
	return t.ID == 1
}

func GetEssayByUid(uid int, after, before string) []Text {
	var essay []Text
	utils.DB.Where("author = ?", uid).
		Where("publish_time > ? AND publish_time < ?", after, before).
		Order("publish_time DESC").
		Find(&essay)
	return essay
}

func GetAllEssayLimitTime(page, pageSize int, after, before string) (t []Text) {
	utils.DB.Offset((page-1)*pageSize).Limit(pageSize).
		Where("type = ?", 0).
		Where("publish_time > ? AND publish_time < ?", after, before).
		Order("publish_time DESC").
		Find(&t)
	return
}

func GetAllText(page, pageSize, t int) ([]Text, error) {
	var essay []Text
	result := utils.DB.
		Where("type = ?", t).
		Offset((page - 1) * pageSize).
		Limit(pageSize).Find(&essay)
	if result.Error != nil {
		return nil, result.Error
	}
	return essay, nil
}

func LookInTitle(key string) (res []Search, err error) {
	var t []Text
	if err = utils.DB.Where("title LIKE ?", "%"+key+"%").Find(&t).Error; err != nil {
		return // 查询结果为空不抛出异常
	}
	if len(t) == 0 {
		return
	}
	for _, text := range t {
		found := false
		// range循环判断这个结果是否被找到过
		for i := range res {
			if res[i].Textid == text.ID {
				res[i].Weight += 2
				found = true
				break
			}
		}
		// 未收录的结果
		if !found {
			r := Search{
				Textid: text.ID,
				Weight: 2,
				Text:   text,
			}
			res = append(res, r)
		}
	}
	return
}
