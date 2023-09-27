package service

import (
	"errors"
	"fmt"
	"log"
	"noda/api/data/transfer"
	"noda/failure"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	userService *UserService
}

func NewAuthenticationService(userService *UserService) *AuthenticationService {
	return &AuthenticationService{
		userService: userService,
	}
}

func (s *AuthenticationService) SignUp(next *transfer.UserCreation) (*transfer.User, error) {
	return s.userService.Save(next)
}

func (s *AuthenticationService) SignIn(credentials *transfer.UserCredentials) (*map[string]any, error) {
	user, err := s.userService.GetUserWithPasswordByEmail(credentials.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		switch {
		default:
			log.Println(err)
			return nil, err
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, failure.ErrIncorrectPassord
		}
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
