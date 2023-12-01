package service

import (
	"errors"
	"fmt"
	"log"
	"noda"
	"noda/data/transfer"
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
		err = noda.NewNilParameterError("SignUp", "creation")
		log.Println(err)
		return uuid.Nil, err
	}
	return s.userService.Save(creation)
}

func (s *authenticationService) SignIn(credentials *transfer.UserCredentials) (*map[string]any, error) {
	// TODO: Check credentials.Email is a valid email address.
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
			return nil, noda.ErrIncorrectPassword
		}
	}

	if user.IsBlocked {
		return nil, noda.ErrUserBlocked
	}

	claims := jwt.MapClaims{
		/* Registered claims.  */
		"iss": "noda",
		"sub": "authentication",
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),

		/* Public claims.  */
		"user_id":   user.ID,
		"user_role": user.Role,
	}

	secret := []byte("secret")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := t.SignedString(secret)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &map[string]any{
		"jwt": ss,
		"iat": claims["iat"].(*jwt.NumericDate).String(),
		"exp": claims["exp"].(*jwt.NumericDate).String(),
		"jti": claims["jti"],
	}, nil
}
