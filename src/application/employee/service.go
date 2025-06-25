package employee

import (
	"github.com/Mrityunjoy99/sample-go/src/repository"
	"github.com/Mrityunjoy99/sample-go/src/tools/genericerror"
	"github.com/google/uuid"
)

type service struct {
	employeeRepo repository.EmployeeRepository
}

type Service interface {
	GetEmployeeById(id uuid.UUID) (*EmployeeResponseDto, genericerror.GenericError)
}

func NewService(employeeRepo repository.EmployeeRepository) Service {
	return &service{
		employeeRepo: employeeRepo,
	}
}

func (s service) GetEmployeeById(id uuid.UUID) (*EmployeeResponseDto, genericerror.GenericError) {
	employee, gerr := s.employeeRepo.GetEmployeeById(id)
	if gerr != nil {
		return nil, gerr
	}

	dto := EmployeeResponseDto{
		Id:        id.String(),
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
	}

	return &dto, nil
}
