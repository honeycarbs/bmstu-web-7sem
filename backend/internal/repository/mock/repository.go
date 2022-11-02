// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	model "neatly/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// AuthorizeAccount mocks base method.
func (m *MockAccountRepository) AuthorizeAccount(a *model.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizeAccount", a)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthorizeAccount indicates an expected call of AuthorizeAccount.
func (mr *MockAccountRepositoryMockRecorder) AuthorizeAccount(a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizeAccount", reflect.TypeOf((*MockAccountRepository)(nil).AuthorizeAccount), a)
}

// CreateAccount mocks base method.
func (m *MockAccountRepository) CreateAccount(a *model.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", a)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAccountRepositoryMockRecorder) CreateAccount(a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountRepository)(nil).CreateAccount), a)
}

// GetOne mocks base method.
func (m *MockAccountRepository) GetOne(userID int) (model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", userID)
	ret0, _ := ret[0].(model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockAccountRepositoryMockRecorder) GetOne(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockAccountRepository)(nil).GetOne), userID)
}

// MockNoteRepository is a mock of NoteRepository interface.
type MockNoteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNoteRepositoryMockRecorder
}

// MockNoteRepositoryMockRecorder is the mock recorder for MockNoteRepository.
type MockNoteRepositoryMockRecorder struct {
	mock *MockNoteRepository
}

// NewMockNoteRepository creates a new mock instance.
func NewMockNoteRepository(ctrl *gomock.Controller) *MockNoteRepository {
	mock := &MockNoteRepository{ctrl: ctrl}
	mock.recorder = &MockNoteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteRepository) EXPECT() *MockNoteRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNoteRepository) Create(userID int, note *model.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID, note)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockNoteRepositoryMockRecorder) Create(userID, note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNoteRepository)(nil).Create), userID, note)
}

// Delete mocks base method.
func (m *MockNoteRepository) Delete(userID, noteID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userID, noteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockNoteRepositoryMockRecorder) Delete(userID, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNoteRepository)(nil).Delete), userID, noteID)
}

// GetAll mocks base method.
func (m *MockNoteRepository) GetAll(userID int) ([]model.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userID)
	ret0, _ := ret[0].([]model.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockNoteRepositoryMockRecorder) GetAll(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockNoteRepository)(nil).GetAll), userID)
}

// GetOne mocks base method.
func (m *MockNoteRepository) GetOne(userID, noteID int) (model.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", userID, noteID)
	ret0, _ := ret[0].(model.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockNoteRepositoryMockRecorder) GetOne(userID, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockNoteRepository)(nil).GetOne), userID, noteID)
}

// Update mocks base method.
func (m *MockNoteRepository) Update(userID int, n model.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userID, n)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockNoteRepositoryMockRecorder) Update(userID, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNoteRepository)(nil).Update), userID, n)
}

// MockTagRepository is a mock of TagRepository interface.
type MockTagRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTagRepositoryMockRecorder
}

// MockTagRepositoryMockRecorder is the mock recorder for MockTagRepository.
type MockTagRepositoryMockRecorder struct {
	mock *MockTagRepository
}

// NewMockTagRepository creates a new mock instance.
func NewMockTagRepository(ctrl *gomock.Controller) *MockTagRepository {
	mock := &MockTagRepository{ctrl: ctrl}
	mock.recorder = &MockTagRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagRepository) EXPECT() *MockTagRepositoryMockRecorder {
	return m.recorder
}

// Assign mocks base method.
func (m *MockTagRepository) Assign(tagID, noteID, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Assign", tagID, noteID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Assign indicates an expected call of Assign.
func (mr *MockTagRepositoryMockRecorder) Assign(tagID, noteID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Assign", reflect.TypeOf((*MockTagRepository)(nil).Assign), tagID, noteID, userID)
}

// Create mocks base method.
func (m *MockTagRepository) Create(userID, noteID int, t *model.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID, noteID, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTagRepositoryMockRecorder) Create(userID, noteID, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTagRepository)(nil).Create), userID, noteID, t)
}

// Delete mocks base method.
func (m *MockTagRepository) Delete(userID, tagID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userID, tagID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTagRepositoryMockRecorder) Delete(userID, tagID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTagRepository)(nil).Delete), userID, tagID)
}

// Detach mocks base method.
func (m *MockTagRepository) Detach(userID, tagID, noteID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detach", userID, tagID, noteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Detach indicates an expected call of Detach.
func (mr *MockTagRepositoryMockRecorder) Detach(userID, tagID, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detach", reflect.TypeOf((*MockTagRepository)(nil).Detach), userID, tagID, noteID)
}

// GetAll mocks base method.
func (m *MockTagRepository) GetAll(userID int) ([]model.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userID)
	ret0, _ := ret[0].([]model.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTagRepositoryMockRecorder) GetAll(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTagRepository)(nil).GetAll), userID)
}

// GetAllByNote mocks base method.
func (m *MockTagRepository) GetAllByNote(userID, noteID int) ([]model.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByNote", userID, noteID)
	ret0, _ := ret[0].([]model.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByNote indicates an expected call of GetAllByNote.
func (mr *MockTagRepositoryMockRecorder) GetAllByNote(userID, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByNote", reflect.TypeOf((*MockTagRepository)(nil).GetAllByNote), userID, noteID)
}

// GetOne mocks base method.
func (m *MockTagRepository) GetOne(userID, tagID int) (model.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", userID, tagID)
	ret0, _ := ret[0].(model.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockTagRepositoryMockRecorder) GetOne(userID, tagID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockTagRepository)(nil).GetOne), userID, tagID)
}

// Update mocks base method.
func (m *MockTagRepository) Update(userID, tagID int, t model.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userID, tagID, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTagRepositoryMockRecorder) Update(userID, tagID, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTagRepository)(nil).Update), userID, tagID, t)
}
