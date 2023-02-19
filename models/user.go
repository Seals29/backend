package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	IsBan     bool   `json:"isban"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type ResetUser struct {
	gorm.Model
	ExpiredDate time.Time `json:"expireddate"`
	ResetCode   int       `json:"resetcode"`
	UserID      int       `json:"userid"`
	Used        bool      `json:"used"`
}
type UserSubscribe struct {
	gorm.Model
	UserEmail string `json:"useremail"`
}
