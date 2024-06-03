package repository

import (
	"context"
	"errors"
)

func (r *Repository) CreateEstate(ctx context.Context, newEstate *Estate) (err error) {
	result := r.Db.WithContext(ctx).Create(newEstate)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("Insert operation failed because rows affected is 0")
		return
	}

	return
}

func (r *Repository) GetEstateByID(ctx context.Context, estateID string) (estate Estate, err error) {
	result := r.Db.WithContext(ctx).Select("id", "width", "length").Where("id", estateID).First(&estate)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (r *Repository) CreateTree(ctx context.Context, newTree *Tree) (err error) {
	result := r.Db.WithContext(ctx).Create(newTree)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected < 1 {
		err = errors.New("Insert operation failed because rows affected is 0")
		return
	}

	return
}

func (r *Repository) GetTreeHeightsByEstateID(ctx context.Context, estateID string) (treeHeights []int, err error) {
	result := r.Db.WithContext(ctx).Table("trees").Select("height").
		Where("estate_id", estateID).Order("height asc").Find(&treeHeights)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}
