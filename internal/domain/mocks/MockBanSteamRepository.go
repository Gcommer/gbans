// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/leighmacdonald/gbans/internal/domain"
	mock "github.com/stretchr/testify/mock"

	net "net"

	steamid "github.com/leighmacdonald/steamid/v3/steamid"

	time "time"
)

// MockBanSteamRepository is an autogenerated mock type for the BanSteamRepository type
type MockBanSteamRepository struct {
	mock.Mock
}

type MockBanSteamRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBanSteamRepository) EXPECT() *MockBanSteamRepository_Expecter {
	return &MockBanSteamRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, ban, hardDelete
func (_m *MockBanSteamRepository) Delete(ctx context.Context, ban *domain.BanSteam, hardDelete bool) error {
	ret := _m.Called(ctx, ban, hardDelete)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.BanSteam, bool) error); ok {
		r0 = rf(ctx, ban, hardDelete)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBanSteamRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockBanSteamRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - ban *domain.BanSteam
//   - hardDelete bool
func (_e *MockBanSteamRepository_Expecter) Delete(ctx interface{}, ban interface{}, hardDelete interface{}) *MockBanSteamRepository_Delete_Call {
	return &MockBanSteamRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, ban, hardDelete)}
}

func (_c *MockBanSteamRepository_Delete_Call) Run(run func(ctx context.Context, ban *domain.BanSteam, hardDelete bool)) *MockBanSteamRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.BanSteam), args[2].(bool))
	})
	return _c
}

func (_c *MockBanSteamRepository_Delete_Call) Return(_a0 error) *MockBanSteamRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBanSteamRepository_Delete_Call) RunAndReturn(run func(context.Context, *domain.BanSteam, bool) error) *MockBanSteamRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// ExpiredBans provides a mock function with given fields: ctx
func (_m *MockBanSteamRepository) ExpiredBans(ctx context.Context) ([]domain.BanSteam, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ExpiredBans")
	}

	var r0 []domain.BanSteam
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.BanSteam, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.BanSteam); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BanSteam)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBanSteamRepository_ExpiredBans_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExpiredBans'
type MockBanSteamRepository_ExpiredBans_Call struct {
	*mock.Call
}

// ExpiredBans is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockBanSteamRepository_Expecter) ExpiredBans(ctx interface{}) *MockBanSteamRepository_ExpiredBans_Call {
	return &MockBanSteamRepository_ExpiredBans_Call{Call: _e.mock.On("ExpiredBans", ctx)}
}

func (_c *MockBanSteamRepository_ExpiredBans_Call) Run(run func(ctx context.Context)) *MockBanSteamRepository_ExpiredBans_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockBanSteamRepository_ExpiredBans_Call) Return(_a0 []domain.BanSteam, _a1 error) *MockBanSteamRepository_ExpiredBans_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBanSteamRepository_ExpiredBans_Call) RunAndReturn(run func(context.Context) ([]domain.BanSteam, error)) *MockBanSteamRepository_ExpiredBans_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, filter
func (_m *MockBanSteamRepository) Get(ctx context.Context, filter domain.SteamBansQueryFilter) ([]domain.BannedSteamPerson, int64, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []domain.BannedSteamPerson
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.SteamBansQueryFilter) ([]domain.BannedSteamPerson, int64, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.SteamBansQueryFilter) []domain.BannedSteamPerson); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BannedSteamPerson)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.SteamBansQueryFilter) int64); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, domain.SteamBansQueryFilter) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockBanSteamRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockBanSteamRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - filter domain.SteamBansQueryFilter
func (_e *MockBanSteamRepository_Expecter) Get(ctx interface{}, filter interface{}) *MockBanSteamRepository_Get_Call {
	return &MockBanSteamRepository_Get_Call{Call: _e.mock.On("Get", ctx, filter)}
}

func (_c *MockBanSteamRepository_Get_Call) Run(run func(ctx context.Context, filter domain.SteamBansQueryFilter)) *MockBanSteamRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.SteamBansQueryFilter))
	})
	return _c
}

func (_c *MockBanSteamRepository_Get_Call) Return(_a0 []domain.BannedSteamPerson, _a1 int64, _a2 error) *MockBanSteamRepository_Get_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockBanSteamRepository_Get_Call) RunAndReturn(run func(context.Context, domain.SteamBansQueryFilter) ([]domain.BannedSteamPerson, int64, error)) *MockBanSteamRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByBanID provides a mock function with given fields: ctx, banID, deletedOk
func (_m *MockBanSteamRepository) GetByBanID(ctx context.Context, banID int64, deletedOk bool) (domain.BannedSteamPerson, error) {
	ret := _m.Called(ctx, banID, deletedOk)

	if len(ret) == 0 {
		panic("no return value specified for GetByBanID")
	}

	var r0 domain.BannedSteamPerson
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, bool) (domain.BannedSteamPerson, error)); ok {
		return rf(ctx, banID, deletedOk)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, bool) domain.BannedSteamPerson); ok {
		r0 = rf(ctx, banID, deletedOk)
	} else {
		r0 = ret.Get(0).(domain.BannedSteamPerson)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, bool) error); ok {
		r1 = rf(ctx, banID, deletedOk)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBanSteamRepository_GetByBanID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByBanID'
type MockBanSteamRepository_GetByBanID_Call struct {
	*mock.Call
}

// GetByBanID is a helper method to define mock.On call
//   - ctx context.Context
//   - banID int64
//   - deletedOk bool
func (_e *MockBanSteamRepository_Expecter) GetByBanID(ctx interface{}, banID interface{}, deletedOk interface{}) *MockBanSteamRepository_GetByBanID_Call {
	return &MockBanSteamRepository_GetByBanID_Call{Call: _e.mock.On("GetByBanID", ctx, banID, deletedOk)}
}

func (_c *MockBanSteamRepository_GetByBanID_Call) Run(run func(ctx context.Context, banID int64, deletedOk bool)) *MockBanSteamRepository_GetByBanID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(bool))
	})
	return _c
}

func (_c *MockBanSteamRepository_GetByBanID_Call) Return(_a0 domain.BannedSteamPerson, _a1 error) *MockBanSteamRepository_GetByBanID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBanSteamRepository_GetByBanID_Call) RunAndReturn(run func(context.Context, int64, bool) (domain.BannedSteamPerson, error)) *MockBanSteamRepository_GetByBanID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByLastIP provides a mock function with given fields: ctx, lastIP, deletedOk
func (_m *MockBanSteamRepository) GetByLastIP(ctx context.Context, lastIP net.IP, deletedOk bool) (domain.BannedSteamPerson, error) {
	ret := _m.Called(ctx, lastIP, deletedOk)

	if len(ret) == 0 {
		panic("no return value specified for GetByLastIP")
	}

	var r0 domain.BannedSteamPerson
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, net.IP, bool) (domain.BannedSteamPerson, error)); ok {
		return rf(ctx, lastIP, deletedOk)
	}
	if rf, ok := ret.Get(0).(func(context.Context, net.IP, bool) domain.BannedSteamPerson); ok {
		r0 = rf(ctx, lastIP, deletedOk)
	} else {
		r0 = ret.Get(0).(domain.BannedSteamPerson)
	}

	if rf, ok := ret.Get(1).(func(context.Context, net.IP, bool) error); ok {
		r1 = rf(ctx, lastIP, deletedOk)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBanSteamRepository_GetByLastIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByLastIP'
type MockBanSteamRepository_GetByLastIP_Call struct {
	*mock.Call
}

// GetByLastIP is a helper method to define mock.On call
//   - ctx context.Context
//   - lastIP net.IP
//   - deletedOk bool
func (_e *MockBanSteamRepository_Expecter) GetByLastIP(ctx interface{}, lastIP interface{}, deletedOk interface{}) *MockBanSteamRepository_GetByLastIP_Call {
	return &MockBanSteamRepository_GetByLastIP_Call{Call: _e.mock.On("GetByLastIP", ctx, lastIP, deletedOk)}
}

func (_c *MockBanSteamRepository_GetByLastIP_Call) Run(run func(ctx context.Context, lastIP net.IP, deletedOk bool)) *MockBanSteamRepository_GetByLastIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(net.IP), args[2].(bool))
	})
	return _c
}

func (_c *MockBanSteamRepository_GetByLastIP_Call) Return(_a0 domain.BannedSteamPerson, _a1 error) *MockBanSteamRepository_GetByLastIP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBanSteamRepository_GetByLastIP_Call) RunAndReturn(run func(context.Context, net.IP, bool) (domain.BannedSteamPerson, error)) *MockBanSteamRepository_GetByLastIP_Call {
	_c.Call.Return(run)
	return _c
}

// GetBySteamID provides a mock function with given fields: ctx, sid64, deletedOk
func (_m *MockBanSteamRepository) GetBySteamID(ctx context.Context, sid64 steamid.SID64, deletedOk bool) (domain.BannedSteamPerson, error) {
	ret := _m.Called(ctx, sid64, deletedOk)

	if len(ret) == 0 {
		panic("no return value specified for GetBySteamID")
	}

	var r0 domain.BannedSteamPerson
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, steamid.SID64, bool) (domain.BannedSteamPerson, error)); ok {
		return rf(ctx, sid64, deletedOk)
	}
	if rf, ok := ret.Get(0).(func(context.Context, steamid.SID64, bool) domain.BannedSteamPerson); ok {
		r0 = rf(ctx, sid64, deletedOk)
	} else {
		r0 = ret.Get(0).(domain.BannedSteamPerson)
	}

	if rf, ok := ret.Get(1).(func(context.Context, steamid.SID64, bool) error); ok {
		r1 = rf(ctx, sid64, deletedOk)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBanSteamRepository_GetBySteamID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBySteamID'
type MockBanSteamRepository_GetBySteamID_Call struct {
	*mock.Call
}

// GetBySteamID is a helper method to define mock.On call
//   - ctx context.Context
//   - sid64 steamid.SID64
//   - deletedOk bool
func (_e *MockBanSteamRepository_Expecter) GetBySteamID(ctx interface{}, sid64 interface{}, deletedOk interface{}) *MockBanSteamRepository_GetBySteamID_Call {
	return &MockBanSteamRepository_GetBySteamID_Call{Call: _e.mock.On("GetBySteamID", ctx, sid64, deletedOk)}
}

func (_c *MockBanSteamRepository_GetBySteamID_Call) Run(run func(ctx context.Context, sid64 steamid.SID64, deletedOk bool)) *MockBanSteamRepository_GetBySteamID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(steamid.SID64), args[2].(bool))
	})
	return _c
}

func (_c *MockBanSteamRepository_GetBySteamID_Call) Return(_a0 domain.BannedSteamPerson, _a1 error) *MockBanSteamRepository_GetBySteamID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBanSteamRepository_GetBySteamID_Call) RunAndReturn(run func(context.Context, steamid.SID64, bool) (domain.BannedSteamPerson, error)) *MockBanSteamRepository_GetBySteamID_Call {
	_c.Call.Return(run)
	return _c
}

// GetOlderThan provides a mock function with given fields: ctx, filter, since
func (_m *MockBanSteamRepository) GetOlderThan(ctx context.Context, filter domain.QueryFilter, since time.Time) ([]domain.BanSteam, error) {
	ret := _m.Called(ctx, filter, since)

	if len(ret) == 0 {
		panic("no return value specified for GetOlderThan")
	}

	var r0 []domain.BanSteam
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.QueryFilter, time.Time) ([]domain.BanSteam, error)); ok {
		return rf(ctx, filter, since)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.QueryFilter, time.Time) []domain.BanSteam); ok {
		r0 = rf(ctx, filter, since)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BanSteam)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.QueryFilter, time.Time) error); ok {
		r1 = rf(ctx, filter, since)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBanSteamRepository_GetOlderThan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOlderThan'
type MockBanSteamRepository_GetOlderThan_Call struct {
	*mock.Call
}

// GetOlderThan is a helper method to define mock.On call
//   - ctx context.Context
//   - filter domain.QueryFilter
//   - since time.Time
func (_e *MockBanSteamRepository_Expecter) GetOlderThan(ctx interface{}, filter interface{}, since interface{}) *MockBanSteamRepository_GetOlderThan_Call {
	return &MockBanSteamRepository_GetOlderThan_Call{Call: _e.mock.On("GetOlderThan", ctx, filter, since)}
}

func (_c *MockBanSteamRepository_GetOlderThan_Call) Run(run func(ctx context.Context, filter domain.QueryFilter, since time.Time)) *MockBanSteamRepository_GetOlderThan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.QueryFilter), args[2].(time.Time))
	})
	return _c
}

func (_c *MockBanSteamRepository_GetOlderThan_Call) Return(_a0 []domain.BanSteam, _a1 error) *MockBanSteamRepository_GetOlderThan_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBanSteamRepository_GetOlderThan_Call) RunAndReturn(run func(context.Context, domain.QueryFilter, time.Time) ([]domain.BanSteam, error)) *MockBanSteamRepository_GetOlderThan_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: ctx, ban
func (_m *MockBanSteamRepository) Save(ctx context.Context, ban *domain.BanSteam) error {
	ret := _m.Called(ctx, ban)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.BanSteam) error); ok {
		r0 = rf(ctx, ban)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBanSteamRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockBanSteamRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - ban *domain.BanSteam
func (_e *MockBanSteamRepository_Expecter) Save(ctx interface{}, ban interface{}) *MockBanSteamRepository_Save_Call {
	return &MockBanSteamRepository_Save_Call{Call: _e.mock.On("Save", ctx, ban)}
}

func (_c *MockBanSteamRepository_Save_Call) Run(run func(ctx context.Context, ban *domain.BanSteam)) *MockBanSteamRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.BanSteam))
	})
	return _c
}

func (_c *MockBanSteamRepository_Save_Call) Return(_a0 error) *MockBanSteamRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBanSteamRepository_Save_Call) RunAndReturn(run func(context.Context, *domain.BanSteam) error) *MockBanSteamRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// Stats provides a mock function with given fields: ctx, stats
func (_m *MockBanSteamRepository) Stats(ctx context.Context, stats *domain.Stats) error {
	ret := _m.Called(ctx, stats)

	if len(ret) == 0 {
		panic("no return value specified for Stats")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Stats) error); ok {
		r0 = rf(ctx, stats)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBanSteamRepository_Stats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stats'
type MockBanSteamRepository_Stats_Call struct {
	*mock.Call
}

// Stats is a helper method to define mock.On call
//   - ctx context.Context
//   - stats *domain.Stats
func (_e *MockBanSteamRepository_Expecter) Stats(ctx interface{}, stats interface{}) *MockBanSteamRepository_Stats_Call {
	return &MockBanSteamRepository_Stats_Call{Call: _e.mock.On("Stats", ctx, stats)}
}

func (_c *MockBanSteamRepository_Stats_Call) Run(run func(ctx context.Context, stats *domain.Stats)) *MockBanSteamRepository_Stats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Stats))
	})
	return _c
}

func (_c *MockBanSteamRepository_Stats_Call) Return(_a0 error) *MockBanSteamRepository_Stats_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBanSteamRepository_Stats_Call) RunAndReturn(run func(context.Context, *domain.Stats) error) *MockBanSteamRepository_Stats_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockBanSteamRepository creates a new instance of MockBanSteamRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBanSteamRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBanSteamRepository {
	mock := &MockBanSteamRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}