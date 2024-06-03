package models

import "time"

type User struct {
	Id        int64     `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(300)" form:"username" valid:"required"`
	Email     string    `gorm:"type:varchar(300);unique" form:"email" valid:"required, email"`
	Password  string    `gorm:"type:varchar(300)" form:"password" valid:"required, minstringlength(6)"`
	ProfilUrl string    `gorm:"type:varchar(300)"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	Photos    []Photo   `gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE;"`
}

type LoginUser struct {
	Email    string `json:"email" valid:"required, email"`
	Password string `json:"password" valid:"required, minstringlength(6)"`
}

type UpdateUser struct {
	Username  string `form:"username"`
	Password  string `form:"password" valid:"minstringlength(6)"`
	ProfilUrl string
	UpdatedAt time.Time
}
