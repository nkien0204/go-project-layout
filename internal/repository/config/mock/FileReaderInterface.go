// Code generated by mockery v2.20.0. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// FileReaderInterface is an autogenerated mock type for the FileReaderInterface type
type FileReaderInterface struct {
	mock.Mock
}

// GetFileName provides a mock function with given fields:
func (_m *FileReaderInterface) GetFileName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ReadFile provides a mock function with given fields:
func (_m *FileReaderInterface) ReadFile() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFileReaderInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewFileReaderInterface creates a new instance of FileReaderInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFileReaderInterface(t mockConstructorTestingTNewFileReaderInterface) *FileReaderInterface {
	mock := &FileReaderInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
