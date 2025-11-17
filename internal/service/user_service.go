package service

import (
	"context"
	"errors"
	"fmt"
	"go-banking-api/internal/domain"
	"go-banking-api/internal/model"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))

	if username == "" {
		return nil, errors.New("kullanıcı adı boş olamaz")
	}

	if email == "" {
		return nil, errors.New("email boş olamaz")
	}

	if len(password) < 6 {
		return nil, errors.New("şifre en az 6 karakter olmalıdır")
	}

	existingUser, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, fmt.Errorf("kullanıcı adı zaten kullanılıyor: %s", username)
	}

	existingUser, _ = s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, fmt.Errorf("email zaten kullanılıyor: %s", email)
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Role:     model.RoleUser,
	}

	err := user.SetPassword(password)
	if err != nil {
		return nil, fmt.Errorf("şifre hash'lenemedi: %w", err)
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı kaydedilemedi: %w", err)
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, usernameOrEmail, password string) (*model.User, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)

	if usernameOrEmail == "" {
		return nil, errors.New("kullanıcı adı veya email boş olamaz")
	}

	if password == "" {
		return nil, errors.New("şifre boş olamaz")
	}

	var user *model.User
	var err error

	if strings.Contains(usernameOrEmail, "@") {
		user, err = s.userRepo.GetByEmail(ctx, usernameOrEmail)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, usernameOrEmail)
	}

	if err != nil {
		return nil, errors.New("kullanıcı adı veya şifre hatalı")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("kullanıcı adı veya şifre hatalı")
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("geçersiz kullanıcı ID")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	username = strings.TrimSpace(username)

	if username == "" {
		return nil, errors.New("kullanıcı adı boş olamaz")
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context, page, pageSize int) ([]*model.User, error) {
	// Varsayılan değerler
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	// Maksimum sayfa boyutu 100
	if pageSize > 100 {
		pageSize = 100
	}

	// Offset hesapla
	offset := (page - 1) * pageSize

	users, err := s.userRepo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("kullanıcılar getirilemedi: %w", err)
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) (*model.User, error) {
	// 1. Kullanıcıyı bul
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	for key, value := range updates {
		switch key {
		case "username":
			newUsername := value.(string)
			newUsername = strings.TrimSpace(newUsername)

			existing, _ := s.userRepo.GetByUsername(ctx, newUsername)
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("kullanıcı adı zaten kullanılıyor: %s", newUsername)
			}

			user.Username = newUsername

		case "email":
			newEmail := value.(string)
			newEmail = strings.TrimSpace(strings.ToLower(newEmail))

			existing, _ := s.userRepo.GetByEmail(ctx, newEmail)
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("email zaten kullanılıyor: %s", newEmail)
			}

			user.Email = newEmail

		case "role":
			newRole := value.(string)
			if newRole != model.RoleAdmin && newRole != model.RoleUser {
				return nil, errors.New("geçersiz rol")
			}
			user.Role = newRole
		}
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("kullanıcı güncellenemedi: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("geçersiz kullanıcı ID")
	}

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("kullanıcı silinemedi: %w", err)
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// 1. Validasyon
	if len(newPassword) < 6 {
		return errors.New("yeni şifre en az 6 karakter olmalıdır")
	}

	// 2. Kullanıcıyı bul
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.New("eski şifre hatalı")
	}

	err = user.SetPassword(newPassword)
	if err != nil {
		return fmt.Errorf("şifre hash'lenemedi: %w", err)
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("şifre değiştirilemedi: %w", err)
	}

	return nil
}
