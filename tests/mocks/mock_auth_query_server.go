// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cosmos/cosmos-sdk/x/auth/types (interfaces: QueryServer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "github.com/cosmos/cosmos-sdk/x/auth/types"
	gomock "github.com/golang/mock/gomock"
)

// MockQueryServer is a mock of QueryServer interface.
type MockQueryServer struct {
	ctrl     *gomock.Controller
	recorder *MockQueryServerMockRecorder
}

// MockQueryServerMockRecorder is the mock recorder for MockQueryServer.
type MockQueryServerMockRecorder struct {
	mock *MockQueryServer
}

// NewMockQueryServer creates a new mock instance.
func NewMockQueryServer(ctrl *gomock.Controller) *MockQueryServer {
	mock := &MockQueryServer{ctrl: ctrl}
	mock.recorder = &MockQueryServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryServer) EXPECT() *MockQueryServerMockRecorder {
	return m.recorder
}

// Account mocks base method.
func (m *MockQueryServer) Account(arg0 context.Context, arg1 *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Account", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Account indicates an expected call of Account.
func (mr *MockQueryServerMockRecorder) Account(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Account", reflect.TypeOf((*MockQueryServer)(nil).Account), arg0, arg1)
}

// AccountAddressByID mocks base method.
func (m *MockQueryServer) AccountAddressByID(arg0 context.Context, arg1 *types.QueryAccountAddressByIDRequest) (*types.QueryAccountAddressByIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountAddressByID", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryAccountAddressByIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountAddressByID indicates an expected call of AccountAddressByID.
func (mr *MockQueryServerMockRecorder) AccountAddressByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountAddressByID", reflect.TypeOf((*MockQueryServer)(nil).AccountAddressByID), arg0, arg1)
}

// AccountInfo mocks base method.
func (m *MockQueryServer) AccountInfo(arg0 context.Context, arg1 *types.QueryAccountInfoRequest) (*types.QueryAccountInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountInfo", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryAccountInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccountInfo indicates an expected call of AccountInfo.
func (mr *MockQueryServerMockRecorder) AccountInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountInfo", reflect.TypeOf((*MockQueryServer)(nil).AccountInfo), arg0, arg1)
}

// Accounts mocks base method.
func (m *MockQueryServer) Accounts(arg0 context.Context, arg1 *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accounts", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryAccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accounts indicates an expected call of Accounts.
func (mr *MockQueryServerMockRecorder) Accounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accounts", reflect.TypeOf((*MockQueryServer)(nil).Accounts), arg0, arg1)
}

// AddressBytesToString mocks base method.
func (m *MockQueryServer) AddressBytesToString(arg0 context.Context, arg1 *types.AddressBytesToStringRequest) (*types.AddressBytesToStringResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddressBytesToString", arg0, arg1)
	ret0, _ := ret[0].(*types.AddressBytesToStringResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddressBytesToString indicates an expected call of AddressBytesToString.
func (mr *MockQueryServerMockRecorder) AddressBytesToString(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddressBytesToString", reflect.TypeOf((*MockQueryServer)(nil).AddressBytesToString), arg0, arg1)
}

// AddressStringToBytes mocks base method.
func (m *MockQueryServer) AddressStringToBytes(arg0 context.Context, arg1 *types.AddressStringToBytesRequest) (*types.AddressStringToBytesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddressStringToBytes", arg0, arg1)
	ret0, _ := ret[0].(*types.AddressStringToBytesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddressStringToBytes indicates an expected call of AddressStringToBytes.
func (mr *MockQueryServerMockRecorder) AddressStringToBytes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddressStringToBytes", reflect.TypeOf((*MockQueryServer)(nil).AddressStringToBytes), arg0, arg1)
}

// Bech32Prefix mocks base method.
func (m *MockQueryServer) Bech32Prefix(arg0 context.Context, arg1 *types.Bech32PrefixRequest) (*types.Bech32PrefixResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bech32Prefix", arg0, arg1)
	ret0, _ := ret[0].(*types.Bech32PrefixResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Bech32Prefix indicates an expected call of Bech32Prefix.
func (mr *MockQueryServerMockRecorder) Bech32Prefix(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bech32Prefix", reflect.TypeOf((*MockQueryServer)(nil).Bech32Prefix), arg0, arg1)
}

// ModuleAccountByName mocks base method.
func (m *MockQueryServer) ModuleAccountByName(arg0 context.Context, arg1 *types.QueryModuleAccountByNameRequest) (*types.QueryModuleAccountByNameResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModuleAccountByName", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryModuleAccountByNameResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModuleAccountByName indicates an expected call of ModuleAccountByName.
func (mr *MockQueryServerMockRecorder) ModuleAccountByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModuleAccountByName", reflect.TypeOf((*MockQueryServer)(nil).ModuleAccountByName), arg0, arg1)
}

// ModuleAccounts mocks base method.
func (m *MockQueryServer) ModuleAccounts(arg0 context.Context, arg1 *types.QueryModuleAccountsRequest) (*types.QueryModuleAccountsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModuleAccounts", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryModuleAccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModuleAccounts indicates an expected call of ModuleAccounts.
func (mr *MockQueryServerMockRecorder) ModuleAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModuleAccounts", reflect.TypeOf((*MockQueryServer)(nil).ModuleAccounts), arg0, arg1)
}

// Params mocks base method.
func (m *MockQueryServer) Params(arg0 context.Context, arg1 *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Params", arg0, arg1)
	ret0, _ := ret[0].(*types.QueryParamsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Params indicates an expected call of Params.
func (mr *MockQueryServerMockRecorder) Params(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Params", reflect.TypeOf((*MockQueryServer)(nil).Params), arg0, arg1)
}
