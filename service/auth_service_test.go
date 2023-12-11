package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/mocks"
	"testing"
	"time"
)

func TestAuthenticationService_SignUp(t *testing.T) {
	defer beQuiet()()
	const routine = "Save"
	var (
		res uuid.UUID
		err error
	)

	t.Run("success", func(t *testing.T) {
		var inserted = uuid.New()
		var creation = &transfer.UserCreation{}
		var s = mocks.NewUserServiceMock()
		s.On(routine, creation).Return(inserted, nil)
		res, err = NewAuthenticationService(s).SignUp(creation)
		assert.Equal(t, inserted, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"creation\" cannot be nil", func(t *testing.T) {
		var s = mocks.NewUserServiceMock()
		s.AssertNotCalled(t, routine)
		res, err = NewAuthenticationService(s).SignUp(nil)
		assert.ErrorContains(t, err, noda.NewNilParameterError("SignUp", "creation").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("got user service error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var creation = &transfer.UserCreation{}
		var s = mocks.NewUserServiceMock()
		s.On(routine, mock.Anything).Return(uuid.Nil, unexpected)
		res, err = NewAuthenticationService(s).SignUp(creation)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestAuthenticationService_SignIn(t *testing.T) {
	const (
		routine  = "FetchRawUserByEmail"
		password = "x@e8[a+*GAUsKBZ!d}>3&"
		email    = "izs16833@zslsz.com"
	)
	var hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var (
		res  *types.TokenPayload
		err  error
		user = &model.User{
			ID:       uuid.New(),
			Email:    email,
			Password: string(hash),
		}
	)

	t.Run("success", func(t *testing.T) {
		var credentials = &transfer.UserCredentials{Email: user.Email, Password: password}
		var us = mocks.NewUserServiceMock()
		us.On(routine, credentials.Email).Return(user, nil)
		res, err = NewAuthenticationService(us).SignIn(credentials)
		assert.NoError(t, err)
		if assert.NotNil(t, res) {
			var token, _ = jwt.Parse(res.Token, func(tk *jwt.Token) (any, error) {
				if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
					t.Fatal(fmt.Errorf("unexpected signing method: %v", tk.Header["alg"]))
				}
				return noda.Secret(), nil
			})
			var claims = token.Claims.(jwt.MapClaims)
			var iat = claims["iat"]
			if assert.NotNil(t, iat, "Missing \"iat\" claim.") {
				assert.Equal(t, jwt.NewNumericDate(time.Unix(int64(iat.(float64)), 0)).Time, res.IssuedAt)
			}
			var exp = claims["exp"]
			if assert.NotNil(t, exp, "Missing \"exp\" claim.") {
				var at = res.Expires.At
				var within = res.Expires.Within
				var unit = res.Expires.Unit
				assert.Equal(t, jwt.NewNumericDate(time.Unix(int64(exp.(float64)), 0)).Time, at)
				if 0 != res.IssuedAt.Compare(time.Time{}) {
					var expiresWithin = at.Sub(res.IssuedAt).Seconds()
					assert.Equal(t, expiresWithin, within)
				} else {
					assert.Equal(t, -1, within)
				}
				assert.Equal(t, "s", unit, "Time unit must be seconds (\"s\").")
			}
			var iss = claims["iss"]
			if assert.Equal(t, "noda", iss) {
				assert.Equal(t, iss, res.Issuer)
			}
			var sub = claims["sub"]
			if assert.Equal(t, "authentication", sub) {
				assert.Equal(t, sub, res.Subject)
			}
			assert.Equal(t, user.ID.String(), claims["user_id"])
			assert.Equal(t, user.Role, types.Role(claims["user_role"].(float64)))
		}
	})

	t.Run("parameter \"credentials\" cannot be nil", func(t *testing.T) {
		var s = mocks.NewUserServiceMock()
		s.AssertNotCalled(t, routine)
		res, err = NewAuthenticationService(s).SignIn(nil)
		assert.ErrorContains(t, err, noda.NewNilParameterError("SignIn", "credentials").Error())
		assert.Nil(t, res)
	})

	t.Run("must trim email and password", func(t *testing.T) {
		var credentials = &transfer.UserCredentials{
			Email:    blankset + email + blankset,
			Password: blankset + password + blankset,
		}
		var s = mocks.NewUserServiceMock()
		s.On(routine, email).Return(user, nil)
		res, err = NewAuthenticationService(s).SignIn(credentials)
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("email must match regexp", func(t *testing.T) {
		var credentials = &transfer.UserCredentials{Email: "wrong"}
		var s = mocks.NewUserServiceMock()
		s.AssertNotCalled(t, routine)
		res, err = NewAuthenticationService(s).SignIn(credentials)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, "Email address does not match regular expression")
	})

	t.Run("satisfies...", func(t *testing.T) {
		var max = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxX"
		var credentials = &transfer.UserCredentials{Password: password}

		max = max + max + max + max + max
		t.Run("240 < creation.Email", func(t *testing.T) {
			credentials.Email = max
			var s = mocks.NewUserServiceMock()
			s.AssertNotCalled(t, routine)
			res, err = NewAuthenticationService(s).SignIn(credentials)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Email", "credentials", 240).Error())
			assert.Nil(t, res)
			credentials.Email = ""
		})

		t.Run("72 < creation.Password", func(t *testing.T) {
			credentials.Password = max + "0*"
			var r = mocks.NewUserServiceMock()
			r.AssertNotCalled(t, routine)
			res, err = NewAuthenticationService(r).SignIn(credentials)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Password", "credentials", 72).Error())
			assert.Nil(t, res)
		})
	})

	t.Run("got user service error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var credentials = &transfer.UserCredentials{Email: email, Password: password}
		var s = mocks.NewUserServiceMock()
		s.On(routine, mock.Anything).Return(nil, unexpected)
		res, err = NewAuthenticationService(s).SignIn(credentials)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}
