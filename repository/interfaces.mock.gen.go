// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=repository/interfaces.go -destination=repository/interfaces.mock.gen.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreateEstate mocks base method.
func (m *MockRepositoryInterface) CreateEstate(ctx context.Context, newEstate *Estate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEstate", ctx, newEstate)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEstate indicates an expected call of CreateEstate.
func (mr *MockRepositoryInterfaceMockRecorder) CreateEstate(ctx, newEstate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEstate", reflect.TypeOf((*MockRepositoryInterface)(nil).CreateEstate), ctx, newEstate)
}

// CreateTree mocks base method.
func (m *MockRepositoryInterface) CreateTree(ctx context.Context, newTree *Tree) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTree", ctx, newTree)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTree indicates an expected call of CreateTree.
func (mr *MockRepositoryInterfaceMockRecorder) CreateTree(ctx, newTree any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTree", reflect.TypeOf((*MockRepositoryInterface)(nil).CreateTree), ctx, newTree)
}

// GetEstateByID mocks base method.
func (m *MockRepositoryInterface) GetEstateByID(ctx context.Context, estateID string) (Estate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEstateByID", ctx, estateID)
	ret0, _ := ret[0].(Estate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEstateByID indicates an expected call of GetEstateByID.
func (mr *MockRepositoryInterfaceMockRecorder) GetEstateByID(ctx, estateID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEstateByID", reflect.TypeOf((*MockRepositoryInterface)(nil).GetEstateByID), ctx, estateID)
}

// GetTreeHeightsByEstateID mocks base method.
func (m *MockRepositoryInterface) GetTreeHeightsByEstateID(ctx context.Context, estateID string) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTreeHeightsByEstateID", ctx, estateID)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTreeHeightsByEstateID indicates an expected call of GetTreeHeightsByEstateID.
func (mr *MockRepositoryInterfaceMockRecorder) GetTreeHeightsByEstateID(ctx, estateID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTreeHeightsByEstateID", reflect.TypeOf((*MockRepositoryInterface)(nil).GetTreeHeightsByEstateID), ctx, estateID)
}

// GetTreesByEstateIDAndPlotsLocations mocks base method.
func (m *MockRepositoryInterface) GetTreesByEstateIDAndPlotsLocations(ctx context.Context, estateID string) ([]Tree, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTreesByEstateIDAndPlotsLocations", ctx, estateID)
	ret0, _ := ret[0].([]Tree)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTreesByEstateIDAndPlotsLocations indicates an expected call of GetTreesByEstateIDAndPlotsLocations.
func (mr *MockRepositoryInterfaceMockRecorder) GetTreesByEstateIDAndPlotsLocations(ctx, estateID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTreesByEstateIDAndPlotsLocations", reflect.TypeOf((*MockRepositoryInterface)(nil).GetTreesByEstateIDAndPlotsLocations), ctx, estateID)
}