package domain

import (
	"time"
)

type UserRepository interface {
	Search(key string) (user []*UserVCC, err error)
	Roles() (user []*Role, err error)
	Unit(level string) (units []*UnitName, err error)
	Update(user *UserVCC) (err error)
}

type UserVCC struct {
	ID              string     `gorm:"primaryKey;column:id" json:"id"`
	Email           string     `gorm:"column:email" json:"email"`
	Username        string     `gorm:"column:username" json:"username"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at" json:"email_verified_at"`
	ApiToken        *string    `gorm:"column:api_token" json:"api_token"`
	LastLogin       *time.Time `gorm:"column:last_login" json:"last_login"`
	FullName        string     `gorm:"column:full_name" json:"full_name"`
	Password        string     `gorm:"column:password" json:"password"`
	Jabatan         string     `gorm:"column:jabatan" json:"jabatan"`
	NIP             string     `gorm:"column:nip" json:"nip"`
	Level           uint       `gorm:"column:level" json:"level"`
	UnitCode        *string    `gorm:"column:unit_code" json:"unit_code"`
	UnitName        *string    `gorm:"column:unit_name" json:"unit_name"`
	Status          uint       `gorm:"column:status" json:"status"`
	RejectedAt      *time.Time `gorm:"column:rejected_at" json:"rejected_at"`
	RememberToken   *string    `gorm:"column:remember_token" json:"remember_token"`
	Phone           *string    `gorm:"column:phone" json:"phone"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Roles           []*Role    `gorm:"many2many:dashboard.role_users;joinForeignKey:user_id;joinReferences:role_id" json:"roles"`
	RoleID          string     `json:"role_id,omitempty"`
}

func (a *UserVCC) TableName() string {
	return "dashboard.users"
}

type Role struct {
	ID                  string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name                string         `gorm:"column:name" json:"name" validate:"required"`
	Description         string         `gorm:"column:description" json:"description" validate:"required"`
	Level               int            `gorm:"column:level" json:"level" validate:"required"`
	IsEnabled           *int           `gorm:"column:is_enabled;default:1" json:"is_enabled"`
	IsCommandCenter     *int           `gorm:"column:is_command_center;default:0" json:"is_command_center"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
	Permissions         []Permissions  `gorm:"many2many:permission_roles;joinForeignKey:role_id;joinReferences:permission_id" json:"permissions,omitempty"`
	PermissionsRelation []*Permissions `gorm:"foreignKey:parent_id" json:"permissions_relation,omitempty"`
}

func (a *Role) TableName() string {
	return "dashboard.roles"
}

type Permissions struct {
	ID         string         `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name       string         `gorm:"column:name" json:"name" validate:"required"`
	Slug       string         `gorm:"column:slug" json:"slug" validate:"required"`
	Parent_id  string         `gorm:"column:parent_id" json:"parent_id"`
	Created_at *time.Time     `gorm:"column:created_at" json:"created_at"`
	Updated_at *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	Children   []*Permissions `gorm:"foreignKey:parent_id" json:"children,omitempty"`
	Parent     *Permissions   `gorm:"foreignKey:Parent_id" json:"parent,omitempty"`
}

type UnitName struct {
	Label string `gorm:"column:label" json:"label" validate:"required"`
	Code  string `gorm:"column:code" json:"code" validate:"required"`
}

type RoleUsers struct {
	UserID string `gorm:"column:user_id" json:"user_id"`
	RoleID string `gorm:"column:role_id" json:"role_id"`
}

func (a *RoleUsers) TableName() string {
	return "dashboard.role_users"
}
