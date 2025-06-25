package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepo     UserRepository
	EmployeeRepo EmployeeRepository
}

func NewRepository(db *gorm.DB) *Repository {
	userRepo := NewUserRepository(db)
	employeeRepo := NewEmployeeRepository(db)

	return &Repository{
		UserRepo:     userRepo,
		EmployeeRepo: employeeRepo,
	}
}
