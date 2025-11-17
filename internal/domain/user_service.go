package domain

import (
	"context"
	"go-banking-api/internal/model"
)

// UserService - Kullanıcı iş mantığı için interface
//
// SERVICE vs REPOSITORY Farkı:
//
// Repository = Sadece veritabanı işlemleri (CRUD)
//   - Create, Read, Update, Delete
//   - Basit, direkt işlemler
//
// Service = İş mantığı (Business Logic)
//   - Register = Create + şifre hash + email kontrolü + validation
//   - Login = şifre kontrolü + kullanıcı bulma
//   - Karmaşık işlemler, birden fazla repository kullanabilir
type UserService interface {
	// Register - Yeni kullanıcı kaydeder
	// İş mantığı:
	//   1. Email ve username'in benzersiz olduğunu kontrol et
	//   2. Şifreyi hash'le
	//   3. Kullanıcıyı oluştur
	//   4. Balance kaydı oluştur (ileride)
	Register(ctx context.Context, username, email, password string) (*model.User, error)

	// Login - Kullanıcı girişi yapar
	// İş mantığı:
	//   1. Kullanıcıyı username veya email ile bul
	//   2. Şifreyi kontrol et
	//   3. JWT token oluştur (ileride)
	Login(ctx context.Context, usernameOrEmail, password string) (*model.User, error)

	// GetUserByID - ID'ye göre kullanıcı getirir
	GetUserByID(ctx context.Context, id uint) (*model.User, error)

	// GetUserByUsername - Username'e göre kullanıcı getirir
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetAllUsers - Tüm kullanıcıları listeler (pagination ile)
	GetAllUsers(ctx context.Context, page, pageSize int) ([]*model.User, error)

	// UpdateUser - Kullanıcı bilgilerini günceller
	// İş mantığı:
	//   1. Yeni email/username benzersiz mi kontrol et
	//   2. Güncelle
	UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) (*model.User, error)

	// DeleteUser - Kullanıcıyı siler
	// İş mantığı:
	//   1. Kullanıcının bakiyesi 0 mı kontrol et (ileride)
	//   2. İlişkili kayıtları sil/güncelle
	//   3. Kullanıcıyı sil
	DeleteUser(ctx context.Context, id uint) error

	// ChangePassword - Kullanıcı şifresini değiştirir
	// İş mantığı:
	//   1. Eski şifreyi doğrula
	//   2. Yeni şifreyi hash'le
	//   3. Güncelle
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}
