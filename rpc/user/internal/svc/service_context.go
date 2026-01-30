package svc

import (
	"zerorequest/pkg"
)

type ServiceContext struct {
	Config pkg.Config
}

func NewServiceContext(c pkg.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
