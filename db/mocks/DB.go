// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import db "github.com/prashanthpai/contactbook/db"
import mock "github.com/stretchr/testify/mock"

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

// All provides a mock function with given fields: _a0
func (_m *DB) All(_a0 int) ([]*db.Entry, error) {
	ret := _m.Called(_a0)

	var r0 []*db.Entry
	if rf, ok := ret.Get(0).(func(int) []*db.Entry); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*db.Entry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *DB) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: _a0
func (_m *DB) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByEmail provides a mock function with given fields: _a0
func (_m *DB) FindByEmail(_a0 string) (*db.Entry, error) {
	ret := _m.Called(_a0)

	var r0 *db.Entry
	if rf, ok := ret.Get(0).(func(string) *db.Entry); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Entry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: _a0, _a1
func (_m *DB) FindByName(_a0 string, _a1 int) ([]*db.Entry, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*db.Entry
	if rf, ok := ret.Get(0).(func(string, int) []*db.Entry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*db.Entry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0
func (_m *DB) Store(_a0 *db.Entry) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*db.Entry) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *DB) Update(_a0 string, _a1 *db.Entry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *db.Entry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}