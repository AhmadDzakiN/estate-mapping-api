// This file contains types that are used in the repository layer.
package repository

import "time"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type Estate struct {
	ID        string    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	Width     int       `gorm:"column:width;not null"`
	Length    int       `gorm:"column:length;not null"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

type Tree struct {
	ID                 string    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	EstateID           string    `gorm:"column:estate_id;type:uuid;not null"`
	HorizontalPosition int       `gorm:"column:horizontal_position;not null"`
	VerticalPosition   int       `gorm:"column:vertical_position;not null"`
	Height             int       `gorm:"column:height;not null"`
	CreatedAt          time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt          time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}
