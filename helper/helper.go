package helper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"pbi-final/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SaveFile menyimpan konten file ke lokasi yang ditentukan
func SaveFile(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Gagal menyimpan file: %v", err)
	}
	return nil
}

// ReadFileContents membaca konten dari file multipart
func ReadFileContents(file multipart.File) ([]byte, error) {
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca konten file: %v", err)
	}
	return fileContents, nil
}

// GenerateUUID menghasilkan UUID yang unik
func GenerateUUID() string {
	return uuid.New().String()
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Token disimpan dalam cookie biasanya memiliki format "Bearer [token]"
	// Periksa dan ambil token dari string tersebut
	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 {
		fmt.Println("Invalid token format:", tokenString)
		return nil, errors.New("Invalid token format")
	}

	tokenString = tokenParts[1]
	fmt.Println("Token after split:", tokenString)

	// Parse dan verifikasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token signing method")
		}
		return config.JWT_KEY, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil, err
	}

	// Periksa apakah token valid
	if !token.Valid {
		fmt.Println("Invalid token")
		return nil, errors.New("Invalid token")
	}

	// Ambil klaim dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to parse token claims")
		return nil, errors.New("Failed to parse token claims")
	}

	return claims, nil
}
