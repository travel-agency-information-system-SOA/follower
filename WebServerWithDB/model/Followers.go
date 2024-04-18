package model

type Followers struct {
	UserId            string `json:"userId" gorm:"not null;type:string"`
	Username          string `json:"username" gorm:"not null;type:string"`
	FollowingUserId   string `json:"followingUserId" gorm:"not null;type:string"`
	FollowingUsername string `json:"followingUsername" gorm:"not null;type:string"`
}
