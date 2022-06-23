package dal

import (
	"time"
)

type User struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	UserName  string    `json:"user_name,omitempty" gorm:"type:varchar(100)"`
	Role      string    `json:"role,omitempty" gorm:"type:varchar(10)"`
	Pwd       string    `json:"password,omitempty" gorm:"type:varchar(100)"`
	Salt      string    `json:"salt,omitempty" gorm:"type:char(4)"`
	Company   string    `json:"company,omitempty" gorm:"type:varchar(100)"`
	Phone     string    `json:"phone,omitempty" gorm:"type:varchar(100)"`
	Address   string    `json:"address,omitempty" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
}

type Information struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	Code      string    `json:"code,omitempty" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"-" gorm:"index:,sort:desc"`
}
