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
	ID        string
	Width     uint
	Length    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tree struct {
	ID                 string
	EstateID           string
	HorizontalPosition uint
	VerticalPosition   uint
	Height             uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
