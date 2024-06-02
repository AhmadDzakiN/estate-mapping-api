package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type RepositoryTestSuite struct {
	suite.Suite
	sqlMock    sqlmock.Sqlmock
	repository RepositoryInterface
	ctx        context.Context
	loc        *time.Location
}

func (r *RepositoryTestSuite) SetupTest() {
	sqlDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	gormdb, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB, DriverName: "postgres", WithoutQuotingCheck: true}), &gorm.Config{})

	r.repository = NewRepository(Repository{Db: gormdb})
	r.sqlMock = mock
	r.ctx = context.Background()
	loc, _ := time.LoadLocation("UTC")
	r.loc = loc
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (r *RepositoryTestSuite) TestCreateEstate() {
	type fields struct {
		mock func(newEstate Estate)
	}

	type args struct {
		ctx       context.Context
		newEstate *Estate
	}

	estate := Estate{
		Width:  10,
		Length: 20,
	}

	query := `INSERT INTO estates (width,length) VALUES ($1,$2) RETURNING id,created_at,updated_at`

	tests := []struct {
		name        string
		args        args
		fields      fields
		expectedErr error
	}{
		{
			name: "Failed, theres an error in db",
			args: args{
				ctx:       r.ctx,
				newEstate: &estate,
			},
			fields: fields{
				mock: func(newEstate Estate) {
					r.sqlMock.ExpectBegin()
					r.sqlMock.ExpectQuery(query).
						WithArgs(estate.Width, estate.Length).
						WillReturnError(gorm.ErrUnsupportedDriver)
					r.sqlMock.ExpectRollback()

				}},
			expectedErr: gorm.ErrUnsupportedDriver,
		},
		{
			name: "Success",
			args: args{
				ctx:       r.ctx,
				newEstate: &estate,
			},
			fields: fields{
				mock: func(newEstate Estate) {
					r.sqlMock.ExpectBegin()
					r.sqlMock.ExpectQuery(query).
						WithArgs(estate.Width, estate.Length).
						WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("f0f40954-d0c8-4a1a-9d54-1b4e57e2e236",
							time.Date(2020, 01, 03, 00, 00, 00, 00, r.loc),
							time.Date(2020, 01, 03, 00, 00, 00, 00, r.loc)))
					r.sqlMock.ExpectCommit()
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		r.Suite.Run(test.name, func() {
			test.fields.mock(*test.args.newEstate)

			actualError := r.repository.CreateEstate(test.args.ctx, test.args.newEstate)

			assert.Equal(r.T(), test.expectedErr, actualError)
		})
	}
}

func (r *RepositoryTestSuite) TestGetEstateByID() {
	type fields struct {
		mock func(estateID string)
	}

	type args struct {
		ctx      context.Context
		estateID string
	}

	query := `SELECT id,width,length FROM estates WHERE id = $1 ORDER BY estates.id LIMIT $2`

	tests := []struct {
		name           string
		args           args
		fields         fields
		expectedResult Estate
		expectedErr    error
	}{
		{
			name: "Failed, theres an error in db",
			args: args{
				ctx:      r.ctx,
				estateID: "f0f40954-d0c8-4a1a-9d54-1b4e57e2e236",
			},
			fields: fields{
				mock: func(estateID string) {
					r.sqlMock.MatchExpectationsInOrder(false)
					r.sqlMock.ExpectQuery(query).WithArgs(estateID, 1).WillReturnError(sql.ErrConnDone)
				}},
			expectedResult: Estate{},
			expectedErr:    sql.ErrConnDone,
		},
		{
			name: "Success",
			args: args{
				ctx:      r.ctx,
				estateID: "f0f40954-d0c8-4a1a-9d54-1b4e57e2e236",
			},
			fields: fields{
				mock: func(id string) {
					r.sqlMock.MatchExpectationsInOrder(false)
					r.sqlMock.ExpectQuery(query).WithArgs(id, 1).
						WillReturnRows(r.sqlMock.NewRows([]string{"id", "width", "length"}).
							AddRow("f0f40954-d0c8-4a1a-9d54-1b4e57e2e236", 10, 20))
				}},
			expectedResult: Estate{
				ID:     "f0f40954-d0c8-4a1a-9d54-1b4e57e2e236",
				Width:  10,
				Length: 20,
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		r.Suite.Run(test.name, func() {
			test.fields.mock(test.args.estateID)

			actualResult, actualErr := r.repository.GetEstateByID(test.args.ctx, test.args.estateID)

			assert.Equal(r.T(), test.expectedErr, actualErr)
			assert.Equal(r.T(), test.expectedResult, actualResult)
		})
	}
}

func (r *RepositoryTestSuite) TestCreateTree() {
	type fields struct {
		mock func(newTree Tree)
	}

	type args struct {
		ctx     context.Context
		newTree *Tree
	}

	tree := Tree{
		EstateID:           "c2dfd742-6a55-41be-b84a-4396f21e2b26",
		HorizontalPosition: 5,
		VerticalPosition:   10,
		Height:             15,
	}

	query := `INSERT INTO trees (estate_id,horizontal_position,vertical_position,height) VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at`

	tests := []struct {
		name        string
		args        args
		fields      fields
		expectedErr error
	}{
		{
			name: "Failed, theres an error in db",
			args: args{
				ctx:     r.ctx,
				newTree: &tree,
			},
			fields: fields{
				mock: func(newTree Tree) {
					r.sqlMock.ExpectBegin()
					r.sqlMock.ExpectQuery(query).
						WithArgs(tree.EstateID, tree.HorizontalPosition, tree.VerticalPosition, tree.Height).
						WillReturnError(gorm.ErrUnsupportedDriver)
					r.sqlMock.ExpectRollback()

				}},
			expectedErr: gorm.ErrUnsupportedDriver,
		},
		{
			name: "Success",
			args: args{
				ctx:     r.ctx,
				newTree: &tree,
			},
			fields: fields{
				mock: func(newTree Tree) {
					r.sqlMock.ExpectBegin()
					r.sqlMock.ExpectQuery(query).
						WithArgs(tree.EstateID, tree.HorizontalPosition, tree.VerticalPosition, tree.Height).
						WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("734c8a10-2c10-404b-b41e-ff6e7f1d0a0b",
							time.Date(2020, 01, 03, 00, 00, 00, 00, r.loc),
							time.Date(2020, 01, 03, 00, 00, 00, 00, r.loc)))
					r.sqlMock.ExpectCommit()
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		r.Suite.Run(test.name, func() {
			test.fields.mock(*test.args.newTree)

			actualError := r.repository.CreateTree(test.args.ctx, test.args.newTree)

			assert.Equal(r.T(), test.expectedErr, actualError)
		})
	}
}

func (r *RepositoryTestSuite) TestGetTreeHeightsByEstateID() {
	type fields struct {
		mock func(estateID string)
	}

	type args struct {
		ctx      context.Context
		estateID string
	}

	query := `SELECT height FROM trees WHERE estate_id = $1 ORDER BY height asc`

	tests := []struct {
		name           string
		args           args
		fields         fields
		expectedResult []int
		expectedErr    error
	}{
		{
			name: "Failed, theres an error in db",
			args: args{
				ctx:      r.ctx,
				estateID: "f0f40954-d0c8-4a1a-9d54-1b4e57e2e236",
			},
			fields: fields{
				mock: func(estateID string) {
					r.sqlMock.MatchExpectationsInOrder(false)
					r.sqlMock.ExpectQuery(query).WithArgs(estateID).WillReturnError(sql.ErrConnDone)
				}},
			expectedResult: []int(nil),
			expectedErr:    sql.ErrConnDone,
		},
		{
			name: "Success",
			args: args{
				ctx:      r.ctx,
				estateID: "f0f40954-d0c8-4a1a-9d54-1b4e57e2e236",
			},
			fields: fields{
				mock: func(id string) {
					r.sqlMock.MatchExpectationsInOrder(false)
					r.sqlMock.ExpectQuery(query).WithArgs(id).
						WillReturnRows(r.sqlMock.NewRows([]string{"height"}).
							AddRow(1).
							AddRow(5).
							AddRow(10))
				}},
			expectedResult: []int{1, 5, 10},
			expectedErr:    nil,
		},
	}

	for _, test := range tests {
		r.Suite.Run(test.name, func() {
			test.fields.mock(test.args.estateID)

			actualResult, actualErr := r.repository.GetTreeHeightsByEstateID(test.args.ctx, test.args.estateID)

			assert.Equal(r.T(), test.expectedErr, actualErr)
			assert.Equal(r.T(), test.expectedResult, actualResult)
		})
	}
}
