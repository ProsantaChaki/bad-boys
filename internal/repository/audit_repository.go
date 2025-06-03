package repository

import (
	"bad_boyes/internal/models"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) CreateLog(auditLog *models.AuditLog) error {
	log.Printf("Creating audit log: %+v", auditLog)

	// Convert JSON fields to strings for storage
	if auditLog.OldValues != nil {
		oldValues, err := json.Marshal(auditLog.OldValues)
		if err != nil {
			log.Printf("Error marshaling old values: %v", err)
			return err
		}
		auditLog.OldValues = models.JSON{"data": string(oldValues)}
	}

	if auditLog.NewValues != nil {
		newValues, err := json.Marshal(auditLog.NewValues)
		if err != nil {
			log.Printf("Error marshaling new values: %v", err)
			return err
		}
		auditLog.NewValues = models.JSON{"data": string(newValues)}
	}

	return r.db.Create(auditLog).Error
}

func (r *AuditRepository) GetLogsByTable(tableName string, recordID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("table_name = ? AND record_id = ?", tableName, recordID).
		Preload("User").
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *AuditRepository) GetLogsByUser(userID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}
