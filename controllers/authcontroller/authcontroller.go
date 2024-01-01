package authcontroller

import (
	"encoding/json"
	"net/http"
	"pbi-final/config"
	"pbi-final/helper"
	"pbi-final/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// cek email
	var user models.User
	if err := models.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Username dan pw salah"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Username dan pw salah"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// proses jwt
	expTime := time.Now().Add(time.Minute * 10)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// deklarasi algoritma login
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "Bearer " + token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login Sukses"}
	helper.ResponseJSON(w, http.StatusOK, response)

}

func Register(w http.ResponseWriter, r *http.Request) {

	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	//hash paswwrod
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "sukses"}
	helper.ResponseJSON(w, http.StatusOK, response)

}
func Logout(w http.ResponseWriter, r *http.Request) {
	// hapus token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "LOGOUT BERHASIL"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func UploadPhoto(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		response := map[string]string{"message": "Token tidak valid"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	tokenValue := cookie.Value
	claims, err := helper.ParseToken(tokenValue)
	if err != nil {
		response := map[string]string{"message": "Token tidak valid"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Membaca id pengguna dari klaim
	userClaim, ok := claims["user_id"].(float64)
	if !ok {
		response := map[string]string{"message": "Token tidak valid"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Membaca id pengguna dari klaim
	userID := uint(userClaim)

	// Menerima file foto dari form
	file, handler, err := r.FormFile("photo")
	if err != nil {
		response := map[string]string{"message": "Gagal menerima file"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	// Baca konten file
	fileContents, err := helper.ReadFileContents(file)
	if err != nil {
		response := map[string]string{"message": "Gagal membaca konten file"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Simpan file foto di server atau tempat penyimpanan yang diinginkan
	photoPath := "uploads/" + helper.GenerateUUID() + handler.Filename
	err = helper.SaveFile(photoPath, fileContents)
	if err != nil {
		response := map[string]string{"message": "Gagal menyimpan file"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Simpan informasi foto di database
	photo := models.Photo{
		Title:    r.FormValue("title"),
		Caption:  r.FormValue("caption"),
		PhotoURL: photoPath,
		UserID:   userID,
	}

	if err := models.DB.Create(&photo).Error; err != nil {
		response := map[string]string{"message": "Gagal menyimpan informasi foto"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Foto berhasil diunggah"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
