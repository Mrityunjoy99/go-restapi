package user_test

import (
	"errors"
	"testing"

	"github.com/Mrityunjoy99/sample-go/src/application/user"
	"github.com/Mrityunjoy99/sample-go/src/common/constant"
	"github.com/Mrityunjoy99/sample-go/src/domain/entity"
	"github.com/Mrityunjoy99/sample-go/src/infrastructure/database"
	mock_repository "github.com/Mrityunjoy99/sample-go/src/mocks/repository"
	"github.com/Mrityunjoy99/sample-go/src/tools/genericerror"
	"github.com/Mrityunjoy99/sample-go/src/tools/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	suite.Suite
	service  user.Service
	userRepo *mock_repository.UserRepository
}

func (suite *UsersTestSuite) SetupTest() {
	lggr, _ := logger.NewZapLogger("test")
	suite.userRepo = mock_repository.NewUserRepository(suite.T())
	suite.service = user.NewService(lggr, suite.userRepo)
}

func (suite *UsersTestSuite) TestGetUserByIdSuccess() {
	expectedUser := entity.User{
		BaseModel: database.BaseModel{
			Id: uuid.New(),
		},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon.doe@gmail.com",
		Phone:     "1234567890",
	}

	suite.userRepo.On("GetUserById", expectedUser.Id).Return(&expectedUser, nil)
	actualUser, err := suite.service.GetUserById(expectedUser.Id)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.FirstName, actualUser.FirstName)
	assert.Equal(suite.T(), expectedUser.LastName, actualUser.LastName)
	assert.Equal(suite.T(), expectedUser.Email, actualUser.Email)
	assert.Equal(suite.T(), expectedUser.Phone, actualUser.Phone)
}

func (suite *UsersTestSuite) TestGetUserByIdFailure() {
	userId := uuid.New()

	suite.userRepo.On("GetUserById", userId).
		Return(nil, genericerror.NewGenericError(constant.ErrorCodeBadRequest, "bad request", nil, nil))

	_, err := suite.service.GetUserById(userId)

	assert.NotNil(suite.T(), err)
}

func (suite *UsersTestSuite) TestCreateUserSuccess() {
	expectedUser := entity.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon.doe@gmail.com",
		Phone:     "1234567890",
	}

	createUserDto := user.CreateUserDto{
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Email:     expectedUser.Email,
		Phone:     expectedUser.Phone,
	}

	suite.userRepo.On("CreateUser", expectedUser).Return(&expectedUser, nil)
	actualUser, err := suite.service.CreateUser(createUserDto)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.FirstName, actualUser.FirstName)
	assert.Equal(suite.T(), expectedUser.LastName, actualUser.LastName)
	assert.Equal(suite.T(), expectedUser.Email, actualUser.Email)
	assert.Equal(suite.T(), expectedUser.Phone, actualUser.Phone)
}

func (suite *UsersTestSuite) TestCreateUserFailure() {
	expectedUser := entity.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon.doe@gmail.com",
		Phone:     "1234567890",
	}

	createUserDto := user.CreateUserDto{
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Email:     expectedUser.Email,
		Phone:     expectedUser.Phone,
	}

	suite.userRepo.On("CreateUser", expectedUser).
		Return(nil, genericerror.NewInternalErrByErr(errors.New("failed to create user")))

	_, err := suite.service.CreateUser(createUserDto)

	assert.NotNil(suite.T(), err)
}

func (suite *UsersTestSuite) TestUpdateUserSuccess() {
	expectedUser := entity.User{
		BaseModel: database.BaseModel{
			Id: uuid.New(),
		},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon.doe@gmail.com",
		Phone:     "1234567890",
	}

	updateUserDto := user.UpdateUserDto{
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Email:     expectedUser.Email,
		Phone:     expectedUser.Phone,
	}

	suite.userRepo.On("UpdateUser", expectedUser).Return(&expectedUser, nil)
	actualUser, err := suite.service.UpdateUser(expectedUser.Id, updateUserDto)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.FirstName, actualUser.FirstName)
}

func (suite *UsersTestSuite) TestUpdateUserFailureUserNotFound() {
	expectedUser := entity.User{
		BaseModel: database.BaseModel{
			Id: uuid.New(),
		},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon.doe@gmail.com",
		Phone:     "1234567890",
	}

	updateUserDto := user.UpdateUserDto{
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Email:     expectedUser.Email,
		Phone:     expectedUser.Phone,
	}

	suite.userRepo.On("UpdateUser", expectedUser).
		Return(nil, genericerror.NewGenericError(constant.ErrorCodeResourceNotFound, "user not found", nil, nil))

	actualUser, err := suite.service.UpdateUser(expectedUser.Id, updateUserDto)

	assert.Equal(suite.T(), actualUser.FirstName, "")
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.GetCode(), constant.ErrorCodeResourceNotFound)
}

func (suite *UsersTestSuite) TestUpdateUserFailureGeneral() {
	expectedUser := entity.User{
		BaseModel: database.BaseModel{
			Id: uuid.New(),
		},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jhon.doe@gmail.com",
		Phone:     "1234567890",
	}

	updateUserDto := user.UpdateUserDto{
		FirstName: expectedUser.FirstName,
		LastName:  expectedUser.LastName,
		Email:     expectedUser.Email,
		Phone:     expectedUser.Phone,
	}

	suite.userRepo.On("UpdateUser", expectedUser).
		Return(nil, genericerror.NewGenericError(constant.ErrorCodeInternalServerError, "network failure", nil, nil))

	actualUser, err := suite.service.UpdateUser(expectedUser.Id, updateUserDto)

	assert.Equal(suite.T(), actualUser.FirstName, "")
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.GetCode(), constant.ErrorCodeInternalServerError)
}

func (suite *UsersTestSuite) TestDeleteUserSuccess() {
	userId := uuid.New()

	suite.userRepo.On("DeleteUser", userId).Return(nil)

	err := suite.service.DeleteUser(userId)

	assert.Nil(suite.T(), err)
}

func (suite *UsersTestSuite) TestDeleteUserFailureNotFound() {
	userId := uuid.New()

	suite.userRepo.On("DeleteUser", userId).
		Return(genericerror.NewGenericError(constant.ErrorCodeBadRequest, "user not found", nil, nil))

	err := suite.service.DeleteUser(userId)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.GetCode(), constant.ErrorCodeBadRequest)
}

func (suite *UsersTestSuite) TestDeleteUserFailureGeneral() {
	userId := uuid.New()

	suite.userRepo.On("DeleteUser", userId).
		Return(genericerror.NewGenericError(constant.ErrorCodeInternalServerError, "network failure", nil, nil))

	err := suite.service.DeleteUser(userId)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), err.GetCode(), constant.ErrorCodeInternalServerError)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
