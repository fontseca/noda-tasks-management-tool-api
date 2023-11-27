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

type UserService interface {
	Save(creation *transfer.UserCreation) (insertedID uuid.UUID, err error)
	FetchByID(id uuid.UUID) (user *transfer.User, err error)
	FetchByEmail(email string) (user *transfer.User, err error)
	FetchRawUserByEmail(email string) (user *model.User, err error)
	Fetch(pagination *types.Pagination) (result *types.Result[transfer.User], err error)
	FetchBlocked(pagination *types.Pagination) (result *types.Result[transfer.User], err error)
	FetchSettings(userID uuid.UUID, pagination *types.Pagination) (result *types.Result[transfer.UserSetting], err error)
	FetchOneSetting(userID uuid.UUID, settingKey string) (setting *transfer.UserSetting, err error)
	Search(pagination *types.Pagination, needle, sortExpr string) (users *types.Result[transfer.User], err error)
	Update(id uuid.UUID, update *transfer.UserUpdate) (ok bool, err error)
	UpdateUserSetting(userID uuid.UUID, settingKey string, update *transfer.UserSettingUpdate) (ok bool, err error)
	Block(id uuid.UUID) (ok bool, err error)
	Unblock(id uuid.UUID) (ok bool, err error)
	PromoteToAdmin(id uuid.UUID) (ok bool, err error)
	DegradeToUser(id uuid.UUID) (ok bool, err error)
	RemoveHardly(id uuid.UUID) error
	RemoveSoftly(id uuid.UUID) error
}

type userService struct {
	r repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository}
}

func (s *userService) Save(creation *transfer.UserCreation) (insertedID uuid.UUID, err error) {
	if nil == creation {
		err = noda.NewNilParameterError("Save", "creation")
		log.Println(err)
		return uuid.Nil, err
	}
	creation.FirstName = strings.Trim(creation.FirstName, " \a\b\f\n\r\t\v")
	creation.MiddleName = strings.Trim(creation.MiddleName, " \a\b\f\n\r\t\v")
	creation.LastName = strings.Trim(creation.LastName, " \a\b\f\n\r\t\v")
	creation.Surname = strings.Trim(creation.Surname, " \a\b\f\n\r\t\v")
	creation.Email = strings.Trim(creation.Email, " \a\b\f\n\r\t\v")
	creation.Password = strings.Trim(creation.Password, " \a\b\f\n\r\t\v")
	switch {
	case 50 < len(creation.FirstName):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("FirstName", "user", 50)
	case 50 < len(creation.MiddleName):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("MiddleName", "user", 50)
	case 50 < len(creation.LastName):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("LastName", "user", 50)
	case 50 < len(creation.Surname):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("Surname", "user", 50)
	case 72 < len(creation.Password):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("Password", "user", 72)
	case 240 < len(creation.Email):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("Email", "user", 240)
	}
	if err := assertPasswordIsValid(&creation.Password, &creation.Email); err != nil {
		return uuid.Nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creation.Password), bcrypt.DefaultCost)
	if err != nil {
		switch {
		default:
			log.Println(err)
			return uuid.Nil, err
		case errors.Is(err, bcrypt.ErrPasswordTooLong):
			return uuid.Nil, noda.ErrPasswordTooLong
		}
	}
	creation.Password = string(hashedPassword)
	insertedIDStr, err := s.r.Save(creation)
	if err != nil {
		return uuid.Nil, err
	}
	parsed, err := uuid.Parse(insertedIDStr)
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

func (s *userService) Update(userID uuid.UUID, up *transfer.UserUpdate) (bool, error) {
	return s.r.Update(userID.String(), up)
}

func (s *userService) PromoteToAdmin(userID uuid.UUID) (bool, error) {
	return s.r.PromoteToAdmin(userID.String())
}

func (s *userService) DegradeToUser(userID uuid.UUID) (bool, error) {
	return s.r.DegradeToUser(userID.String())
}

func (s *userService) Block(userID uuid.UUID) (bool, error) {
	return s.r.Block(userID.String())
}

func (s *userService) Unblock(userID uuid.UUID) (bool, error) {
	return s.r.Unblock(userID.String())
}

func (s *userService) FetchByEmail(email string) (user *transfer.User, err error) {
	email = strings.Trim(email, " \a\b\f\n\r\t\v")
	if "" == email {
		return nil, noda.ErrUserNotFound
	}
	return s.r.FetchShallowUserByEmail(email)
}

func (s *userService) FetchByID(id uuid.UUID) (user *transfer.User, err error) {
	if uuid.Nil == id {
		err = noda.NewNilParameterError("FetchByID", "id")
		log.Println(err)
		return nil, err
	}
	return s.r.FetchShallowUserByID(id.String())
}

func (s *userService) FetchRawUserByEmail(email string) (user *model.User, err error) {
	email = strings.Trim(email, " \a\b\f\n\r\t\v")
	if "" == email {
		return nil, noda.ErrUserNotFound
	}
	return s.r.FetchByEmail(email)
}

func (s *userService) Fetch(pag *types.Pagination) (*types.Result[transfer.User], error) {
	users, err := s.r.Fetch(pag.Page, pag.RPP)
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

func (s *userService) Search(pag *types.Pagination, needle, sortExpr string) (*types.Result[transfer.User], error) {
	users, err := s.r.Search(pag.Page, pag.RPP, needle, sortExpr)
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

func (s *userService) FetchBlocked(pag *types.Pagination) (*types.Result[transfer.User], error) {
	users, err := s.r.FetchBlocked(pag.Page, pag.RPP)
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

func (s *userService) FetchSettings(
	userID uuid.UUID,
	pag *types.Pagination,
) (*types.Result[transfer.UserSetting], error) {
	settings, err := s.r.FetchSettings(userID.String(), pag.Page, pag.RPP)
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

func (s *userService) FetchOneSetting(userID uuid.UUID, settingKey string) (*transfer.UserSetting, error) {
	setting, err := s.r.FetchOneSetting(userID.String(), settingKey)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(setting.Value.([]byte), &setting.Value); err != nil {
		log.Println(err)
		return nil, err
	}
	return setting, nil
}

func (s *userService) UpdateUserSetting(
	userID uuid.UUID,
	settingKey string,
	update *transfer.UserSettingUpdate,
) (bool, error) {
	buf, err := json.Marshal(update.Value)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return s.r.UpdateUserSetting(userID.String(), settingKey, string(buf))
}

func (s *userService) RemoveHardly(id uuid.UUID) error {
	return s.r.RemoveHardly(id.String())
}

func (s *userService) RemoveSoftly(id uuid.UUID) error {
	return s.r.RemoveSoftly(id.String())
}
