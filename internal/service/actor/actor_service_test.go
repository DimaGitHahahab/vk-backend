package actor

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"vk-backend/internal/domain"
	"vk-backend/mocks"
)

func TestActorService_AddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	birthDate := time.Now()
	repo.
		EXPECT().
		AddActor(gomock.Any(), "name", 1, birthDate).
		Return(&domain.Actor{
			Id:        1,
			Name:      "name",
			Gender:    1,
			BirthDate: birthDate,
		}, nil)

	act, err := service.AddActor(context.Background(), "name", 1, birthDate)
	assert.NoError(t, err)
	assert.Equal(t, &domain.Actor{
		Id:        1,
		Name:      "name",
		Gender:    1,
		BirthDate: birthDate,
	}, act)
}

func TestActorService_AddActor_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	birthDate := time.Now()
	act, err := service.AddActor(context.Background(), "", 1, birthDate)
	assert.ErrorIs(t, err, domain.ErrEmptyName)
	assert.Nil(t, act)

	act, err = service.AddActor(context.Background(), "name", 1, time.Time{})
	assert.ErrorIs(t, err, domain.ErrEmptyBirthDate)
	assert.Nil(t, act)

	act, err = service.AddActor(context.Background(), "name", 1, time.Now().Add(time.Hour))
	assert.ErrorIs(t, err, domain.ErrFutureBirthDate)
	assert.Nil(t, act)
}

func TestActorService_GetActorById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(true, nil)

	birthDate := time.Now()
	repo.
		EXPECT().
		GetActorById(gomock.Any(), 1).
		Return(&domain.Actor{
			Id:        1,
			Name:      "name",
			Gender:    1,
			BirthDate: birthDate,
		}, nil)

	act, err := service.GetActorById(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, &domain.Actor{
		Id:        1,
		Name:      "name",
		Gender:    1,
		BirthDate: birthDate,
	}, act)
}

func TestActorService_GetActorById_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(false, nil)

	act, err := service.GetActorById(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrActorNotExists)
	assert.Nil(t, act)
}

func TestActorService_UpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)
	birthDate := time.Now()
	repo.
		EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		UpdateActor(gomock.Any(), &domain.Actor{
			Id:        1,
			Name:      "name",
			Gender:    1,
			BirthDate: birthDate,
		}).
		Return(nil)

	err := service.UpdateActor(context.Background(), &domain.Actor{
		Id:        1,
		Name:      "name",
		Gender:    1,
		BirthDate: birthDate,
	})

	assert.NoError(t, err)

}

func TestActorService_UpdateActor_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)
	birthDate := time.Now()
	repo.
		EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(false, nil)

	err := service.UpdateActor(context.Background(), &domain.Actor{
		Id:        1,
		Name:      "name",
		Gender:    1,
		BirthDate: birthDate,
	})

	assert.ErrorIs(t, err, domain.ErrActorNotExists)
}

func TestActorService_DeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		DeleteActor(gomock.Any(), 1).
		Return(nil)

	err := service.DeleteActor(context.Background(), 1)
	assert.NoError(t, err)
}

func TestActorService_DeleteActor_NotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	repo.
		EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(false, nil)

	err := service.DeleteActor(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrActorNotExists)
}

func TestActorService_RepoReturnsInnerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockActorRepository(ctrl)
	service := NewService(repo)

	birthDate := time.Now()

	repo.
		EXPECT().
		AddActor(gomock.Any(), "name", 1, birthDate).
		Return(nil, assert.AnError)

	act, err := service.AddActor(context.Background(), "name", 1, birthDate)
	assert.ErrorIs(t, err, assert.AnError)
	assert.Nil(t, act)

	repo.EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		GetActorById(gomock.Any(), 1).
		Return(nil, assert.AnError)

	act, err = service.GetActorById(context.Background(), 1)
	assert.ErrorIs(t, err, assert.AnError)
	assert.Nil(t, act)

	repo.EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		UpdateActor(gomock.Any(), &domain.Actor{
			Id:        1,
			Name:      "name",
			Gender:    1,
			BirthDate: birthDate,
		}).
		Return(assert.AnError)

	err = service.UpdateActor(context.Background(), &domain.Actor{
		Id:        1,
		Name:      "name",
		Gender:    1,
		BirthDate: birthDate,
	})
	assert.ErrorIs(t, err, assert.AnError)

	repo.EXPECT().
		ActorExists(gomock.Any(), 1).
		Return(true, nil)

	repo.
		EXPECT().
		DeleteActor(gomock.Any(), 1).
		Return(assert.AnError)

	err = service.DeleteActor(context.Background(), 1)
	assert.ErrorIs(t, err, assert.AnError)
}
