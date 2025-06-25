package repository

import (
	"errors"

	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/Mrityunjoy99/sample-go/src/domain/entity"
	"github.com/Mrityunjoy99/sample-go/src/tools/genericerror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	GetEmployeeById(id uuid.UUID) (*entity.Employee, genericerror.GenericError)
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (u *employeeRepository) GetEmployeeById(id uuid.UUID) (*entity.Employee, genericerror.GenericError) {
	var employee entity.Employee

	result := u.db.Model(&employee).Where("id=?", id).Scan(&employee)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, genericerror.NewGenericError(constant.ErrorCodeResourceNotFound, result.Error.Error(), nil, result.Error)
	}

	if result.Error != nil {
		return nil, genericerror.NewInternalErrByErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, genericerror.NewGenericError(constant.ErrorCodeResourceNotFound, "No employee found", nil, errors.New("No employee found"))
	}

	return &employee, nil
}
