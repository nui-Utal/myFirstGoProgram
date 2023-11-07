package common

import (
	"gorm.io/gorm"
	"xinyeOfficalWebsite/pkg/utils"
)

func IncrementField(table, field string, id int) error {
	return utils.DB.Table(table).
		Where("id = ?", id).UpdateColumn(field, gorm.Expr(field+" + ?", 1)).Error
}

func DecrementField(table, field string, id int) error {
	return utils.DB.Table(table).
		Where("id = ?", id).UpdateColumn(field, gorm.Expr(field+" - ?", 1)).Error
}
