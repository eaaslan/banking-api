package repository

import (
	"context"
	"errors"
	"fmt"
	"go-banking-api/internal/domain"
	"go-banking-api/internal/model"
	"gorm.io/gorm"
)

// userRepository - GORM kullanarak UserRepository interface'ini uygular
//
// STRUCT NEDİR?
// - Go'da class yoktur, struct vardır
// - struct = veri tutar
// - struct'a method ekleyerek "object" gibi davranır
type userRepository struct {
	db *gorm.DB // Veritabanı bağlantısı
}

// NewUserRepository - CONSTRUCTOR (Yapıcı) fonksiyon
//
// NEDEN Constructor?
// - Repository oluştururken db bağlantısını veriyoruz
// - Dışarıdan interface döner, içeride implementation kullanılır
// - Bu "Dependency Injection" pattern'idir
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create - Yeni kullanıcı oluşturur
//
// (r *userRepository) = Bu method userRepository struct'ına aittir
// r = receiver, "this" ya da "self" gibi
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// WithContext = timeout kontrolü için context kullan
	// Create = GORM'un INSERT komutu
	result := r.db.WithContext(ctx).Create(user)

	// Hata kontrolü
	if result.Error != nil {
		// %w = error wrapping, hatayı sarmalayarak döner
		return fmt.Errorf("kullanıcı oluşturulamadı: %w", result.Error)
	}

	return nil
}

// GetByID - ID'ye göre kullanıcı bulur
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	// First = İlk kaydı getir
	// Eğer bulamazsa gorm.ErrRecordNotFound hatası döner
	result := r.db.WithContext(ctx).First(&user, id)

	// Özel hata kontrolü: Kayıt bulunamadı mı?
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("kullanıcı bulunamadı: %d", id)
	}

	// Genel hata kontrolü
	if result.Error != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", result.Error)
	}

	return &user, nil
}

// GetByUsername - Kullanıcı adına göre bulur
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	// Where = SQL WHERE koşulu
	// ? = placeholder, SQL injection'dan korur
	result := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("kullanıcı bulunamadı: %s", username)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", result.Error)
	}

	return &user, nil
}

// GetByEmail - Email'e göre bulur
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	result := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("kullanıcı bulunamadı: %s", email)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("kullanıcı getirilemedi: %w", result.Error)
	}

	return &user, nil
}

// GetAll - Tüm kullanıcıları getirir (sayfalama ile)
//
// PAGINATION NEDİR?
// - Tüm kayıtları getirmek yerine sayfa sayfa getirmek
// - limit: Her sayfada kaç kayıt (örn: 10)
// - offset: Kaç kayıt atla (örn: sayfa 2 için 10 atla)
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*model.User, error) {
	var users []*model.User

	// Varsayılan değerler
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Order = Sıralama
	// Limit = Maksimum kayıt sayısı
	// Offset = Kaç kayıt atla
	result := r.db.WithContext(ctx).
		Order("created_at DESC"). // En yeniden eskiye
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("kullanıcılar getirilemedi: %w", result.Error)
	}

	return users, nil
}

// Update - Kullanıcı bilgilerini günceller
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Save = GORM'un UPDATE komutu
	// Tüm alanları günceller
	result := r.db.WithContext(ctx).Save(user)

	if result.Error != nil {
		return fmt.Errorf("kullanıcı güncellenemedi: %w", result.Error)
	}

	// RowsAffected = Kaç satır etkilendi?
	// 0 ise kullanıcı bulunamadı demektir
	if result.RowsAffected == 0 {
		return fmt.Errorf("kullanıcı bulunamadı veya değişiklik yok")
	}

	return nil
}

// Delete - Kullanıcıyı siler (SOFT DELETE)
//
// SOFT DELETE NEDİR?
// - Kaydı gerçekten silmez
// - deleted_at alanına şu anki zamanı yazar
// - Sorgulamalarda otomatik olarak deleted_at NULL olanlar gelir
// - Veri kaybı olmaz, geri getirilebilir
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	// Delete = GORM'un soft delete komutu
	result := r.db.WithContext(ctx).Delete(&model.User{}, id)

	if result.Error != nil {
		return fmt.Errorf("kullanıcı silinemedi: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("kullanıcı bulunamadı: %d", id)
	}

	return nil
}
