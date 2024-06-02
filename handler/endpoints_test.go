package handler

import (
	"context"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type EndpointsTestSuite struct {
	suite.Suite
	repositoryMock *repository.MockRepositoryInterface
	//validator      *validator.Validate
	ctx context.Context
}

func (e *EndpointsTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(e.T())
	defer mockCtrl.Finish()

	e.repositoryMock = repository.NewMockRepositoryInterface(mockCtrl)

	e.ctx = context.Background()
}
