package repository

import (
	"context"
	"go-banking-api/internal/model"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB - Test için hafıza içi SQLite veritabanı oluşturur
//
// NEDEN SQLite?
// - Hafızada çalışır, dosya oluşturmaz
// - Hızlı
// - Test için ideal
// - PostgreSQL kurmaya gerek yok
func setupTestDB(t *testing.T) *gorm.DB {
	// ":memory:" = Hafızada veritabanı oluştur
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Test veritabanı oluşturulamadı: %v", err)
	}

	// Tabloları oluştur (migration)
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("Migration başarısız: %v", err)
	}

	return db
}

// TestUserRepository_Create - Create metodunu test eder
//
// TEST İSİMLENDİRME:
// Test + StructAdı + _ + MetodAdı
func TestUserRepository_Create(t *testing.T) {
	// Arrange (Hazırlık)
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     model.RoleUser,
	}
	user.SetPassword("password123")

	// Act (İşlem)
	err := repo.Create(ctx, user)

	// Assert (Kontrol)
	if err != nil {
		t.Errorf("Create başarısız oldu: %v", err)
	}

	// ID atanmış mı kontrol et
	if user.ID == 0 {
		t.Error("User ID atanmamış")
	}

	t.Logf("✅ Kullanıcı başarıyla oluşturuldu: ID=%d", user.ID)
}

// TestUserRepository_GetByID - GetByID metodunu test eder
func TestUserRepository_GetByID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Önce bir kullanıcı oluştur
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     model.RoleUser,
	}
	user.SetPassword("password123")
	repo.Create(ctx, user)

	// Act - ID ile getir
	foundUser, err := repo.GetByID(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("GetByID başarısız: %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Yanlış kullanıcı geldi. Beklenen: %d, Gelen: %d", user.ID, foundUser.ID)
	}

	if foundUser.Username != user.Username {
		t.Errorf("Username yanlış. Beklenen: %s, Gelen: %s", user.Username, foundUser.Username)
	}

	t.Logf("✅ Kullanıcı başarıyla bulundu: %s", foundUser.Username)
}

// TestUserRepository_GetByID_NotFound - Olmayan kullanıcıyı arama testi
func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Act - Olmayan ID ile getirmeye çalış
	_, err := repo.GetByID(ctx, 999)

	// Assert - Hata almalıyız
	if err == nil {
		t.Error("Olmayan kullanıcı için hata alınmalıydı")
	}

	t.Logf("✅ Beklenen hata alındı: %v", err)
}

// TestUserRepository_GetByUsername - GetByUsername metodunu test eder
func TestUserRepository_GetByUsername(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "johndoe",
		Email:    "john@example.com",
		Role:     model.RoleUser,
	}
	user.SetPassword("password123")
	repo.Create(ctx, user)

	// Act
	foundUser, err := repo.GetByUsername(ctx, "johndoe")

	// Assert
	if err != nil {
		t.Fatalf("GetByUsername başarısız: %v", err)
	}

	if foundUser.Username != "johndoe" {
		t.Errorf("Yanlış kullanıcı. Beklenen: johndoe, Gelen: %s", foundUser.Username)
	}

	t.Logf("✅ Username ile kullanıcı bulundu: %s", foundUser.Email)
}

// TestUserRepository_GetByEmail - GetByEmail metodunu test eder
func TestUserRepository_GetByEmail(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "janedoe",
		Email:    "jane@example.com",
		Role:     model.RoleUser,
	}
	user.SetPassword("password123")
	repo.Create(ctx, user)

	// Act
	foundUser, err := repo.GetByEmail(ctx, "jane@example.com")

	// Assert
	if err != nil {
		t.Fatalf("GetByEmail başarısız: %v", err)
	}

	if foundUser.Email != "jane@example.com" {
		t.Errorf("Yanlış kullanıcı. Beklenen: jane@example.com, Gelen: %s", foundUser.Email)
	}

	t.Logf("✅ Email ile kullanıcı bulundu: %s", foundUser.Username)
}

// TestUserRepository_GetAll - GetAll metodunu test eder
func TestUserRepository_GetAll(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 3 kullanıcı oluştur
	users := []*model.User{
		{Username: "user1", Email: "user1@example.com", Role: model.RoleUser},
		{Username: "user2", Email: "user2@example.com", Role: model.RoleUser},
		{Username: "user3", Email: "user3@example.com", Role: model.RoleAdmin},
	}

	for _, u := range users {
		u.SetPassword("password123")
		repo.Create(ctx, u)
	}

	// Act - İlk 10 kullanıcıyı getir
	allUsers, err := repo.GetAll(ctx, 10, 0)

	// Assert
	if err != nil {
		t.Fatalf("GetAll başarısız: %v", err)
	}

	if len(allUsers) != 3 {
		t.Errorf("Yanlış kullanıcı sayısı. Beklenen: 3, Gelen: %d", len(allUsers))
	}

	t.Logf("✅ Toplam %d kullanıcı listelendi", len(allUsers))
}

// TestUserRepository_Update - Update metodunu test eder
func TestUserRepository_Update(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "oldname",
		Email:    "old@example.com",
		Role:     model.RoleUser,
	}
	user.SetPassword("password123")
	repo.Create(ctx, user)

	// Act - Email değiştir
	user.Email = "new@example.com"
	err := repo.Update(ctx, user)

	// Assert
	if err != nil {
		t.Fatalf("Update başarısız: %v", err)
	}

	// Değişikliği kontrol et
	updatedUser, _ := repo.GetByID(ctx, user.ID)
	if updatedUser.Email != "new@example.com" {
		t.Errorf("Email güncellenmedi. Beklenen: new@example.com, Gelen: %s", updatedUser.Email)
	}

	t.Logf("✅ Kullanıcı başarıyla güncellendi: %s", updatedUser.Email)
}

// TestUserRepository_Delete - Delete metodunu test eder
func TestUserRepository_Delete(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username: "tobedeleted",
		Email:    "delete@example.com",
		Role:     model.RoleUser,
	}
	user.SetPassword("password123")
	repo.Create(ctx, user)

	// Act - Kullanıcıyı sil
	err := repo.Delete(ctx, user.ID)

	// Assert
	if err != nil {
		t.Fatalf("Delete başarısız: %v", err)
	}

	// Silindiğini kontrol et (soft delete, bulunamaz olmalı)
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Error("Silinen kullanıcı hala bulunuyor")
	}

	t.Logf("✅ Kullanıcı başarıyla silindi (soft delete)")
}
