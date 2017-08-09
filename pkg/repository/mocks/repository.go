// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/trussle/snowy/pkg/repository (interfaces: Repository)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	document "github.com/trussle/snowy/pkg/document"
	repository "github.com/trussle/snowy/pkg/repository"
	uuid "github.com/trussle/snowy/pkg/uuid"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return _m.recorder
}

// Close mocks base method
func (_m *MockRepository) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (_mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

// GetContent mocks base method
func (_m *MockRepository) GetContent(_param0 uuid.UUID) (document.Content, error) {
	ret := _m.ctrl.Call(_m, "GetContent", _param0)
	ret0, _ := ret[0].(document.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContent indicates an expected call of GetContent
func (_mr *MockRepositoryMockRecorder) GetContent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetContent", arg0)
}

// GetDocument mocks base method
func (_m *MockRepository) GetDocument(_param0 uuid.UUID, _param1 repository.Query) (document.Document, error) {
	ret := _m.ctrl.Call(_m, "GetDocument", _param0, _param1)
	ret0, _ := ret[0].(document.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDocument indicates an expected call of GetDocument
func (_mr *MockRepositoryMockRecorder) GetDocument(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDocument", arg0, arg1)
}

// GetDocuments mocks base method
func (_m *MockRepository) GetDocuments(_param0 uuid.UUID, _param1 repository.Query) ([]document.Document, error) {
	ret := _m.ctrl.Call(_m, "GetDocuments", _param0, _param1)
	ret0, _ := ret[0].([]document.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDocuments indicates an expected call of GetDocuments
func (_mr *MockRepositoryMockRecorder) GetDocuments(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDocuments", arg0, arg1)
}

// InsertDocument mocks base method
func (_m *MockRepository) InsertDocument(_param0 document.Document) (document.Document, error) {
	ret := _m.ctrl.Call(_m, "InsertDocument", _param0)
	ret0, _ := ret[0].(document.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertDocument indicates an expected call of InsertDocument
func (_mr *MockRepositoryMockRecorder) InsertDocument(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InsertDocument", arg0)
}

// PutContent mocks base method
func (_m *MockRepository) PutContent(_param0 document.Content) (uuid.UUID, error) {
	ret := _m.ctrl.Call(_m, "PutContent", _param0)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutContent indicates an expected call of PutContent
func (_mr *MockRepositoryMockRecorder) PutContent(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutContent", arg0)
}
