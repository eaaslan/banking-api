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

func (r *PostgresRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, role = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Role, user.ID)
	return err
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// --- Transaction Repository ---

func (r *PostgresRepository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	query := `INSERT INTO transactions (from_user_id, to_user_id, amount, type, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, tx.FromUserID, tx.ToUserID, tx.Amount, tx.Type, tx.Status).Scan(&tx.ID, &tx.CreatedAt)
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

func (r *PostgresRepository) GetTransactionsByUserID(ctx context.Context, userID int64) ([]*models.Transaction, error) {
	query := `SELECT id, from_user_id, to_user_id, amount, type, status, created_at FROM transactions WHERE from_user_id = $1 OR to_user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
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

func (r *PostgresRepository) GetAuditLogsByEntity(ctx context.Context, entityType string, entityID int64) ([]*models.AuditLog, error) {
   query := `SELECT id, entity_type, entity_id, action, details, created_at FROM audit_logs WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		l := &models.AuditLog{}
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.Details, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
