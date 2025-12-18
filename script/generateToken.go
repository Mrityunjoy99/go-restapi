package main

import (
	"time"

	"github.com/Mrityunjoy99/sample-go/src/common/config"
	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/Mrityunjoy99/sample-go/src/domain/entity"
	"github.com/Mrityunjoy99/sample-go/src/domain/service"
	"github.com/google/uuid"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	domainServiceRegistry, gerr := service.NewServiceRegistry(config)
	if gerr != nil {
		panic("Failed to create domain service registry: " + gerr.Error())
	}

	userId := uuid.NewString()
	userType := constant.UserTypeAdmin

	jwtTokenEntiry := &entity.JwtToken{
		UserId:    userId,
		UserType:  userType,
		ExpiredAt: time.Now().Add(time.Duration(24*365) * time.Hour),
	}

	token, err := domainServiceRegistry.JwtService.GenerateToken(jwtTokenEntiry)
	if err != nil {
		panic("Failed to generate token: " + err.Error())
	}

	println("Generated token:", token)
}
