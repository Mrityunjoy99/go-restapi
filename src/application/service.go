package application

import (
	"errors"

	"github.com/Mrityunjoy99/sample-go/src/application/admin"
	"github.com/Mrityunjoy99/sample-go/src/application/employee"
	"github.com/Mrityunjoy99/sample-go/src/application/user"
	"github.com/Mrityunjoy99/sample-go/src/common/config"
	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/Mrityunjoy99/sample-go/src/domain/service"
	"github.com/Mrityunjoy99/sample-go/src/repository"
	"github.com/Mrityunjoy99/sample-go/src/tools/genericerror"
)

type Service struct {
	UserService  user.Service
	AdminService admin.Service
	EmployeeService employee.Service
}

func NewService(c *config.Config, r *repository.Repository, domainService *service.ServiceRegistry) (*Service, genericerror.GenericError) {
	if domainService == nil {
		return nil, genericerror.NewGenericError(constant.ErrorCodeBadRequest, "domainService is required", nil, errors.New("domainService is required"))
	}

	if domainService.JwtService == nil {
		return nil, genericerror.NewGenericError(constant.ErrorCodeBadRequest, "JwtService is required in domainService", nil, errors.New("JwtService is required in domainService"))
	}

	userService := user.NewService(r.UserRepo)
	adminService := admin.NewService(domainService.JwtService, c.Jwt.ExpireTimeSec)
	employeeService := employee.NewService(r.EmployeeRepo)

	return &Service{
		UserService:  userService,
		AdminService: adminService,
		EmployeeService: employeeService,
	}, nil
}
