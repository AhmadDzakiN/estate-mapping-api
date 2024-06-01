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
	Width     int
	Length    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tree struct {
	ID                 string
	EstateID           string
	HorizontalPosition int
	VerticalPosition   int
	Height             int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
