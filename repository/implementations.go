package repository

import (
	"context"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateEstate(ctx context.Context, length, width int) (newEstateID string, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO estates (width, length) VALUES ($1, $2) RETURNING id", width, length).Scan(&newEstateID)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetEstate(ctx context.Context, estateID string) (estate Estate, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, width, length FROM estates where id = $1", estateID).Scan(&estate.ID, &estate.Width, &estate.Length)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateTree(ctx context.Context, estateID string, horizontal, vertical, height int) (newTreeID string, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO trees (estate_id, horizontal_position, vertical_position, height) VALUES ($1, $2, $3, $4) RETURNING id", estateID, horizontal, vertical, height).
		Scan(&newTreeID)
	if err != nil {
		return
	}
	return
}
