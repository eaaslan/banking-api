package service

import (
	"context"
	"fmt"
	"go-banking-api/internal/model"
	"go-banking-api/internal/repository"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestService - Test için service ve repository oluşturur
func setupTestService(t *testing.T) *userService {
	// SQLite hafıza veritabanı
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Test veritabanı oluşturulamadı: %v", err)
	}

	// Migration
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("Migration başarısız: %v", err)
	}

	// Repository oluştur
	userRepo := repository.NewUserRepository(db)

	// Service oluştur (userService döner çünkü test içinde private metodlara erişmek isteyebiliriz)
	return &userService{
		userRepo: userRepo,
	}
}

// TestUserService_Register - Register metodunu test eder
func TestUserService_Register(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Act - Yeni kullanıcı kaydet
	user, err := service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Assert
	if err != nil {
		t.Fatalf("Register başarısız: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID atanmamış")
	}

	if user.Username != "johndoe" {
		t.Errorf("Username yanlış. Beklenen: johndoe, Gelen: %s", user.Username)
	}

	if user.Email != "john@example.com" {
		t.Errorf("Email yanlış. Beklenen: john@example.com, Gelen: %s", user.Email)
	}

	if user.Role != model.RoleUser {
		t.Errorf("Role yanlış. Beklenen: %s, Gelen: %s", model.RoleUser, user.Role)
	}

	t.Logf("✅ Kullanıcı başarıyla kaydedildi: %s", user.Username)
}

// TestUserService_Register_DuplicateUsername - Aynı username ile kayıt testi
func TestUserService_Register_DuplicateUsername(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// İlk kullanıcıyı kaydet
	service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Aynı username ile tekrar kayıt dene
	_, err := service.Register(ctx, "johndoe", "different@example.com", "password123")

	// Assert - Hata almalıyız
	if err == nil {
		t.Error("Aynı username ile kayıt olunmamalıydı")
	}

	t.Logf("✅ Beklenen hata alındı: %v", err)
}

// TestUserService_Register_DuplicateEmail - Aynı email ile kayıt testi
func TestUserService_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// İlk kullanıcıyı kaydet
	service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Aynı email ile tekrar kayıt dene
	_, err := service.Register(ctx, "differentuser", "john@example.com", "password123")

	// Assert - Hata almalıyız
	if err == nil {
		t.Error("Aynı email ile kayıt olunmamalıydı")
	}

	t.Logf("✅ Beklenen hata alındı: %v", err)
}

// TestUserService_Register_ShortPassword - Kısa şifre testi
func TestUserService_Register_ShortPassword(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Act - Kısa şifre ile kayıt dene
	_, err := service.Register(ctx, "johndoe", "john@example.com", "123")

	// Assert - Hata almalıyız
	if err == nil {
		t.Error("Kısa şifre ile kayıt olunmamalıydı")
	}

	t.Logf("✅ Beklenen hata alındı: %v", err)
}

// TestUserService_Login - Login metodunu test eder
func TestUserService_Login(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Önce kullanıcı kaydet
	service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Username ile login
	user, err := service.Login(ctx, "johndoe", "password123")

	// Assert
	if err != nil {
		t.Fatalf("Login başarısız: %v", err)
	}

	if user.Username != "johndoe" {
		t.Errorf("Yanlış kullanıcı. Beklenen: johndoe, Gelen: %s", user.Username)
	}

	t.Logf("✅ Username ile login başarılı: %s", user.Username)
}

// TestUserService_Login_WithEmail - Email ile login testi
func TestUserService_Login_WithEmail(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Önce kullanıcı kaydet
	service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Email ile login
	user, err := service.Login(ctx, "john@example.com", "password123")

	// Assert
	if err != nil {
		t.Fatalf("Login başarısız: %v", err)
	}

	if user.Email != "john@example.com" {
		t.Errorf("Yanlış kullanıcı. Beklenen: john@example.com, Gelen: %s", user.Email)
	}

	t.Logf("✅ Email ile login başarılı: %s", user.Email)
}

// TestUserService_Login_WrongPassword - Yanlış şifre ile login testi
func TestUserService_Login_WrongPassword(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Önce kullanıcı kaydet
	service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Yanlış şifre ile login dene
	_, err := service.Login(ctx, "johndoe", "wrongpassword")

	// Assert - Hata almalıyız
	if err == nil {
		t.Error("Yanlış şifre ile login olunmamalıydı")
	}

	t.Logf("✅ Beklenen hata alındı: %v", err)
}

// TestUserService_GetAllUsers - Pagination testi
func TestUserService_GetAllUsers(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// 5 kullanıcı kaydet
	for i := 1; i <= 5; i++ {
		username := fmt.Sprintf("user%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		service.Register(ctx, username, email, "password123")
	}

	// Act - İlk 3 kullanıcıyı getir (sayfa 1, sayfa boyutu 3)
	users, err := service.GetAllUsers(ctx, 1, 3)

	// Assert
	if err != nil {
		t.Fatalf("GetAllUsers başarısız: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Yanlış kullanıcı sayısı. Beklenen: 3, Gelen: %d", len(users))
	}

	t.Logf("✅ Sayfa 1: %d kullanıcı listelendi", len(users))

	// Act - İkinci sayfayı getir
	users, err = service.GetAllUsers(ctx, 2, 3)

	// Assert
	if err != nil {
		t.Fatalf("GetAllUsers başarısız: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Yanlış kullanıcı sayısı. Beklenen: 2, Gelen: %d", len(users))
	}

	t.Logf("✅ Sayfa 2: %d kullanıcı listelendi", len(users))
}

// TestUserService_UpdateUser - Kullanıcı güncelleme testi
func TestUserService_UpdateUser(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Kullanıcı kaydet
	user, _ := service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Email güncelle
	updates := map[string]interface{}{
		"email": "newemail@example.com",
	}
	updatedUser, err := service.UpdateUser(ctx, user.ID, updates)

	// Assert
	if err != nil {
		t.Fatalf("UpdateUser başarısız: %v", err)
	}

	if updatedUser.Email != "newemail@example.com" {
		t.Errorf("Email güncellenmedi. Beklenen: newemail@example.com, Gelen: %s", updatedUser.Email)
	}

	t.Logf("✅ Email başarıyla güncellendi: %s", updatedUser.Email)
}

// TestUserService_ChangePassword - Şifre değiştirme testi
func TestUserService_ChangePassword(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Kullanıcı kaydet
	user, _ := service.Register(ctx, "johndoe", "john@example.com", "oldpassword")

	// Act - Şifre değiştir
	err := service.ChangePassword(ctx, user.ID, "oldpassword", "newpassword123")

	// Assert
	if err != nil {
		t.Fatalf("ChangePassword başarısız: %v", err)
	}

	// Yeni şifre ile login dene
	_, err = service.Login(ctx, "johndoe", "newpassword123")
	if err != nil {
		t.Error("Yeni şifre ile login başarısız")
	}

	t.Logf("✅ Şifre başarıyla değiştirildi")
}

// TestUserService_ChangePassword_WrongOldPassword - Yanlış eski şifre testi
func TestUserService_ChangePassword_WrongOldPassword(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Kullanıcı kaydet
	user, _ := service.Register(ctx, "johndoe", "john@example.com", "oldpassword")

	// Act - Yanlış eski şifre ile değiştirmeye çalış
	err := service.ChangePassword(ctx, user.ID, "wrongoldpassword", "newpassword123")

	// Assert - Hata almalıyız
	if err == nil {
		t.Error("Yanlış eski şifre ile şifre değiştirilmemeli")
	}

	t.Logf("✅ Beklenen hata alındı: %v", err)
}

// TestUserService_DeleteUser - Kullanıcı silme testi
func TestUserService_DeleteUser(t *testing.T) {
	// Arrange
	service := setupTestService(t)
	ctx := context.Background()

	// Kullanıcı kaydet
	user, _ := service.Register(ctx, "johndoe", "john@example.com", "password123")

	// Act - Kullanıcıyı sil
	err := service.DeleteUser(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("DeleteUser başarısız: %v", err)
	}

	// Silindiğini kontrol et
	_, err = service.GetUserByID(ctx, user.ID)
	if err == nil {
		t.Error("Silinen kullanıcı hala bulunuyor")
	}

	t.Logf("✅ Kullanıcı başarıyla silindi")
}
