package models

import (
	"echoGo/db"
	"time"
)

type Anggota struct {
	ID             int             `json:"id" gorm:"primarykey"`
	Nim            string          `json:"nim" gorm:"unique"`
	Password       string          `json:"password"`
	RoleID         int             `json:"role_id"`
	Role           *Roles          `json:"role" gorm:"foreignKey:RoleID"`
	AnggotaDetails *AnggotaDetails `json:"anggota_details" gorm:"foreignKey:AnggotaID"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type AnggotaDetails struct {
	ID              int       `json:"id"`
	AnggotaID       int       `json:"anggota_id" gorm:"index"`
	Name            string    `json:"name"`
	TahunPengurusan string    `json:"tahun_pengurusan"`
	Divisi          string    `json:"divisi"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Roles struct {
	ID   int    `json:"id" gorm:"primarykey"`
	Name string `json:"name"`
}
type ReqAuth struct {
	AccesToken   string `json:"acces_token"`
	RefreshToken string `json:"refresh_token"`
}

func FindUserByNIM(Nim string) (*Anggota, error) {
	var user Anggota
	err := db.Db.Where("nim = ?", Nim).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
