package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"noda"
	"noda/data/transfer"
	"noda/mocks"
	"testing"
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
