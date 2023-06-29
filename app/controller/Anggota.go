package controller

import (
	"echoGo/app/models"
	"echoGo/db"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	req := models.Anggota{}

	if err := c.Bind(&req); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	_, err := models.FindUserByNIM(req.Nim)
	if err == nil {
		errorMsg := "Nim Sudah digunakan"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	Create := &models.Anggota{
		Nim:      req.Nim,
		Password: req.Password,
		RoleID:   req.RoleID,
		AnggotaDetails: &models.AnggotaDetails{
			Name:            req.AnggotaDetails.Name,
			TahunPengurusan: req.AnggotaDetails.TahunPengurusan,
			Divisi:          req.AnggotaDetails.Divisi,
			Email:           req.AnggotaDetails.Email,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.Db.Create(&Create).Error
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(500, res)
	}
	// Fetch user's details and associated role
	err = db.Db.Model(&Create).Preload("Role").Preload("AnggotaDetails").First(&Create).Error
	if err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(500, res)
	}
	res.Data = map[string]interface{}{
		"Nim":             Create.Nim,
		"Name":            Create.AnggotaDetails.Name,
		"TahunPengurusan": Create.AnggotaDetails.TahunPengurusan,
		"Divisi":          Create.AnggotaDetails.Divisi,
		"Email":           Create.AnggotaDetails.Email,
	}

	res.Success = true
	return c.JSON(200, res)
}
