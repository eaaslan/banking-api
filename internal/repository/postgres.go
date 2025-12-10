package repository

import (
	"context"
	"database/sql"
	"backend/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// --- User Repository ---

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)
	return err
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// --- Transaction Repository ---

func (r *PostgresRepository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	query := `INSERT INTO transactions (from_user_id, to_user_id, amount, type, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, tx.FromUserID, tx.ToUserID, tx.Amount, tx.Type, tx.Status).Scan(&tx.ID)
	return err
}

func (r *PostgresRepository) GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `SELECT id, from_user_id, to_user_id, amount, type, status, created_at FROM transactions WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *PostgresRepository) UpdateTransactionStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// --- Balance Repository ---

func (r *PostgresRepository) GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error) {
	b := &models.Balance{}
	query := `SELECT user_id, amount, last_updated_at FROM balances WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&b.UserID, &b.Amount, &b.LastUpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *PostgresRepository) CreateBalance(ctx context.Context, balance *models.Balance) error {
	query := `INSERT INTO balances (user_id, amount) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, balance.UserID, balance.Amount)
	return err
}

func (r *PostgresRepository) UpdateBalance(ctx context.Context, balance *models.Balance) error {
	query := `UPDATE balances SET amount = $1, last_updated_at = CURRENT_TIMESTAMP WHERE user_id = $2`
	_, err := r.db.ExecContext(ctx, query, balance.Amount, balance.UserID)
	return err
}

// --- Audit Repository ---

func (r *PostgresRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	query := `INSERT INTO audit_logs (entity_type, entity_id, action, details) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, log.EntityType, log.EntityID, log.Action, log.Details)
	return err
}
