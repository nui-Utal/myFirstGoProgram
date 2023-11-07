package models

import (
	"gorm.io/gorm"
	"time"
	utils2 "xinyeOfficalWebsite/pkg/utils"
)

type User struct {
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Username   string     `json:"username,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	Sno        string     `json:"sno,omitempty"`
	Special    string     `json:"special,omitempty"`
	Pwd        string     `json:"pwd,omitempty" mysql:"default:123456"`
	DeletedAt  *time.Time `gorm:"type:datetime"`
	UpdateTime string     `json:"update_time,omitempty"`
	Avatar     string     `json:"avatar,omitempty"`
	Salt       string     `json:"salt,omitempty"`
	FanNum     int        `json:"fan_num,omitempty" mysql:"default:0"`
	FollowNum  int        `json:"follow_num,omitempty" mysql:"default:0"`
}

type Temp struct {
	Pwd  string
	Salt string
}

// 添加用户
func CreateUser(username, phone, sno, special, pwd string) (error, User) {
	salt, _ := utils2.GenerateSalt(8)
	u := User{
		Username:   username,
		Phone:      phone,
		Sno:        sno,
		Special:    special,
		Pwd:        utils2.Md5Encrypt(pwd + salt),
		Salt:       salt,
		UpdateTime: utils2.GetLocalTime(),
	}
	result := utils2.DB.Model(&User{}).Create(&u)
	return result.Error, u
}

// 通过id获取用户
func GetUserById(id int) User {
	var u User
	utils2.DB.Model(User{}).Where("id = ?", id).Find(&u)
	return u
}

// 通过学号获取用户
func GetUserBySno(sno string) User {
	var u User
	utils2.DB.Model(User{}).Where("sno = ?", sno).Find(&u)
	return u
}

// 通过用户名获取用户
func GetUserByName(sno string) User {
	var u User
	utils2.DB.Model(User{}).Where("username = ?", sno).Find(&u)
	return u
}

// 通过用户名搜索用户
func SearchUserByName(name string) []User {
	var u []User
	utils2.DB.Model(User{}).Where("username LIKE ?", name+"%").Find(&u)
	return u
}

// 获取用户列表
func GetUsers(page, pageSize int) ([]User, error) {
	var users []User
	result := utils2.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// 使用用户名和密码/登录
func GetSaltAndPwd(field, value string) (salt, pwd string, e error) {
	var temp Temp

	if err := utils2.DB.Model(&User{}).Where(field+" = ?", value).Find(&temp).Error; err != nil {
		return "", "", err
	}

	return temp.Salt, temp.Pwd, nil
}

// 使用密码登录
func LoginWithPassword(pwd string) User {
	var u User
	utils2.DB.Model(&User{}).Where("pwd = ?", pwd).First(&u).Update("update_time", utils2.GetLocalTime())
	return u
}

func GetSaltById(id int) string {
	var u User
	utils2.DB.Model(&User{}).Where("id = ?", id).Find(&u)
	return u.Salt
}

// 更新密码
func UpdatePwdById(id int, pwd string) (u User) {
	utils2.DB.Model(&User{}).Where("id = ?", id).Update("pwd", pwd).Find(&u)
	return
}

// 修改手机号
func UpdatePhoneById(id int, phone string) User {
	var u User
	utils2.DB.Model(&User{}).Where("id = ?", id).Update("phone", phone).Find(&u)
	return u
}

// 修改学号
func UpdateSnoById(id int, sno string) User {
	var u User
	utils2.DB.Model(&User{}).Where("id = ?", id).Update("sno", sno).Find(&u)
	return u
}

// 更新头像
func UploadAvatar(fileName string, id int) {
	utils2.DB.Model(&User{}).Where("id = ?", id).Update("avatar", fileName)
}

// 查看头像
func GetAvatar(id int) string {
	var u User
	utils2.DB.Model(&User{}).Where("id = ?", id).Find(&u)
	return u.Avatar
}

func DeleteUser(userID uint) error {
	// 使用 GORM 提供的 Delete 方法进行逻辑删除
	result := utils2.DB.Delete(&User{}, userID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
