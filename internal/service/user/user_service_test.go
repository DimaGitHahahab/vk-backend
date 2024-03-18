package user

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"vk-backend/internal/domain"
	"vk-backend/mocks"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)

	repo.EXPECT().UserExists(gomock.Any(), "test").Return(false, nil)
	repo.EXPECT().AddUser(gomock.Any(), "test", gomock.Any()).Return(nil, nil)

	s := NewService(repo)
	_, err := s.Register(nil, "test", "test")
	assert.NoError(t, err)
}

func TestUserService_Register_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)

	s := NewService(repo)
	_, err := s.Register(nil, "", "test")
	assert.ErrorIs(t, err, domain.ErrEmptyName)

	_, err = s.Register(nil, "test", "")
	assert.ErrorIs(t, err, domain.ErrEmptyPassword)
}

func TestUserService_Register_UserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)

	repo.EXPECT().UserExists(gomock.Any(), "test").Return(true, nil)

	s := NewService(repo)
	_, err := s.Register(nil, "test", "test")
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
}

func TestUserService_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)

	repo.EXPECT().GetUserById(gomock.Any(), 1).Return(&domain.User{}, nil)

	s := NewService(repo)
	_, err := s.GetUserById(nil, 1)
	assert.NoError(t, err)
}

func TestUserService_GetUserByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)

	repo.EXPECT().UserExists(gomock.Any(), "test").Return(true, nil)
	repo.EXPECT().GetUserByName(gomock.Any(), "test").Return(&domain.User{}, nil)

	s := NewService(repo)
	_, err := s.GetUserByName(nil, "test")
	assert.NoError(t, err)
}
