package common

import "xinyeOfficalWebsite/pkg/models"

type UserContainer struct {
	Id        int
	Username  string
	Sno       string
	Special   string
	Avatar    string
	Official  int
	FanNum    int
	FollowNum int
}

type TextContainer struct {
	Liked   bool
	Collect bool
	Text    models.Text
}

func UserConversion(u models.User) UserContainer {
	uc := UserContainer{
		Id:        u.ID,
		Username:  u.Username,
		Sno:       u.Sno,
		Special:   u.Special,
		Avatar:    GetPath("head") + u.Avatar,
		FanNum:    u.FanNum,
		FollowNum: u.FollowNum,
	}
	return uc
}

func BatchUserConversion(u []models.User) []UserContainer {
	var us []UserContainer
	for _, user := range u {
		us = append(us, UserConversion(user))
	}
	return us
}
