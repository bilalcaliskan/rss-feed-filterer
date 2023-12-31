// Code generated by mockery v2.39.1. DO NOT EDIT.

package feed

import (
	gofeed "github.com/mmcdole/gofeed"
	mock "github.com/stretchr/testify/mock"
)

// MockParser is an autogenerated mock type for the Parser type
type MockParser struct {
	mock.Mock
}

type MockParser_Expecter struct {
	mock *mock.Mock
}

func (_m *MockParser) EXPECT() *MockParser_Expecter {
	return &MockParser_Expecter{mock: &_m.Mock}
}

// ParseURL provides a mock function with given fields: url
func (_m *MockParser) ParseURL(url string) (*gofeed.Feed, error) {
	ret := _m.Called(url)

	if len(ret) == 0 {
		panic("no return value specified for ParseURL")
	}

	var r0 *gofeed.Feed
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*gofeed.Feed, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) *gofeed.Feed); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gofeed.Feed)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockParser_ParseURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseURL'
type MockParser_ParseURL_Call struct {
	*mock.Call
}

// ParseURL is a helper method to define mock.On call
//   - url string
func (_e *MockParser_Expecter) ParseURL(url interface{}) *MockParser_ParseURL_Call {
	return &MockParser_ParseURL_Call{Call: _e.mock.On("ParseURL", url)}
}

func (_c *MockParser_ParseURL_Call) Run(run func(url string)) *MockParser_ParseURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockParser_ParseURL_Call) Return(_a0 *gofeed.Feed, _a1 error) *MockParser_ParseURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockParser_ParseURL_Call) RunAndReturn(run func(string) (*gofeed.Feed, error)) *MockParser_ParseURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockParser creates a new instance of MockParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockParser(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockParser {
	mock := &MockParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
