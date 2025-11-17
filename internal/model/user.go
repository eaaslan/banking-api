package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Role         string `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}

func (u *User) Validate() error {
	if err := u.validateUsername(); err != nil {
		return err
	}

	if err := u.validateEmail(); err != nil {
		return err
	}

	if err := u.validatePassword(); err != nil {
		return err
	}

	if err := u.validateRole(); err != nil {
		return err
	}

	return nil
}

func (u *User) validateUsername() error {
	u.Username = strings.TrimSpace(u.Username)

	if u.Username == "" {
		return errors.New("kullanici adi bos olamaz")
	}

	if len(u.Username) < 3 {
		return errors.New("kullanici adi en az 3 karakter olmalidir")
	}

	if len(u.Username) > 50 {
		return errors.New("kullanici adi en fazla 50 karakter olabilir")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, u.Username)
	if !matched {
		return errors.New("kullanici adi sadece harf, rakam, alt cizgi ve tire icermelidir")
	}

	return nil
}

func (u *User) validateEmail() error {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	if u.Email == "" {
		return errors.New("email bos olamaz")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, u.Email)
	if !matched {
		return errors.New("gecersiz email formati")
	}

	if len(u.Email) > 100 {
		return errors.New("email en fazla 100 karakter olabilir")
	}

	return nil
}

func (u *User) validatePassword() error {
	if u.ID == 0 && u.PasswordHash == "" {
		return errors.New("sifre bos olamaz")
	}

	if u.PasswordHash != "" && len(u.PasswordHash) < 6 {
		return errors.New("sifre en az 6 karakter olmalidir")
	}

	return nil
}

func (u *User) validateRole() error {
	u.Role = strings.TrimSpace(strings.ToLower(u.Role))

	if u.Role == "" {
		u.Role = RoleUser
		return nil
	}

	if u.Role != RoleAdmin && u.Role != RoleUser {
		return errors.New("role sadece 'admin' veya 'user' olabilir")
	}

	return nil
}

func (u *User) SetPassword(plainPassword string) error {
	if len(plainPassword) < 6 {
		return errors.New("sifre en az 6 karakter olmalidir")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)

	return nil
}

func (User) TableName() string {
	return "users"
}
