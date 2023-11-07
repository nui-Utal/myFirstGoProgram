package models

import "xinyeOfficalWebsite/pkg/utils"

type Carousel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Order      int    `json:"index"`
	UploadTime string
}

func AddCarousel(c Carousel) bool {
	return utils.DB.Save(&c).RowsAffected != 1
}

func GetCarousels(page, pageSize int) (cl []Carousel, err error) {
	result := utils.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&cl)
	if result.Error != nil {
		return nil, result.Error
	}
	return cl, nil
}

func UpdateOrderByName(name string, order int) bool {
	return utils.DB.Model(&Carousel{}).Where("name = ?", name).Update("order", order).RowsAffected != 1
}

func DeleteCarousel(name string) bool {
	return utils.DB.Where("name = ?", name).Delete(&Carousel{}).RowsAffected != 1
}
