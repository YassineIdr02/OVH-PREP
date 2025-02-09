package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/YassineIdr02/ovh-prep/E2E-Tests/models"
	"github.com/YassineIdr02/ovh-prep/E2E-Tests/storage"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secretkey")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 3).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	if err := storage.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	fmt.Println("User found:", user.Username)
	fmt.Println("Stored password hash:", user.Password)

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	fmt.Println("Stored password hash:", string(bytes))
	if err != nil {
		if string(bytes) != user.Password {
			fmt.Println("Invalid credentials")
			return nil, errors.New("invalid credentials")
		}
	}
	return &user, nil
}
