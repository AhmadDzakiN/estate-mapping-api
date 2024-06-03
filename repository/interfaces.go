// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, newEstate *Estate) (err error)
	GetEstateByID(ctx context.Context, estateID string) (estate Estate, err error)
	CreateTree(ctx context.Context, newTree *Tree) (err error)
	GetTreeHeightsByEstateID(ctx context.Context, estateID string) (treeHeights []int, err error)
}
