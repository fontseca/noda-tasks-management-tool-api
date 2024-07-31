package service

import (
	"errors"
	"fmt"
	"log"
	"noda/data/transfer"
	"noda/data/types"
	"noda/failure"
	"noda/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService interface {
	SignUp(creation *transfer.UserCreation) (insertedID uuid.UUID, err error)
	SignIn(credentials *transfer.UserCredentials) (payload *types.TokenPayload, err error)
}

type authenticationService struct {
	userService UserService
}

func NewAuthenticationService(userService UserService) AuthenticationService {
	return &authenticationService{
		userService: userService,
	}
}

func (s *authenticationService) SignUp(creation *transfer.UserCreation) (insertedID uuid.UUID, err error) {
	if nil == creation {
		err = failure.NewNilParameterError("SignUp", "creation")
		log.Println(err)
		return uuid.Nil, err
	}
	return s.userService.Save(creation)
}

func (s *authenticationService) SignIn(credentials *transfer.UserCredentials) (payload *types.TokenPayload, err error) {
	if nil == credentials {
		return nil, failure.NewNilParameterError("SignIn", "credentials")
	}
	doTrim(&credentials.Email, &credentials.Password)
	switch {
	case 72 < len(credentials.Password):
		return nil, failure.ErrTooLong.Clone().FormatDetails("Password", "credentials", 72)
	case 240 < len(credentials.Email):
		return nil, failure.ErrTooLong.Clone().FormatDetails("Email", "credentials", 240)
	case !emailRegexp.MatchString(credentials.Email):
		return nil, failure.ErrBadRequest.
			Clone().
			SetDetails(fmt.Sprintf("Email address does not match regular expression: %q.", emailRegexp.String()))
	}
	user, err := s.userService.FetchRawUserByEmail(credentials.Email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, failure.ErrIncorrectPassword
		}
	}
	var claims = jwt.MapClaims{
		"iss":       "noda",
		"sub":       "authentication",
		"iat":       jwt.NewNumericDate(time.Now()),
		"exp":       jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"user_id":   user.UUID,
		"user_role": user.Role,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := t.SignedString(global.Secret())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var jti, _ = claims["jti"].(string)
	var sub, _ = claims["sub"].(string)
	var iss, _ = claims["iss"].(string)
	payload = &types.TokenPayload{
		ID:      jti,
		Token:   ss,
		Subject: sub,
		Issuer:  iss,
	}
	var ok bool
	var iat, exp *jwt.NumericDate
	iat, ok = claims["iat"].(*jwt.NumericDate)
	if ok {
		payload.IssuedAt = iat.Time
	}
	exp, ok = claims["exp"].(*jwt.NumericDate)
	if ok {
		var expires = types.TokenExpires{
			At:     exp.Time,
			Within: -1,
			Unit:   "s",
		}
		if nil != iat {
			expires.Within = exp.Sub(iat.Time).Seconds()
		}
		payload.Expires = expires
	}
	return payload, nil
}
