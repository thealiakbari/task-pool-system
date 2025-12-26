package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UniversalModel struct {
	Id        uuid.UUID      `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;index"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;index"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

type DBWrapper struct {
	DB *gorm.DB
}

func NewDBWrapper(db *gorm.DB) DBWrapper {
	return DBWrapper{
		DB: db,
	}
}

type Entity interface {
	GetDomain() string
}
