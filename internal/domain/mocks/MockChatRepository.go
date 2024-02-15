// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/leighmacdonald/gbans/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockChatRepository is an autogenerated mock type for the ChatRepository type
type MockChatRepository struct {
	mock.Mock
}

type MockChatRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockChatRepository) EXPECT() *MockChatRepository_Expecter {
	return &MockChatRepository_Expecter{mock: &_m.Mock}
}

// AddChatHistory provides a mock function with given fields: ctx, message
func (_m *MockChatRepository) AddChatHistory(ctx context.Context, message *domain.PersonMessage) error {
	ret := _m.Called(ctx, message)

	if len(ret) == 0 {
		panic("no return value specified for AddChatHistory")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.PersonMessage) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockChatRepository_AddChatHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddChatHistory'
type MockChatRepository_AddChatHistory_Call struct {
	*mock.Call
}

// AddChatHistory is a helper method to define mock.On call
//   - ctx context.Context
//   - message *domain.PersonMessage
func (_e *MockChatRepository_Expecter) AddChatHistory(ctx interface{}, message interface{}) *MockChatRepository_AddChatHistory_Call {
	return &MockChatRepository_AddChatHistory_Call{Call: _e.mock.On("AddChatHistory", ctx, message)}
}

func (_c *MockChatRepository_AddChatHistory_Call) Run(run func(ctx context.Context, message *domain.PersonMessage)) *MockChatRepository_AddChatHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.PersonMessage))
	})
	return _c
}

func (_c *MockChatRepository_AddChatHistory_Call) Return(_a0 error) *MockChatRepository_AddChatHistory_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockChatRepository_AddChatHistory_Call) RunAndReturn(run func(context.Context, *domain.PersonMessage) error) *MockChatRepository_AddChatHistory_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonMessage provides a mock function with given fields: ctx, messageID
func (_m *MockChatRepository) GetPersonMessage(ctx context.Context, messageID int64) (domain.QueryChatHistoryResult, error) {
	ret := _m.Called(ctx, messageID)

	if len(ret) == 0 {
		panic("no return value specified for GetPersonMessage")
	}

	var r0 domain.QueryChatHistoryResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (domain.QueryChatHistoryResult, error)); ok {
		return rf(ctx, messageID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.QueryChatHistoryResult); ok {
		r0 = rf(ctx, messageID)
	} else {
		r0 = ret.Get(0).(domain.QueryChatHistoryResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, messageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatRepository_GetPersonMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonMessage'
type MockChatRepository_GetPersonMessage_Call struct {
	*mock.Call
}

// GetPersonMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - messageID int64
func (_e *MockChatRepository_Expecter) GetPersonMessage(ctx interface{}, messageID interface{}) *MockChatRepository_GetPersonMessage_Call {
	return &MockChatRepository_GetPersonMessage_Call{Call: _e.mock.On("GetPersonMessage", ctx, messageID)}
}

func (_c *MockChatRepository_GetPersonMessage_Call) Run(run func(ctx context.Context, messageID int64)) *MockChatRepository_GetPersonMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockChatRepository_GetPersonMessage_Call) Return(_a0 domain.QueryChatHistoryResult, _a1 error) *MockChatRepository_GetPersonMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChatRepository_GetPersonMessage_Call) RunAndReturn(run func(context.Context, int64) (domain.QueryChatHistoryResult, error)) *MockChatRepository_GetPersonMessage_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonMessageContext provides a mock function with given fields: ctx, serverID, messageID, paddedMessageCount
func (_m *MockChatRepository) GetPersonMessageContext(ctx context.Context, serverID int, messageID int64, paddedMessageCount int) ([]domain.QueryChatHistoryResult, error) {
	ret := _m.Called(ctx, serverID, messageID, paddedMessageCount)

	if len(ret) == 0 {
		panic("no return value specified for GetPersonMessageContext")
	}

	var r0 []domain.QueryChatHistoryResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int64, int) ([]domain.QueryChatHistoryResult, error)); ok {
		return rf(ctx, serverID, messageID, paddedMessageCount)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int64, int) []domain.QueryChatHistoryResult); ok {
		r0 = rf(ctx, serverID, messageID, paddedMessageCount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.QueryChatHistoryResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int64, int) error); ok {
		r1 = rf(ctx, serverID, messageID, paddedMessageCount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatRepository_GetPersonMessageContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonMessageContext'
type MockChatRepository_GetPersonMessageContext_Call struct {
	*mock.Call
}

// GetPersonMessageContext is a helper method to define mock.On call
//   - ctx context.Context
//   - serverID int
//   - messageID int64
//   - paddedMessageCount int
func (_e *MockChatRepository_Expecter) GetPersonMessageContext(ctx interface{}, serverID interface{}, messageID interface{}, paddedMessageCount interface{}) *MockChatRepository_GetPersonMessageContext_Call {
	return &MockChatRepository_GetPersonMessageContext_Call{Call: _e.mock.On("GetPersonMessageContext", ctx, serverID, messageID, paddedMessageCount)}
}

func (_c *MockChatRepository_GetPersonMessageContext_Call) Run(run func(ctx context.Context, serverID int, messageID int64, paddedMessageCount int)) *MockChatRepository_GetPersonMessageContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int64), args[3].(int))
	})
	return _c
}

func (_c *MockChatRepository_GetPersonMessageContext_Call) Return(_a0 []domain.QueryChatHistoryResult, _a1 error) *MockChatRepository_GetPersonMessageContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChatRepository_GetPersonMessageContext_Call) RunAndReturn(run func(context.Context, int, int64, int) ([]domain.QueryChatHistoryResult, error)) *MockChatRepository_GetPersonMessageContext_Call {
	_c.Call.Return(run)
	return _c
}

// QueryChatHistory provides a mock function with given fields: ctx, filters
func (_m *MockChatRepository) QueryChatHistory(ctx context.Context, filters domain.ChatHistoryQueryFilter) ([]domain.QueryChatHistoryResult, int64, error) {
	ret := _m.Called(ctx, filters)

	if len(ret) == 0 {
		panic("no return value specified for QueryChatHistory")
	}

	var r0 []domain.QueryChatHistoryResult
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChatHistoryQueryFilter) ([]domain.QueryChatHistoryResult, int64, error)); ok {
		return rf(ctx, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ChatHistoryQueryFilter) []domain.QueryChatHistoryResult); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.QueryChatHistoryResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ChatHistoryQueryFilter) int64); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, domain.ChatHistoryQueryFilter) error); ok {
		r2 = rf(ctx, filters)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockChatRepository_QueryChatHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryChatHistory'
type MockChatRepository_QueryChatHistory_Call struct {
	*mock.Call
}

// QueryChatHistory is a helper method to define mock.On call
//   - ctx context.Context
//   - filters domain.ChatHistoryQueryFilter
func (_e *MockChatRepository_Expecter) QueryChatHistory(ctx interface{}, filters interface{}) *MockChatRepository_QueryChatHistory_Call {
	return &MockChatRepository_QueryChatHistory_Call{Call: _e.mock.On("QueryChatHistory", ctx, filters)}
}

func (_c *MockChatRepository_QueryChatHistory_Call) Run(run func(ctx context.Context, filters domain.ChatHistoryQueryFilter)) *MockChatRepository_QueryChatHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ChatHistoryQueryFilter))
	})
	return _c
}

func (_c *MockChatRepository_QueryChatHistory_Call) Return(_a0 []domain.QueryChatHistoryResult, _a1 int64, _a2 error) *MockChatRepository_QueryChatHistory_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockChatRepository_QueryChatHistory_Call) RunAndReturn(run func(context.Context, domain.ChatHistoryQueryFilter) ([]domain.QueryChatHistoryResult, int64, error)) *MockChatRepository_QueryChatHistory_Call {
	_c.Call.Return(run)
	return _c
}

// TopChatters provides a mock function with given fields: ctx, count
func (_m *MockChatRepository) TopChatters(ctx context.Context, count uint64) ([]domain.TopChatterResult, error) {
	ret := _m.Called(ctx, count)

	if len(ret) == 0 {
		panic("no return value specified for TopChatters")
	}

	var r0 []domain.TopChatterResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) ([]domain.TopChatterResult, error)); ok {
		return rf(ctx, count)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) []domain.TopChatterResult); ok {
		r0 = rf(ctx, count)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TopChatterResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, count)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatRepository_TopChatters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TopChatters'
type MockChatRepository_TopChatters_Call struct {
	*mock.Call
}

// TopChatters is a helper method to define mock.On call
//   - ctx context.Context
//   - count uint64
func (_e *MockChatRepository_Expecter) TopChatters(ctx interface{}, count interface{}) *MockChatRepository_TopChatters_Call {
	return &MockChatRepository_TopChatters_Call{Call: _e.mock.On("TopChatters", ctx, count)}
}

func (_c *MockChatRepository_TopChatters_Call) Run(run func(ctx context.Context, count uint64)) *MockChatRepository_TopChatters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *MockChatRepository_TopChatters_Call) Return(_a0 []domain.TopChatterResult, _a1 error) *MockChatRepository_TopChatters_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChatRepository_TopChatters_Call) RunAndReturn(run func(context.Context, uint64) ([]domain.TopChatterResult, error)) *MockChatRepository_TopChatters_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockChatRepository creates a new instance of MockChatRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockChatRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockChatRepository {
	mock := &MockChatRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}