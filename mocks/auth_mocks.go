package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"noda/data/transfer"
	"noda/data/types"
)

type AuthenticationServiceMock struct {
	mock.Mock
}

func NewAuthenticationServiceMock() *AuthenticationServiceMock {
	return new(AuthenticationServiceMock)
}

func (m *AuthenticationServiceMock) SignUp(creation *transfer.UserCreation) (insertedID uuid.UUID, err error) {
	var args = m.Called(creation)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *AuthenticationServiceMock) SignIn(credentials *transfer.UserCredentials) (payload *types.TokenPayload, err error) {
	var args = m.Called(credentials)
	var arg0 = args.Get(0)
	if nil != arg0 {
		payload = arg0.(*types.TokenPayload)
	}
	return payload, args.Error(1)
}
