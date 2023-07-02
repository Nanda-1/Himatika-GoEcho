package models

import (
	"time"

	"gorm.io/gorm"
)

// status ada 3
type Workshop struct {
	ID            int                   `json:"id"`
	JudulWorkshop string                `json:"judul_workshop"`
	Pembicara     string                `json:"pembicara"`
	WorkshopDate  string                `json:"workshop_date"`
	Status        string                `json:"status"`
	Absensi       []AbsensiWorkshop     `json:"absensi" gorm:"foreignkey:WorkshopID"`
	Pendaftaran   []PendaftaranWorkshop `json:"pendaftaran" gorm:"foreignkey:WorkshopID"`
	PropWorkshop  *PropWorkshop         `json:"prop_workshop" gorm:"foreignkey:WorkshopID"`
	LpjWorkshop   *LpjWorkshop          `json:"lpj_workshop" gorm:"foreignkey:WorkshopID"`
	Dokumentasi   *DokumentasiWorkshop  `json:"dokumentasi" gorm:"foreignkey:WorkshopID"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

type PendaftaranWorkshop struct {
	ID            int       `json:"id"`
	WorkshopID    int       `json:"workshop_id"`
	Nim           string    `json:"nim"`
	Name          string    `json:"name"`
	ProgramStudi  string    `json:"program_studi"`
	Angatan       string    `json:"angatan"`
	JudulWorkshop string    `json:"judul_workshop"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AbsensiWorkshop struct {
	ID            int       `json:"id"`
	WorkshopID    int       `json:"workshop_id"`
	Nim           string    `json:"nim"`
	Name          string    `json:"name"`
	ProgramStudi  string    `json:"program_studi"`
	JudulWorkshop string    `json:"judul_workshop"`
	Angatan       string    `json:"angatan"`
	WorkshopDate  string    `json:"workshop_date"`
	Masukan       string    `json:"masukan" gorm:"text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PropWorkshop struct {
	ID              int       `json:"id"`
	WorkshopID      int       `json:"workshop_id"`
	JudulWorkshop   string    `json:"judul_workshop"`
	Nim             string    `json:"nim"`
	Name            string    `json:"name"`
	TahunPengurusan string    `json:"tahun_pengurusan"`
	PathFile        string    `json:"path_file"`
	NameFile        string    `json:"name_file"`
	UploadComplete  bool      `json:"upload_complete"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LpjWorkshop struct {
	ID              int       `json:"id"`
	WorkshopID      int       `json:"workshop_id"`
	Nim             string    `json:"nim"`
	Name            string    `json:"name"`
	TahunPengurusan string    `json:"tahun_pengurusan"`
	PathFile        string    `json:"path_file"`
	NameFile        string    `json:"name_file"`
	UploadComplete  bool      `json:"upload_complete"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DokumentasiWorkshop struct {
	ID              int       `json:"id"`
	WorkshopID      int       `json:"workshop_id"`
	Nim             string    `json:"nim"`
	Name            string    `json:"name"`
	TahunPengurusan string    `json:"tahun_pengurusan"`
	PathFile        string    `json:"path_file"`
	NameFile        string    `json:"name_file"`
	UploadComplete  bool      `json:"upload_complete"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func CreateWorkshop(db *gorm.DB, workshop *Workshop) error {
	return db.Create(workshop).Error
}

func GetWorkshopByJudul(db *gorm.DB, judul string) (*Workshop, error) {
	workshop := &Workshop{}
	err := db.Where("judul_workshop = ?", judul).First(workshop).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Workshop tidak ditemukan
		}
		return nil, err
	}

	return workshop, nil
}

// pendaftaan
func CreatePendaftaran(db *gorm.DB, P *PendaftaranWorkshop) error {
	return db.Create(P).Error
}

// absen
func CreateAbsensi(db *gorm.DB, P *AbsensiWorkshop) error {
	return db.Create(P).Error
}

// Props
func CreateProps(db *gorm.DB, P *PropWorkshop) error {
	return db.Create(P).Error
}

// Lpj
func CreateLpj(db *gorm.DB, P *LpjWorkshop) error {
	return db.Create(P).Error
}

// Lpj
func CreateDukumentasi(db *gorm.DB, P *DokumentasiWorkshop) error {
	return db.Create(P).Error
}
