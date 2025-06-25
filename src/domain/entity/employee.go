package entity

import "github.com/Mrityunjoy99/sample-go/src/infrastructure/database"

type Employee struct {
	database.BaseModel
	FirstName string
	LastName  string
}

func (Employee) TableName() string {
	return "employees"
}
