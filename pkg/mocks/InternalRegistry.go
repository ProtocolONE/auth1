// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import database "github.com/ProtocolONE/auth1.protocol.one/pkg/database"
import mock "github.com/stretchr/testify/mock"
import persist "github.com/ProtocolONE/auth1.protocol.one/pkg/persist"
import service "github.com/ProtocolONE/auth1.protocol.one/pkg/service"

// InternalRegistry is an autogenerated mock type for the InternalRegistry type
type InternalRegistry struct {
	mock.Mock
}

// ApplicationService provides a mock function with given fields:
func (_m *InternalRegistry) ApplicationService() service.ApplicationServiceInterface {
	ret := _m.Called()

	var r0 service.ApplicationServiceInterface
	if rf, ok := ret.Get(0).(func() service.ApplicationServiceInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.ApplicationServiceInterface)
		}
	}

	return r0
}

// GeoIpService provides a mock function with given fields:
func (_m *InternalRegistry) GeoIpService() service.GeoIp {
	ret := _m.Called()

	var r0 service.GeoIp
	if rf, ok := ret.Get(0).(func() service.GeoIp); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.GeoIp)
		}
	}

	return r0
}

// HydraAdminApi provides a mock function with given fields:
func (_m *InternalRegistry) HydraAdminApi() service.HydraAdminApi {
	ret := _m.Called()

	var r0 service.HydraAdminApi
	if rf, ok := ret.Get(0).(func() service.HydraAdminApi); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.HydraAdminApi)
		}
	}

	return r0
}

// LauncherServer provides a mock function with given fields:
func (_m *InternalRegistry) LauncherServer() service.LauncherServerService {
	ret := _m.Called()

	var r0 service.LauncherServerService
	if rf, ok := ret.Get(0).(func() service.LauncherServerService); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.LauncherServerService)
		}
	}

	return r0
}

// LauncherTokenService provides a mock function with given fields:
func (_m *InternalRegistry) LauncherTokenService() service.LauncherTokenServiceInterface {
	ret := _m.Called()

	var r0 service.LauncherTokenServiceInterface
	if rf, ok := ret.Get(0).(func() service.LauncherTokenServiceInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.LauncherTokenServiceInterface)
		}
	}

	return r0
}

// Mailer provides a mock function with given fields:
func (_m *InternalRegistry) Mailer() service.MailerInterface {
	ret := _m.Called()

	var r0 service.MailerInterface
	if rf, ok := ret.Get(0).(func() service.MailerInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.MailerInterface)
		}
	}

	return r0
}

// MfaService provides a mock function with given fields:
func (_m *InternalRegistry) MfaService() service.MfaApiInterface {
	ret := _m.Called()

	var r0 service.MfaApiInterface
	if rf, ok := ret.Get(0).(func() service.MfaApiInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.MfaApiInterface)
		}
	}

	return r0
}

// MgoSession provides a mock function with given fields:
func (_m *InternalRegistry) MgoSession() database.MgoSession {
	ret := _m.Called()

	var r0 database.MgoSession
	if rf, ok := ret.Get(0).(func() database.MgoSession); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(database.MgoSession)
		}
	}

	return r0
}

// OneTimeTokenService provides a mock function with given fields:
func (_m *InternalRegistry) OneTimeTokenService() service.OneTimeTokenServiceInterface {
	ret := _m.Called()

	var r0 service.OneTimeTokenServiceInterface
	if rf, ok := ret.Get(0).(func() service.OneTimeTokenServiceInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.OneTimeTokenServiceInterface)
		}
	}

	return r0
}

// Watcher provides a mock function with given fields:
func (_m *InternalRegistry) Watcher() persist.Watcher {
	ret := _m.Called()

	var r0 persist.Watcher
	if rf, ok := ret.Get(0).(func() persist.Watcher); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(persist.Watcher)
		}
	}

	return r0
}
