package repository

import "context"

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
