package model

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	EntityType string `gorm:"type:varchar(50);not null;index" json:"entity_type"` // "user", "transaction", "balance"
	EntityID   uint   `gorm:"not null;index" json:"entity_id"`                    // Related entity's ID
	Action     string `gorm:"type:varchar(20);not null;index" json:"action"`      // "create", "update", "delete"
	Details    string `gorm:"type:text" json:"details"`                           // JSON formatted details
}

// Entity types
const (
	EntityTypeUser        = "user"
	EntityTypeTransaction = "transaction"
	EntityTypeBalance     = "balance"
)

// Action types
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionRead   = "read"
)

// BeforeCreate is a GORM hook that runs before creating a record
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	return a.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (a *AuditLog) BeforeUpdate(tx *gorm.DB) error {
	return a.Validate()
}

// Validate validates all audit log fields
func (a *AuditLog) Validate() error {
	// Validate entity type
	if err := a.validateEntityType(); err != nil {
		return err
	}

	// Validate entity ID
	if err := a.validateEntityID(); err != nil {
		return err
	}

	// Validate action
	if err := a.validateAction(); err != nil {
		return err
	}

	// Validate details
	if err := a.validateDetails(); err != nil {
		return err
	}

	return nil
}

// validateEntityType validates entity type
func (a *AuditLog) validateEntityType() error {
	a.EntityType = strings.TrimSpace(strings.ToLower(a.EntityType))

	if a.EntityType == "" {
		return errors.New("entity_type cannot be empty")
	}

	validTypes := []string{EntityTypeUser, EntityTypeTransaction, EntityTypeBalance}
	for _, vt := range validTypes {
		if a.EntityType == vt {
			return nil
		}
	}

	return errors.New("entity_type must be 'user', 'transaction', or 'balance'")
}

// validateEntityID validates entity ID
func (a *AuditLog) validateEntityID() error {
	if a.EntityID == 0 {
		return errors.New("entity_id cannot be empty or zero")
	}

	return nil
}

// validateAction validates action type
func (a *AuditLog) validateAction() error {
	a.Action = strings.TrimSpace(strings.ToLower(a.Action))

	if a.Action == "" {
		return errors.New("action cannot be empty")
	}

	validActions := []string{ActionCreate, ActionUpdate, ActionDelete, ActionRead}
	for _, va := range validActions {
		if a.Action == va {
			return nil
		}
	}

	return errors.New("action must be 'create', 'update', 'delete', or 'read'")
}

// validateDetails validates details field
func (a *AuditLog) validateDetails() error {
	// Details can be empty, but if provided should not exceed limit
	if len(a.Details) > 5000 {
		return errors.New("details cannot exceed 5000 characters")
	}

	return nil
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(entityType string, entityID uint, action string, details string) *AuditLog {
	return &AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Details:    details,
	}
}

// IsCreate checks if action is create
func (a *AuditLog) IsCreate() bool {
	return a.Action == ActionCreate
}

// IsUpdate checks if action is update
func (a *AuditLog) IsUpdate() bool {
	return a.Action == ActionUpdate
}

// IsDelete checks if action is delete
func (a *AuditLog) IsDelete() bool {
	return a.Action == ActionDelete
}

// IsRead checks if action is read
func (a *AuditLog) IsRead() bool {
	return a.Action == ActionRead
}

// TableName specifies the table name for GORM
func (AuditLog) TableName() string {
	return "audit_logs"
}
