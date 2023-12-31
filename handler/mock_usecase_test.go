// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/usecase.go

// Package handler_test is a generated GoMock package.
package handler_test

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// GenerateURL mocks base method.
func (m *MockUsecase) GenerateURL(ctx context.Context, originalURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateURL", ctx, originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateURL indicates an expected call of GenerateURL.
func (mr *MockUsecaseMockRecorder) GenerateURL(ctx, originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateURL", reflect.TypeOf((*MockUsecase)(nil).GenerateURL), ctx, originalURL)
}

// Health mocks base method.
func (m *MockUsecase) Health(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockUsecaseMockRecorder) Health(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockUsecase)(nil).Health), ctx)
}

// SearchOriginalURL mocks base method.
func (m *MockUsecase) SearchOriginalURL(ctx context.Context, shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchOriginalURL", ctx, shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchOriginalURL indicates an expected call of SearchOriginalURL.
func (mr *MockUsecaseMockRecorder) SearchOriginalURL(ctx, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchOriginalURL", reflect.TypeOf((*MockUsecase)(nil).SearchOriginalURL), ctx, shortURL)
}
