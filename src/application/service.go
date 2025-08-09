package application

import (
	"errors"

	"github.com/Mrityunjoy99/sample-go/src/application/admin"
	"github.com/Mrityunjoy99/sample-go/src/application/user"
	"github.com/Mrityunjoy99/sample-go/src/common/config"
	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/Mrityunjoy99/sample-go/src/domain/service"
	"github.com/Mrityunjoy99/sample-go/src/repository"
	"github.com/Mrityunjoy99/sample-go/src/tools/genericerror"
	"github.com/Mrityunjoy99/sample-go/src/tools/logger"
)

type Service struct {
	UserService  user.Service
	AdminService admin.Service
}

func NewService(c *config.Config, r *repository.Repository, logger logger.Logger, domainService *service.ServiceRegistry) (*Service, genericerror.GenericError) {
	if domainService == nil {
		return nil, genericerror.NewGenericError(constant.ErrorCodeBadRequest, "domainService is required", nil, errors.New("domainService is required"))
	}

	if domainService.JwtService == nil {
		return nil, genericerror.NewGenericError(constant.ErrorCodeBadRequest, "JwtService is required in domainService", nil, errors.New("JwtService is required in domainService"))
	}

	userService := user.NewService(logger, r.UserRepo)

	return &Service{
		UserService:  userService,
		AdminService: admin.NewService(domainService.JwtService, c.Jwt.ExpireTimeSec),
	}, nil
}
