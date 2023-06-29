package controller

import (
	"echoGo/app/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var SECRET = []byte(os.Getenv("SECRET.KEY"))

func GenerateToken(Nim string, tokenType string, exp time.Time) string {
	// user := models.Auth{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nim":        Nim,
		"token_type": tokenType,
		"exp":        exp.Unix(),
	})

	tokenStr, _ := token.SignedString(SECRET)
	return tokenStr
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return SECRET, nil
	})

	claims, oke := token.Claims.(jwt.MapClaims)
	if oke && token.Valid {
		fmt.Println("token_type and exp :", claims["token_type"], claims["exp"])
	} else {
		log.Println("this not valid")
		log.Println(err)
		return nil, err
	}
	return claims, nil
}

func Login(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	req := models.Anggota{}

	if err := c.Bind(&req); err != nil {
		errorMsg := err.Error()
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(400, res)
	}

	Nim, err := models.FindUserByNIM(req.Nim)
	if err != nil {
		errorMsg := "Nim Tidak Terdaftar"
		res.Success = false
		res.Error = &errorMsg
		return c.JSON(http.StatusBadRequest, res)
	}

	accessToken := GenerateToken(Nim.Nim, "access_token", time.Now().Add(15*time.Minute))
	refreshToken := GenerateToken(Nim.Nim, "refresh_token", time.Now().AddDate(0, 0, 5))

	resObj := map[string]interface{}{
		"user": map[string]interface{}{
			"Nim":     Nim.Nim,
			"Role_id": Nim.Role,
		},
		"token": map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}
	res.Success = true
	res.Data = resObj
	return c.JSON(200, res)
}
func RefreshToken(c echo.Context) error {
	res := models.JsonResponse{Success: true}
	req := models.ReqAuth{}
	_ = c.Bind(&req)
	if req.RefreshToken == "" {
		errMsg := "refresh token is required"
		res.Success = false
		res.Error = &errMsg
		c.JSON(400, res)
	}

	claims, err := DecodeToken(req.RefreshToken)
	if err != nil {
		errMsg := err.Error()
		res.Success = false
		res.Error = &errMsg
		c.JSON(401, res)
	}

	nim, found := claims["nim"]
	if !found {
		res.Success = false
		msgErr := "Refresh Token is invalid"
		res.Error = &msgErr
		return c.JSON(401, res)
	}

	if claims["token_type"] != "refresh_token" {
		res.Success = false
		msgErr := "Refresh Token is invalid"
		res.Error = &msgErr
		return c.JSON(401, res)
	}

	res.Data = map[string]string{
		"access_token": GenerateToken(nim.(string), "access_token", time.Now().Add(15*time.Minute)),
	}

	return c.JSON(200, res)
}
