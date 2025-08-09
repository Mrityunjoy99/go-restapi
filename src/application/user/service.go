package user

import (
	"fmt"

	"github.com/Mrityunjoy99/sample-go/src/repository"
	"github.com/Mrityunjoy99/sample-go/src/tools/genericerror"
	"github.com/Mrityunjoy99/sample-go/src/tools/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type service struct {
	logger   logger.Logger
	userRepo repository.UserRepository
}

type Service interface {
	GetUserById(id uuid.UUID) (UserResponseDto, genericerror.GenericError)
	CreateUser(dto CreateUserDto) (UserResponseDto, genericerror.GenericError)
	UpdateUser(id uuid.UUID, dto UpdateUserDto) (UserResponseDto, genericerror.GenericError)
	DeleteUser(id uuid.UUID) genericerror.GenericError
}

func NewService(logger logger.Logger, userRepo repository.UserRepository) Service {
	return &service{logger: logger, userRepo: userRepo}
}

func (s *service) GetUserById(id uuid.UUID) (UserResponseDto, genericerror.GenericError) {
	s.logger.Info("GetUserById. id: " + id.String())
	user, gerr := s.userRepo.GetUserById(id)
	if gerr != nil {
		return UserResponseDto{}, gerr
	}

	return newDtoFromEntity(*user), nil
}

func (s *service) CreateUser(dto CreateUserDto) (UserResponseDto, genericerror.GenericError) {
	s.logger.Info("CreateUser called", zap.Any("dto", dto))
	user := dto.toUserEntity()

	createdUser, gerr := s.userRepo.CreateUser(user)
	if gerr != nil {
		return UserResponseDto{}, gerr
	}

	return newDtoFromEntity(*createdUser), nil
}

func (s *service) UpdateUser(id uuid.UUID, dto UpdateUserDto) (UserResponseDto, genericerror.GenericError) {
	s.logger.Info("UpdateUser. id: " + id.String() + ", dto: " + fmt.Sprint(dto))
	user := dto.toUserEntity()
	user.Id = id

	updatedUser, gerr := s.userRepo.UpdateUser(user)
	if gerr != nil {
		return UserResponseDto{}, gerr
	}

	return newDtoFromEntity(*updatedUser), nil
}

func (s *service) DeleteUser(id uuid.UUID) genericerror.GenericError {
	s.logger.Info("DeleteUser. id: " + id.String())
	return s.userRepo.DeleteUser(id)
}
