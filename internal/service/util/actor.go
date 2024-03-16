package util

import (
	"time"
	"vk-backend/internal/domain"
)

func ValidateActorData(name string, birthDate time.Time) error {
	if name == "" {
		return domain.ErrEmptyName
	}
	if birthDate.After(time.Now()) {
		return domain.ErrFutureBirthDate
	}
	if birthDate.IsZero() {
		return domain.ErrEmptyBirthDate
	}
	return nil
}
