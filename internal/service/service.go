package service

import (
	"github.com/go-nunu/nunu-layout-mcp/pkg/log"
	"github.com/go-nunu/nunu-layout-mcp/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	//jwt    *jwt.JWT
	//tm     repository.Transaction
}

func NewService(
	//tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	// jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		//jwt:    jwt,
		//tm:     tm,
	}
}
