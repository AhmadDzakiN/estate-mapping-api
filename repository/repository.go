// This file contains the repository implementation layer.
package repository

import (
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

//type NewRepositoryOptions struct {
//	Dsn string
//}

func NewRepository(opts Repository) *Repository {
	//db, err := sql.Open("postgres", opts.Dsn)
	//if err != nil {
	//	panic(err)
	//}
	//return &Repository{
	//	Db: db,
	//}

	return &Repository{
		Db: opts.Db,
	}
}
