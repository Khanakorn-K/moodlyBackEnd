package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"strings"
)

type CustomCauseService struct {
	repo *repositories.CustomCauseRepository
}

func NewCustomCauseService(repo *repositories.CustomCauseRepository) *CustomCauseService {
	return &CustomCauseService{repo: repo}
}

func (s *CustomCauseService) CreateCause(cause *models.CustomCause) error {
	cause.Name = strings.TrimSpace(cause.Name)

	if cause.UserID == 0 {
		return errors.New("user id is required")
	}

	if cause.Name == "" {
		return errors.New("cause name is required")
	}

	return s.repo.Create(cause)
}

func (s *CustomCauseService) GetCauses(userID uint) ([]models.CustomCause, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}

	return s.repo.FindByUserID(userID)
}

func (s *CustomCauseService) UpdateCause(cause *models.CustomCause) error {
	cause.Name = strings.TrimSpace(cause.Name)

	if cause.ID == 0 {
		return errors.New("cause id is required")
	}

	if cause.UserID == 0 {
		return errors.New("user id is required")
	}

	if cause.Name == "" {
		return errors.New("cause name is required")
	}

	return s.repo.Update(cause)
}

func (s *CustomCauseService) DeleteCause(id uint, userID uint) error {
	if id == 0 {
		return errors.New("cause id is required")
	}

	if userID == 0 {
		return errors.New("user id is required")
	}

	return s.repo.Delete(id, userID)
}
