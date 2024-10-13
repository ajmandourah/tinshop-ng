// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ajmandourah/tinshop-ng/repository (interfaces: Config)

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	repository "github.com/ajmandourah/tinshop-ng/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// AddBeforeHook mocks base method.
func (m *MockConfig) AddBeforeHook(arg0 func(repository.Config)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddBeforeHook", arg0)
}

// AddBeforeHook indicates an expected call of AddBeforeHook.
func (mr *MockConfigMockRecorder) AddBeforeHook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBeforeHook", reflect.TypeOf((*MockConfig)(nil).AddBeforeHook), arg0)
}

// AddHook mocks base method.
func (m *MockConfig) AddHook(arg0 func(repository.Config)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddHook", arg0)
}

// AddHook indicates an expected call of AddHook.
func (mr *MockConfigMockRecorder) AddHook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHook", reflect.TypeOf((*MockConfig)(nil).AddHook), arg0)
}

// BannedTheme mocks base method.
func (m *MockConfig) BannedTheme() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BannedTheme")
	ret0, _ := ret[0].([]string)
	return ret0
}

// BannedTheme indicates an expected call of BannedTheme.
func (mr *MockConfigMockRecorder) BannedTheme() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BannedTheme", reflect.TypeOf((*MockConfig)(nil).BannedTheme))
}

// CustomDB mocks base method.
func (m *MockConfig) CustomDB() map[string]repository.TitleDBEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CustomDB")
	ret0, _ := ret[0].(map[string]repository.TitleDBEntry)
	return ret0
}

// CustomDB indicates an expected call of CustomDB.
func (mr *MockConfigMockRecorder) CustomDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CustomDB", reflect.TypeOf((*MockConfig)(nil).CustomDB))
}

// DebugNfs mocks base method.
func (m *MockConfig) DebugNfs() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebugNfs")
	ret0, _ := ret[0].(bool)
	return ret0
}

// DebugNfs indicates an expected call of DebugNfs.
func (mr *MockConfigMockRecorder) DebugNfs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebugNfs", reflect.TypeOf((*MockConfig)(nil).DebugNfs))
}

// DebugNoSecurity mocks base method.
func (m *MockConfig) DebugNoSecurity() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebugNoSecurity")
	ret0, _ := ret[0].(bool)
	return ret0
}

// DebugNoSecurity indicates an expected call of DebugNoSecurity.
func (mr *MockConfigMockRecorder) DebugNoSecurity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebugNoSecurity", reflect.TypeOf((*MockConfig)(nil).DebugNoSecurity))
}

// DebugTicket mocks base method.
func (m *MockConfig) DebugTicket() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebugTicket")
	ret0, _ := ret[0].(bool)
	return ret0
}

// DebugTicket indicates an expected call of DebugTicket.
func (mr *MockConfigMockRecorder) DebugTicket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebugTicket", reflect.TypeOf((*MockConfig)(nil).DebugTicket))
}

// Directories mocks base method.
func (m *MockConfig) Directories() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Directories")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Directories indicates an expected call of Directories.
func (mr *MockConfigMockRecorder) Directories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Directories", reflect.TypeOf((*MockConfig)(nil).Directories))
}

// ForwardAuthURL mocks base method.
func (m *MockConfig) ForwardAuthURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForwardAuthURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// ForwardAuthURL indicates an expected call of ForwardAuthURL.
func (mr *MockConfigMockRecorder) ForwardAuthURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardAuthURL", reflect.TypeOf((*MockConfig)(nil).ForwardAuthURL))
}

// Host mocks base method.
func (m *MockConfig) Host() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Host")
	ret0, _ := ret[0].(string)
	return ret0
}

// Host indicates an expected call of Host.
func (mr *MockConfigMockRecorder) Host() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Host", reflect.TypeOf((*MockConfig)(nil).Host))
}

// IsBannedTheme mocks base method.
func (m *MockConfig) IsBannedTheme(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBannedTheme", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBannedTheme indicates an expected call of IsBannedTheme.
func (mr *MockConfigMockRecorder) IsBannedTheme(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBannedTheme", reflect.TypeOf((*MockConfig)(nil).IsBannedTheme), arg0)
}

// IsBlacklisted mocks base method.
func (m *MockConfig) IsBlacklisted(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBlacklisted", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsBlacklisted indicates an expected call of IsBlacklisted.
func (mr *MockConfigMockRecorder) IsBlacklisted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBlacklisted", reflect.TypeOf((*MockConfig)(nil).IsBlacklisted), arg0)
}

// IsWhitelisted mocks base method.
func (m *MockConfig) IsWhitelisted(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsWhitelisted", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsWhitelisted indicates an expected call of IsWhitelisted.
func (mr *MockConfigMockRecorder) IsWhitelisted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsWhitelisted", reflect.TypeOf((*MockConfig)(nil).IsWhitelisted), arg0)
}

// LoadConfig mocks base method.
func (m *MockConfig) LoadConfig() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadConfig")
}

// LoadConfig indicates an expected call of LoadConfig.
func (mr *MockConfigMockRecorder) LoadConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockConfig)(nil).LoadConfig))
}

// NfsShares mocks base method.
func (m *MockConfig) NfsShares() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NfsShares")
	ret0, _ := ret[0].([]string)
	return ret0
}

// NfsShares indicates an expected call of NfsShares.
func (mr *MockConfigMockRecorder) NfsShares() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NfsShares", reflect.TypeOf((*MockConfig)(nil).NfsShares))
}

// NoWelcomeMessage mocks base method.
func (m *MockConfig) NoWelcomeMessage() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NoWelcomeMessage")
	ret0, _ := ret[0].(bool)
	return ret0
}

// NoWelcomeMessage indicates an expected call of NoWelcomeMessage.
func (mr *MockConfigMockRecorder) NoWelcomeMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NoWelcomeMessage", reflect.TypeOf((*MockConfig)(nil).NoWelcomeMessage))
}

// Port mocks base method.
func (m *MockConfig) Port() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Port")
	ret0, _ := ret[0].(int)
	return ret0
}

// Port indicates an expected call of Port.
func (mr *MockConfigMockRecorder) Port() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Port", reflect.TypeOf((*MockConfig)(nil).Port))
}

// Protocol mocks base method.
func (m *MockConfig) Protocol() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Protocol")
	ret0, _ := ret[0].(string)
	return ret0
}

// Protocol indicates an expected call of Protocol.
func (mr *MockConfigMockRecorder) Protocol() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Protocol", reflect.TypeOf((*MockConfig)(nil).Protocol))
}

// ReverseProxy mocks base method.
func (m *MockConfig) ReverseProxy() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReverseProxy")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ReverseProxy indicates an expected call of ReverseProxy.
func (mr *MockConfigMockRecorder) ReverseProxy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReverseProxy", reflect.TypeOf((*MockConfig)(nil).ReverseProxy))
}

// RootShop mocks base method.
func (m *MockConfig) RootShop() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RootShop")
	ret0, _ := ret[0].(string)
	return ret0
}

// RootShop indicates an expected call of RootShop.
func (mr *MockConfigMockRecorder) RootShop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RootShop", reflect.TypeOf((*MockConfig)(nil).RootShop))
}

// SetRootShop mocks base method.
func (m *MockConfig) SetRootShop(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRootShop", arg0)
}

// SetRootShop indicates an expected call of SetRootShop.
func (mr *MockConfigMockRecorder) SetRootShop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRootShop", reflect.TypeOf((*MockConfig)(nil).SetRootShop), arg0)
}

// SetShopTemplateData mocks base method.
func (m *MockConfig) SetShopTemplateData(arg0 repository.ShopTemplate) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetShopTemplateData", arg0)
}

// SetShopTemplateData indicates an expected call of SetShopTemplateData.
func (mr *MockConfigMockRecorder) SetShopTemplateData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetShopTemplateData", reflect.TypeOf((*MockConfig)(nil).SetShopTemplateData), arg0)
}

// ShopTemplateData mocks base method.
func (m *MockConfig) ShopTemplateData() repository.ShopTemplate {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShopTemplateData")
	ret0, _ := ret[0].(repository.ShopTemplate)
	return ret0
}

// ShopTemplateData indicates an expected call of ShopTemplateData.
func (mr *MockConfigMockRecorder) ShopTemplateData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShopTemplateData", reflect.TypeOf((*MockConfig)(nil).ShopTemplateData))
}

// ShopTitle mocks base method.
func (m *MockConfig) ShopTitle() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShopTitle")
	ret0, _ := ret[0].(string)
	return ret0
}

// ShopTitle indicates an expected call of ShopTitle.
func (mr *MockConfigMockRecorder) ShopTitle() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShopTitle", reflect.TypeOf((*MockConfig)(nil).ShopTitle))
}

// Sources mocks base method.
func (m *MockConfig) Sources() repository.ConfigSources {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sources")
	ret0, _ := ret[0].(repository.ConfigSources)
	return ret0
}

// Sources indicates an expected call of Sources.
func (mr *MockConfigMockRecorder) Sources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sources", reflect.TypeOf((*MockConfig)(nil).Sources))
}

// VerifyNSP mocks base method.
func (m *MockConfig) VerifyNSP() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyNSP")
	ret0, _ := ret[0].(bool)
	return ret0
}

// VerifyNSP indicates an expected call of VerifyNSP.
func (mr *MockConfigMockRecorder) VerifyNSP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyNSP", reflect.TypeOf((*MockConfig)(nil).VerifyNSP))
}

// WelcomeMessage mocks base method.
func (m *MockConfig) WelcomeMessage() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WelcomeMessage")
	ret0, _ := ret[0].(string)
	return ret0
}

// WelcomeMessage indicates an expected call of WelcomeMessage.
func (mr *MockConfigMockRecorder) WelcomeMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WelcomeMessage", reflect.TypeOf((*MockConfig)(nil).WelcomeMessage))
}
