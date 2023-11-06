package service

import (
	"encoding/json"
	"errors"
	"log"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/repository"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	r *repository.UserRepository
}

type us = UserService

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository}
}

func (s *us) Save(next *transfer.UserCreation) (uuid.UUID, error) {
	if err := assertPasswordIsValid(&next.Password, &next.Email); err != nil {
		return uuid.Nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(next.Password), bcrypt.DefaultCost)
	if err != nil {
		switch {
		default:
			log.Println(err)
			return uuid.Nil, err
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return uuid.Nil, noda.ErrPasswordTooLong
		}
	}
	next.Password = string(hashedPassword)
	insertedID, err := s.r.InsertUser(next)
	if err != nil {
		return uuid.Nil, err
	}
	parsed, err := uuid.Parse(insertedID)
	if nil != err {
		log.Println(err)
		return uuid.Nil, err
	}
	return parsed, nil
}

func assertPasswordIsValid(password, email *string) *noda.AggregateDetails {
	passwordErrors := new(noda.AggregateDetails)
	emailWithoutAt := strings.Split(*email, "@")[0]
	if strings.Contains(emailWithoutAt, *password) {
		passwordErrors.Append("Password seems to be similar to email.")
		return passwordErrors
	}
	lengthPattern, _ := regexp.Compile(`^.{8,}$`)
	digitPattern, _ := regexp.Compile(`.*\d`)
	upperCasePattern, _ := regexp.Compile(`.*[A-ZÁÉÍÓÚ]`)
	lowerCasePattern, _ := regexp.Compile(`.*[a-záéíóú]`)
	specialCharPattern, _ := regexp.Compile(`.*[!@#$%^&*? ]`)
	if !lengthPattern.MatchString(*password) {
		passwordErrors.Append("Password must be at least 8 characters long.")
	}
	if !digitPattern.MatchString(*password) {
		passwordErrors.Append("Password must contain at least one digit.")
	}
	if !upperCasePattern.MatchString(*password) {
		passwordErrors.Append("Password must contain at least one uppercase letter.")
	}
	if !lowerCasePattern.MatchString(*password) {
		passwordErrors.Append("Password must contain at least one lowercase letter.")
	}
	if !specialCharPattern.MatchString(*password) {
		passwordErrors.Append("Password must contain at least one special character (!@#$%^&*?).")
	}
	if passwordErrors.Has() {
		return passwordErrors
	}
	return nil
}

func (s *us) Update(userID uuid.UUID, up *transfer.UserUpdate) (bool, error) {
	return s.r.UpdateUser(userID.String(), up)
}

func (s *us) PromoteToAdmin(userID uuid.UUID) (bool, error) {
	return s.r.PromoteUserToAdmin(userID.String())
}

func (s *us) DegradeToNormalUser(userID uuid.UUID) (bool, error) {
	return s.r.DegradeAdminToNormalUser(userID.String())
}

func (s *us) Block(userID uuid.UUID) (bool, error) {
	return s.r.BlockUser(userID.String())
}

func (s *us) Unblock(userID uuid.UUID) (bool, error) {
	return s.r.UnblockUser(userID.String())
}

func (s *us) GetByEmail(email string) (*transfer.User, error) {
	return s.r.FetchTransferUserByEmail(email)
}

func (s *us) GetByID(id uuid.UUID) (*transfer.User, error) {
	return s.r.FetchTransferUserByID(id.String())
}

func (s *us) GetUserWithPasswordByEmail(email string) (*model.User, error) {
	return s.r.FetchUserByEmail(email)
}

func (s *us) GetAll(pag *types.Pagination) (*types.Result[transfer.User], error) {
	users, err := s.r.FetchUsers(pag.Page, pag.RPP)
	if err != nil {
		return nil, err
	}
	return &types.Result[transfer.User]{
		Page:      pag.Page,
		RPP:       pag.RPP,
		Retrieved: int64(len(users)),
		Payload:   users,
	}, nil
}

func (s *us) SearchUsers(pag *types.Pagination, needle, sortExpr string) (*types.Result[transfer.User], error) {
	users, err := s.r.SearchUsers(pag.Page, pag.RPP, needle, sortExpr)
	if err != nil {
		return nil, err
	}
	return &types.Result[transfer.User]{
		Page:      pag.Page,
		RPP:       pag.RPP,
		Retrieved: int64(len(users)),
		Payload:   users,
	}, nil
}

func (s *us) GetAllBlocked(pag *types.Pagination) (*types.Result[transfer.User], error) {
	users, err := s.r.FetchBlockedUsers(pag.Page, pag.RPP)
	if err != nil {
		return nil, err
	}
	return &types.Result[transfer.User]{
		Page:      pag.Page,
		RPP:       pag.RPP,
		Retrieved: int64(len(users)),
		Payload:   users,
	}, nil
}

func (s *us) GetUserSettings(pag *types.Pagination, userID uuid.UUID) (*types.Result[transfer.UserSetting], error) {
	settings, err := s.r.FetchUserSettings(userID.String(), pag.Page, pag.RPP)
	if err != nil {
		return nil, err
	}
	for _, setting := range settings {
		if err := json.Unmarshal(setting.Value.([]byte), &setting.Value); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return &types.Result[transfer.UserSetting]{
		Page:      pag.Page,
		RPP:       pag.RPP,
		Retrieved: int64(len(settings)),
		Payload:   settings,
	}, nil
}

func (s *us) GetOneSetting(userID uuid.UUID, settingKey string) (*transfer.UserSetting, error) {
	setting, err := s.r.FetchOneUserSetting(userID.String(), settingKey)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(setting.Value.([]byte), &setting.Value); err != nil {
		log.Println(err)
		return nil, err
	}
	return setting, nil
}

func (s *us) UpdateUserSetting(userID uuid.UUID, settingKey string, update *transfer.UserSettingUpdate) (bool, error) {
	buf, err := json.Marshal(update.Value)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return s.r.UpdateUserSetting(userID.String(), settingKey, string(buf))
}

func (s *us) HardDelete(id uuid.UUID) error {
	return s.r.HardlyDeleteUser(id.String())
}

func (s *us) SoftDelete(id uuid.UUID) (string, error) {
	return s.r.SoftlyDeleteUser(id.String())
}
