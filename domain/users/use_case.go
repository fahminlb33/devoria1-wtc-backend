package users

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/utils"
)

type UserUseCase interface {
	Register(c context.Context, model RegisterModel) (resp utils.Response)
	Login(c context.Context, model LoginModel) (resp utils.Response)
	GetProfile(c context.Context, model GetProfileModel) (resp utils.Response)
}

type UserUseCaseImpl struct {
	Database *gorm.DB
	JwtAuth  authentication.IJwtAuth
}

func ConstructUserUseCase(db *gorm.DB, jwtAuth authentication.IJwtAuth) UserUseCase {
	return &UserUseCaseImpl{
		Database: db,
		JwtAuth:  jwtAuth,
	}
}

func (u *UserUseCaseImpl) Login(c context.Context, model LoginModel) (resp utils.Response) {
	db := u.Database.WithContext(c)

	// check if the user is already registered
	var user User
	dbResult := db.First(&user, "email = ?", model.Email)

	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return utils.WrapResponse(http.StatusBadRequest, "User not found", nil)
	}

	// verify password
	if !authentication.VerifyPassword(model.Password, *user.Password) {
		return utils.WrapResponse(http.StatusBadRequest, "Invalid password", nil)
	}

	// generate acccess token
	accessToken, _ := u.JwtAuth.Sign(strconv.Itoa(int(user.ID)), user.FirstName)
	finalResponse := UserLoginDto{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AccessToken: accessToken,
	}

	return utils.WrapResponse(http.StatusCreated, "User registered", finalResponse)
}

func (u *UserUseCaseImpl) Register(c context.Context, model RegisterModel) (response utils.Response) {
	db := u.Database.WithContext(c)

	// check if the user is already registered
	var user User
	dbResult := db.First(&user, "email = ?", model.Email)

	if !errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return utils.WrapResponse(http.StatusBadRequest, "User already registered", nil)
	}

	// save user to database
	hashedPassword, _ := authentication.HashPassword(model.Password)
	user = User{
		Email:     model.Email,
		Password:  &hashedPassword,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Role:      CONTRIBUTOR,
	}

	db.Create(&user)

	// generate acccess token
	accessToken, _ := u.JwtAuth.Sign(strconv.Itoa(int(user.ID)), user.FirstName)
	finalResponse := UserLoginDto{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AccessToken: accessToken,
	}

	return utils.WrapResponse(http.StatusCreated, "User registered", finalResponse)
}

func (u *UserUseCaseImpl) GetProfile(c context.Context, model GetProfileModel) (resp utils.Response) {
	db := u.Database.WithContext(c)

	// get user by ID
	var user User
	db.First(&user, model.ID)

	// generate acccess token
	finalResponse := UserProfileDto{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return utils.WrapResponse(http.StatusOK, "OK", finalResponse)
}
