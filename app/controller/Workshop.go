package controller

import (
	"echoGo/app/models"
	"echoGo/db"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateWorkshop(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	req := models.Workshop{}
	if err := c.Bind(&req); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	if req.JudulWorkshop == "" || req.Pembicara == "" || req.WorkshopDate == "" {
		errorMsg := "Tidak Boleh Kosong"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	// Check if workshop with the same title already exists
	existingWorkshop, err := models.GetWorkshopByJudul(db.Db, req.JudulWorkshop)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	if existingWorkshop != nil {
		errorMsg := "Judul Workshop sudah ada"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}
	// parsedDate, err := time.Parse("2006-01-02", req.WorkshopDate.Format("2006-01-02"))
	// if err != nil {
	// 	errorMsg := "Format tanggal tidak valid"
	// 	res.Success = false
	// 	res.Error = &errorMsg
	// 	return c.JSON(http.StatusBadRequest, res)
	// }

	CreateWorkshop := models.Workshop{
		JudulWorkshop: req.JudulWorkshop,
		Pembicara:     req.Pembicara,
		WorkshopDate:  req.WorkshopDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := models.CreateWorkshop(db.Db, &CreateWorkshop); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Success = true
	res.Data = map[string]interface{}{
		"Judul Workshop": CreateWorkshop.JudulWorkshop,
		"Pembicara":      CreateWorkshop.Pembicara,
		"Workshop Date":  CreateWorkshop.WorkshopDate,
	}
	return c.JSON(http.StatusOK, res)
}

func GetAllWorkshop(c echo.Context) error {
	workshops := []models.Workshop{}

	if err := db.Db.Preload("Pendaftaran").Preload("Absensi").Preload("PropWorkshop").Preload("LpjWorkshop").Preload("Dokumentasi").Find(&workshops).Error; err != nil {
		errorMsg := err.Error()
		res := models.JsonResponse{
			Success: false,
			Error:   &errorMsg,
		}
		return c.JSON(http.StatusInternalServerError, res)
	}

	res := models.JsonResponse{
		Success: true,
		Data:    workshops,
	}
	return c.JSON(http.StatusOK, res)
}

func CreatePendaftaran(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	req := models.PendaftaranWorkshop{}
	if err := c.Bind(&req); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	// Mencari workshop berdasarkan judul
	judul, err := models.GetWorkshopByJudul(db.Db, req.JudulWorkshop)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	if judul == nil {
		errorMsg := "Judul tidak ditemukan"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusNotFound, res)
	}

	// Membuat Pendaftaran baru
	newPendaftaran := models.PendaftaranWorkshop{
		WorkshopID:    judul.ID,
		Nim:           req.Nim,
		Name:          req.Name,
		ProgramStudi:  req.ProgramStudi,
		Angatan:       req.Angatan,
		JudulWorkshop: req.JudulWorkshop,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Simpan Pendaftaran ke database
	if err := models.CreatePendaftaran(db.Db, &newPendaftaran); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Mengupdate Workshop dengan menambahkan Pendaftaran baru ke slice Pendaftaran
	judul.Pendaftaran = append(judul.Pendaftaran, newPendaftaran)
	if err := db.Db.Save(judul).Error; err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Data = newPendaftaran
	return c.JSON(http.StatusOK, res)
}

func CreateAbsensi(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	req := models.AbsensiWorkshop{}
	if err := c.Bind(&req); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	// Mencari workshop berdasarkan judul
	judul, err := models.GetWorkshopByJudul(db.Db, req.JudulWorkshop)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	if judul == nil {
		errorMsg := "Judul tidak ditemukan"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusNotFound, res)
	}

	// Membuat Absensi baru

	absensi := models.AbsensiWorkshop{
		WorkshopID:    judul.ID,
		Nim:           req.Nim,
		Name:          req.Name,
		ProgramStudi:  req.ProgramStudi,
		JudulWorkshop: req.JudulWorkshop,
		Angatan:       req.Angatan,
		WorkshopDate:  req.WorkshopDate,
		Masukan:       req.Masukan,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Simpan Absensi ke database
	if err := models.CreateAbsensi(db.Db, &absensi); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	// Mengupdate Workshop dengan menambahkan Absensi baru ke slice Absensi
	judul.Absensi = append(judul.Absensi, absensi)
	if err := db.Db.Save(judul).Error; err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Data = absensi
	return c.JSON(http.StatusOK, res)
}

func UploadProposal(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	tema := c.FormValue("judul_workshop")
	nim := c.FormValue("nim")
	name := c.FormValue("name")
	ThnPengurusan := c.FormValue("tahun_pengurusan")

	// Mencari workshop berdasarkan judul
	judul, err := models.GetWorkshopByJudul(db.Db, tema)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	if judul == nil {
		errorMsg := "Judul tidak ditemukan"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusNotFound, res)
	}

	fileDir := "./uploads"
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, os.ModePerm)
	}

	form, err := c.MultipartForm()
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	files := form.File["file"]
	if len(files) == 0 {
		errorMsg := "No files uploaded. Please upload at least one file."
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	var uploads []*models.PropWorkshop
	for _, file := range files {
		fileUUID := uuid.New().String() // Membuat UUID unik untuk setiap file
		// uuidString := fileUUID.String()
		filename := fileUUID + "_" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "_") // Mengganti spasi dengan garis bawah
		filePath := filepath.Join(fileDir, filename)

		src, err := file.Open()
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		upload := &models.PropWorkshop{
			WorkshopID:      judul.ID,
			JudulWorkshop:   tema,
			Nim:             nim,
			Name:            name,
			TahunPengurusan: ThnPengurusan,
			PathFile:        filePath,
			NameFile:        filename,
			UploadComplete:  true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		uploads = append(uploads, upload)
	}

	for _, upload := range uploads {
		if err := models.CreateProps(db.Db, upload); err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		// Mengupdate Workshop dengan menambahkan PropWorkshop baru
		judul.PropWorkshop = upload
		if err := db.Db.Save(judul).Error; err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		res.Data = upload
	}

	return c.JSON(http.StatusOK, res)
}

func UploadLpj(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	tema := c.FormValue("judul_workshop")
	nim := c.FormValue("nim")
	name := c.FormValue("name")
	ThnPengurusan := c.FormValue("tahun_pengurusan")

	// Mencari workshop berdasarkan judul
	judul, err := models.GetWorkshopByJudul(db.Db, tema)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	if judul == nil {
		errorMsg := "Judul tidak ditemukan"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusNotFound, res)
	}

	fileDir := "./uploads"
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, os.ModePerm)
	}

	form, err := c.MultipartForm()
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	files := form.File["file"]
	if len(files) == 0 {
		errorMsg := "No files uploaded. Please upload at least one file."
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	var uploads []*models.LpjWorkshop
	for _, file := range files {
		fileUUID := uuid.New().String() // Membuat UUID unik untuk setiap file
		// uuidString := fileUUID.String()
		filename := fileUUID + "_" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "_") // Mengganti spasi dengan garis bawah
		filePath := filepath.Join(fileDir, filename)

		src, err := file.Open()
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		upload := &models.LpjWorkshop{
			WorkshopID:      judul.ID,
			Nim:             nim,
			Name:            name,
			TahunPengurusan: ThnPengurusan,
			PathFile:        filePath,
			NameFile:        filename,
			UploadComplete:  true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		uploads = append(uploads, upload)
	}

	for _, upload := range uploads {
		if err := models.CreateLpj(db.Db, upload); err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		// Mengupdate Workshop dengan menambahkan PropWorkshop baru
		judul.LpjWorkshop = upload
		if err := db.Db.Save(judul).Error; err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		res.Data = upload
	}

	return c.JSON(http.StatusOK, res)
}

func UploadDokumentasi(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	tema := c.FormValue("judul_workshop")
	nim := c.FormValue("nim")
	name := c.FormValue("name")
	ThnPengurusan := c.FormValue("tahun_pengurusan")

	// Mencari workshop berdasarkan judul
	judul, err := models.GetWorkshopByJudul(db.Db, tema)
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusInternalServerError, res)
	}

	if judul == nil {
		errorMsg := "Judul tidak ditemukan"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusNotFound, res)
	}

	fileDir := "./uploads"
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, os.ModePerm)
	}

	form, err := c.MultipartForm()
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	files := form.File["file"]
	if len(files) == 0 {
		errorMsg := "No files uploaded. Please upload at least one file."
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	var uploads []*models.DokumentasiWorkshop
	for _, file := range files {
		fileUUID := uuid.New().String() // Membuat UUID unik untuk setiap file
		// uuidString := fileUUID.String()
		filename := fileUUID + "_" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "_") // Mengganti spasi dengan garis bawah
		filePath := filepath.Join(fileDir, filename)

		src, err := file.Open()
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		upload := &models.DokumentasiWorkshop{
			WorkshopID:      judul.ID,
			Nim:             nim,
			Name:            name,
			TahunPengurusan: ThnPengurusan,
			PathFile:        filePath,
			NameFile:        filename,
			UploadComplete:  true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		uploads = append(uploads, upload)
	}

	for _, upload := range uploads {
		if err := models.CreateDukumentasi(db.Db, upload); err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		// Mengupdate Workshop dengan menambahkan PropWorkshop baru
		judul.Dokumentasi = upload
		if err := db.Db.Save(judul).Error; err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg
			return c.JSON(http.StatusInternalServerError, res)
		}

		res.Data = upload
	}

	return c.JSON(http.StatusOK, res)
}
