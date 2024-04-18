package model

type User struct {
	ID       int    `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"not null;type:string"`
}
