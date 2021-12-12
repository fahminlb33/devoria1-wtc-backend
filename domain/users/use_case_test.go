package users_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/users"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/fahminlb33/devoria1-wtc-backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// --- Login

type LoginUseCaseSuite struct {
	suite.Suite
}

func (s *LoginUseCaseSuite) SetupSuite() {
	publicKey, privateKey := mocks.GenerateRSAKeyPair()
	config.GlobalConfig.Authentication.PublicKey = publicKey
	config.GlobalConfig.Authentication.PrivateKey = privateKey
}

func (s *LoginUseCaseSuite) TestUCLoginUserNotFound() {
	db, gormdb, mock := mocks.SetupGormMock(s.T())
	defer db.Close()

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := users.ConstructUserUseCase(gormdb, jwtAuth)

	// set queries
	mock.ExpectQuery("SELECT (.+)").WillReturnError(gorm.ErrRecordNotFound)

	// create model
	model := users.LoginModel{
		Email:    "foo@bar.com",
		Password: "foobar",
	}

	// call function
	response := useCase.Login(context.Background(), model)

	// assert
	assert.Equal(s.T(), http.StatusBadRequest, response.HttpStatus)
	assert.Equal(s.T(), "User not found", response.Message)
}

func (s *LoginUseCaseSuite) TestUCLoginPasswordMismatch() {
	db, gormdb, mock := mocks.SetupGormMock(s.T())
	defer db.Close()

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := users.ConstructUserUseCase(gormdb, jwtAuth)

	// set queries
	rows := mock.NewRows([]string{"email", "password"}).AddRow("foo@bar.com", "$2a$12$5hin/ijVdB1KQew5CtBme.Xxxn/GIjJayTO5QaBA385jgybzFSqKy")
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	// create model
	model := users.LoginModel{
		Email:    "foo@bar.com",
		Password: "foobar2",
	}

	// call function
	response := useCase.Login(context.Background(), model)

	// assert
	assert.Equal(s.T(), http.StatusBadRequest, response.HttpStatus)
	assert.Equal(s.T(), "Invalid password", response.Message)
}

func (s *LoginUseCaseSuite) TestUCLoginSuccess() {
	db, gormdb, mock := mocks.SetupGormMock(s.T())
	defer db.Close()

	// construct handler
	jwtAuth, _ := authentication.ConstructJwtAuth()
	useCase := users.ConstructUserUseCase(gormdb, jwtAuth)

	// set queries
	rows := mock.NewRows([]string{"email", "password"}).AddRow("foo@bar.com", "$2a$12$5hin/ijVdB1KQew5CtBme.Xxxn/GIjJayTO5QaBA385jgybzFSqKy")
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(rows)

	// create model
	model := users.LoginModel{
		Email:    "foo@bar.com",
		Password: "foobar",
	}

	// call function
	response := useCase.Login(context.Background(), model)

	// assert
	assert.Equal(s.T(), "Hello!", response.Message)
	assert.NotNil(s.T(), response.Data)
}

func TestLoginUseCaseSuite(t *testing.T) {
	suite.Run(t, new(LoginUseCaseSuite))
}
