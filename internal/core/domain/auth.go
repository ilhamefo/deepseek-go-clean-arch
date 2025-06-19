package domain

import "time"

type AuthRepository interface {
	IsRegistered(email string) (isRegistered bool, err error)
	Register(user User) (err error)
}

type User struct {
	ID              string     `json:"id" gorm:"column:id"`
	Email           string     `json:"email" gorm:"column:email"`
	Password        string     `json:"password" gorm:"column:password"`
	VerifiedEmail   bool       `json:"verified_email" gorm:"-"`
	Name            string     `json:"name" gorm:"column:name"`
	GivenName       string     `json:"given_name" gorm:"-"`
	FamilyName      string     `json:"family_name" gorm:"-"`
	Picture         string     `json:"picture" gorm:"column:picture"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"column:email_verified_at"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

func (a *User) TableName() string {
	return "public.users"
}
