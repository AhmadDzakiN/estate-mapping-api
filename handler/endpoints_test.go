package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EndpointsTestSuite struct {
	suite.Suite
	repositoryMock *repository.MockRepositoryInterface
	echo           *echo.Echo
	server         *Server
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (e *EndpointsTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(e.T())
	defer mockCtrl.Finish()

	e.repositoryMock = repository.NewMockRepositoryInterface(mockCtrl)
	e.server = NewServer(NewServerOptions{
		Repository: e.repositoryMock,
	})

	e.echo = echo.New()
	e.echo.Validator = &CustomValidator{validator: validator.New()}
}

func TestEndpointsSuite(t *testing.T) {
	suite.Run(t, new(EndpointsTestSuite))
}

func (e *EndpointsTestSuite) TestCreateEstate() {
	type fields struct {
		mock func(ctx echo.Context)
	}

	type args struct {
		reqBody string
	}

	tests := []struct {
		name               string
		args               args
		fields             fields
		expectedErr        string
		expectedStatusCode int
	}{
		{
			name: "Failed, invalid request body format",
			args: args{
				reqBody: `{"width": "test", "length": 20}`,
			},
			fields: fields{
				mock: func(ctx echo.Context) {},
			},
			expectedErr:        "Invalid input",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failed, length/width < 1 or > 50000",
			args: args{
				reqBody: `{"width": -1, "length": 20}`,
			},
			fields: fields{
				mock: func(ctx echo.Context) {},
			},
			expectedErr:        "Invalid input",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failed, got error from repo",
			args: args{
				reqBody: `{"width": 10, "length": 20}`,
			},
			fields: fields{
				mock: func(ctx echo.Context) {
					newEstate := repository.Estate{
						Length: 20,
						Width:  10,
					}
					e.repositoryMock.EXPECT().CreateEstate(ctx.Request().Context(), &newEstate).Return(sql.ErrConnDone)
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Success",
			args: args{
				reqBody: `{"width": 10, "length": 20}`,
			},
			fields: fields{
				mock: func(ctx echo.Context) {
					newEstate := repository.Estate{
						Length: 20,
						Width:  10,
					}
					e.repositoryMock.EXPECT().CreateEstate(ctx.Request().Context(), &newEstate).Return(nil)
				},
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		e.Suite.Run(test.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/estates", strings.NewReader(test.args.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.echo.NewContext(req, rec)

			test.fields.mock(ctx)

			err := e.server.CreateEstate(ctx)
			assert.NoError(e.T(), err)

			var resp generated.InvalidInputErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(e.T(), err)

			assert.Equal(e.T(), test.expectedStatusCode, rec.Code)
			assert.Equal(e.T(), test.expectedErr, resp.Error)
		})
	}
}

func (e *EndpointsTestSuite) TestCreateTree() {
	type fields struct {
		mock func(ctx echo.Context, estateID openapi_types.UUID)
	}

	type args struct {
		reqBody  string
		estateID openapi_types.UUID
	}

	tests := []struct {
		name               string
		args               args
		fields             fields
		expectedErr        string
		expectedStatusCode int
	}{
		{
			name: "Failed, invalid request body format",
			args: args{
				reqBody:  `{"x": "test", "y": 20, "height": 2}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {},
			},
			expectedErr:        "Invalid input",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failed, length/width < 1 or > 50000 or height > 30",
			args: args{
				reqBody:  `{"x": 1, "y": 20, "height": 32}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {},
			},
			expectedErr:        "Invalid input",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failed, got error non record not found for GetEstateByID repo",
			args: args{
				reqBody:  `{"x": 1, "y": 20, "height": 15}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{}, sql.ErrConnDone)
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Failed, estate not found for GetEstateByID repo",
			args: args{
				reqBody:  `{"x": 1, "y": 20, "height": 15}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{}, gorm.ErrRecordNotFound)
				},
			},
			expectedErr:        "Estate not found",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Failed, tree position is out of the estate's area",
			args: args{
				reqBody:  `{"x": 1, "y": 20, "height": 15}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID:     estateID.String(),
						Length: 5,
						Width:  15,
					}, nil)
				},
			},
			expectedErr:        "Tree position is out of the estate's area",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failed, got error from CreateTree repo",
			args: args{
				reqBody:  `{"x": 1, "y": 20, "height": 15}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID:     estateID.String(),
						Length: 5,
						Width:  20,
					}, nil)
					e.repositoryMock.EXPECT().CreateTree(ctx.Request().Context(), &repository.Tree{
						EstateID:           estateID.String(),
						HorizontalPosition: 1,
						VerticalPosition:   20,
						Height:             15,
					}).Return(sql.ErrConnDone)
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Success",
			args: args{
				reqBody:  `{"x": 1, "y": 20, "height": 15}`,
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID:     estateID.String(),
						Length: 5,
						Width:  20,
					}, nil)
					e.repositoryMock.EXPECT().CreateTree(ctx.Request().Context(), &repository.Tree{
						EstateID:           estateID.String(),
						HorizontalPosition: 1,
						VerticalPosition:   20,
						Height:             15,
					}).Return(nil)
				},
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		e.Suite.Run(test.name, func() {

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estates/%s/tree", test.args.estateID), strings.NewReader(test.args.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.echo.NewContext(req, rec)

			test.fields.mock(ctx, test.args.estateID)

			err := e.server.CreateTree(ctx, test.args.estateID)
			assert.NoError(e.T(), err)

			var resp generated.InvalidInputErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(e.T(), err)

			assert.Equal(e.T(), test.expectedStatusCode, rec.Code)
			assert.Equal(e.T(), test.expectedErr, resp.Error)
		})
	}
}

func (e *EndpointsTestSuite) TestGetEstateDronePlan() {
	type fields struct {
		mock func(ctx echo.Context, estateID openapi_types.UUID)
	}

	type args struct {
		estateID openapi_types.UUID
	}

	tests := []struct {
		name               string
		args               args
		fields             fields
		expectedErr        string
		expectedStatusCode int
	}{
		{
			name: "Failed, estate not found for GetEstateByID",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{}, gorm.ErrRecordNotFound)
				},
			},
			expectedErr:        "Estate not found",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Failed, got error non record not found for GetEstateByID repo",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{}, errors.New("random error"))
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Failed, got error for GetTreeHeightsByEstateID repo",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID: estateID.String(),
					}, nil)
					e.repositoryMock.EXPECT().GetTreeHeightsByEstateID(ctx.Request().Context(), estateID.String()).Return([]int(nil), errors.New("random error"))
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Success, with tree count is even",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID: estateID.String(),
					}, nil)
					e.repositoryMock.EXPECT().GetTreeHeightsByEstateID(ctx.Request().Context(), estateID.String()).Return([]int{1, 2, 3, 4}, nil)
				},
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Success, with tree count is odd",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID: estateID.String(),
					}, nil)
					e.repositoryMock.EXPECT().GetTreeHeightsByEstateID(ctx.Request().Context(), estateID.String()).Return([]int{1, 2, 3}, nil)
				},
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		e.Suite.Run(test.name, func() {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estates/%s/stats", test.args.estateID), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.echo.NewContext(req, rec)

			test.fields.mock(ctx, test.args.estateID)

			err := e.server.GetEstateStats(ctx, test.args.estateID)
			assert.NoError(e.T(), err)

			var resp generated.InvalidInputErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(e.T(), err)

			assert.Equal(e.T(), test.expectedStatusCode, rec.Code)
			assert.Equal(e.T(), test.expectedErr, resp.Error)
		})
	}
}

func (e *EndpointsTestSuite) TestGetEstateStats() {
	type fields struct {
		mock func(ctx echo.Context, estateID openapi_types.UUID)
	}

	type args struct {
		estateID openapi_types.UUID
	}

	tests := []struct {
		name               string
		args               args
		fields             fields
		expectedErr        string
		expectedStatusCode int
	}{
		{
			name: "Failed, estate not found for GetEstateByID",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{}, gorm.ErrRecordNotFound)
				},
			},
			expectedErr:        "Estate not found",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Failed, got error non record not found for GetEstateByID repo",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{}, errors.New("random error"))
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Failed, got error for GetTreesByEstateIDAndPlotsLocations repo",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID: estateID.String(),
					}, nil)
					e.repositoryMock.EXPECT().GetTreesByEstateIDAndPlotsLocations(ctx.Request().Context(), estateID.String()).Return([]repository.Tree(nil), errors.New("random error"))
				},
			},
			expectedErr:        "Oops, something wrong with the server. Please try again later",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Success, with estate do not have any single tree",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID: estateID.String(),
					}, nil)
					e.repositoryMock.EXPECT().GetTreesByEstateIDAndPlotsLocations(ctx.Request().Context(), estateID.String()).Return([]repository.Tree(nil), nil)
				},
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Success, with estate have tree in its plot",
			args: args{
				estateID: uuid.New(),
			},
			fields: fields{
				mock: func(ctx echo.Context, estateID openapi_types.UUID) {
					e.repositoryMock.EXPECT().GetEstateByID(ctx.Request().Context(), estateID.String()).Return(repository.Estate{
						ID: estateID.String(),
					}, nil)
					e.repositoryMock.EXPECT().GetTreesByEstateIDAndPlotsLocations(ctx.Request().Context(), estateID.String()).Return([]repository.Tree{
						{
							ID:                 uuid.New().String(),
							EstateID:           uuid.New().String(),
							HorizontalPosition: 2,
							VerticalPosition:   3,
							Height:             5,
						},
					}, nil)
				},
			},
			expectedErr:        "",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		e.Suite.Run(test.name, func() {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/estates/%s/drone-plan", test.args.estateID), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.echo.NewContext(req, rec)

			test.fields.mock(ctx, test.args.estateID)

			err := e.server.GetEstateDronePlan(ctx, test.args.estateID)
			assert.NoError(e.T(), err)

			var resp generated.InvalidInputErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(e.T(), err)

			assert.Equal(e.T(), test.expectedStatusCode, rec.Code)
			assert.Equal(e.T(), test.expectedErr, resp.Error)
		})
	}
}
