package service

import (
	"errors"
	"log"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/api/data/types"
	"noda/api/repository"
	"noda/failure"
	"regexp"
	"strings"

	"github.com/google/uuid"
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
	user, err := s.r.Insert(next)
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

func (s *UserService) Update(userID uuid.UUID, up *transfer.UserUpdate) (bool, error) {
	return s.r.Update(userID.String(), up)
}

func (s *UserService) PromoteToAdmin(userID uuid.UUID) (bool, error) {
	return s.r.PromoteToAdmin(userID.String())
}

func (s *UserService) DegradeToNormalUser(userID uuid.UUID) (bool, error) {
	return s.r.DegradeToNormalUser(userID.String())
}

func (s *UserService) GetByEmail(email string) (*transfer.User, error) {
	return s.r.SelectShallowUserByEmail(email)
}

func (s *UserService) GetByID(id uuid.UUID) (*transfer.User, error) {
	return s.r.SelectShallowUserByID(id.String())
}

func (s *UserService) GetUserWithPasswordByEmail(email string) (*model.User, error) {
	return s.r.SelectRawUserByEmail(email)
}

func (s *UserService) GetAll(pag *types.Pagination) (*types.Result[transfer.User], error) {
	users, err := s.r.SelectAll(pag.RPP, pag.Page)
	if err != nil {
		return nil, err
	}
	return &types.Result[transfer.User]{
		Page:      pag.Page,
		RPP:       pag.RPP,
		Retrieved: int64(len(*users)),
		Payload:   users,
	}, nil
}

func (s *UserService) DeleteUserByID(id uuid.UUID) (string, error) {
	return s.r.Delete(id.String())
}
