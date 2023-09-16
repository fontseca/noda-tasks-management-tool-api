package service

import (
	"errors"
	"log"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/api/repository"
	"noda/failure"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	r *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) Save(next *transfer.UserCreation) (*transfer.User, error) {
	if err := assertPasswordIsValid(&next.Password, &next.Email); err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(next.Password), bcrypt.DefaultCost)
	if err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return nil, failure.ErrPassordTooLong
		}
	}
	next.Password = string(hashedPassword)
	user, err := s.r.Inset(next)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func assertPasswordIsValid(password, email *string) *failure.Aggregation {
	passwordErrors := failure.NewAggregation()
	emailWithoutAt := strings.Split(*email, "@")[0]
	if strings.Contains(emailWithoutAt, *password) {
		passwordErrors.Append(errors.New("password looks similar to email"))
		return passwordErrors
	}
	lengthPattern, _ := regexp.Compile(`^.{8,}$`)
	digitPattern, _ := regexp.Compile(`.*\d`)
	upperCasePattern, _ := regexp.Compile(`.*[A-ZÁÉÍÓÚ]`)
	lowerCasePattern, _ := regexp.Compile(`.*[a-záéíóú]`)
	specialCharPattern, _ := regexp.Compile(`.*[!@#$%^&*? ]`)
	if !lengthPattern.MatchString(*password) {
		passwordErrors.Append(errors.New("password must be at least 8 characters long"))
	}
	if !digitPattern.MatchString(*password) {
		passwordErrors.Append(errors.New("password must contain at least one digit"))
	}
	if !upperCasePattern.MatchString(*password) {
		passwordErrors.Append(errors.New("password must contain at least one uppercase letter"))
	}
	if !lowerCasePattern.MatchString(*password) {
		passwordErrors.Append(errors.New("password must contain at least one lowercase letter"))
	}
	if !specialCharPattern.MatchString(*password) {
		passwordErrors.Append(errors.New("password must contain at least one special character (!@#$%^&*?)"))
	}
	if passwordErrors.Has() {
		return passwordErrors
	}
	return nil
}

func (s *UserService) GetByEmail(email string) (*transfer.User, error) {
	return s.r.SelectByEmail(email)
}

func (s *UserService) GetUserWithPasswordByEmail(email string) (*model.User, error) {
	return s.r.SelectWithPasswordByEmail(email)
}

func (s *UserService) GetAll() (*[]*transfer.User, error) {
	return s.r.SelectAll()
}
