package models

type Photo struct {
	Id       string `gorm:"primaryKey"`
	Title    string `gorm:"type:varchar(300)" form:"title" valid:"required"`
	Caption  string `gorm:"type:varchar(300)" form:"caption" valid:"required"`
	PhotoUrl string `gorm:"type:varchar(300)" valid:"required"`
	UserID   int64
	User     User
}

type UpdatePhoto struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
}
