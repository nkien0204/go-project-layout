// Code generated by mockery v2.35.4. DO NOT EDIT.

package mock

import (
	config "github.com/nkien0204/lets-go/internal/domain/entity/config"

	mock "github.com/stretchr/testify/mock"
)

// ConfigUsecase is an autogenerated mock type for the ConfigUsecase type
type ConfigUsecase struct {
	mock.Mock
}

// LoadConfig provides a mock function with given fields:
func (_m *ConfigUsecase) LoadConfig() *config.Cfg {
	ret := _m.Called()

	var r0 *config.Cfg
	if rf, ok := ret.Get(0).(func() *config.Cfg); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Cfg)
		}
	}

	return r0
}

// NewConfigUsecase creates a new instance of ConfigUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConfigUsecase {
	mock := &ConfigUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
