package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"moodly/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user *models.User) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == nil {
		return errors.New("password is required")
	}

	existingUser, err := s.repo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(*user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	hashedPasswordStr := string(hashedPassword)
	user.Password = &hashedPasswordStr

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(email string, password string) (string, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return "", errors.New("email is required")
	}

	if password == "" {
		return "", errors.New("password is required")
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if user.Password == nil {
		return "", errors.New("this account uses OAuth login")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(*user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) LoginWithOAuthGoogle(
	email string,
	name string,
	provider string,
	providerAccountID string,
) (string, *models.User, error) {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)
	provider = strings.TrimSpace(provider)
	providerAccountID = strings.TrimSpace(providerAccountID)

	if email == "" {
		return "", nil, errors.New("email is required")
	}

	if providerAccountID == "" {
		return "", nil, errors.New("provider account id is required")
	}

	if provider == "" {
		provider = "google"
	}

	if name == "" {
		name = email
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, err
		}

		user = &models.User{
			Name:     name,
			Email:    email,
			Password: nil,
		}

		if err := s.repo.CreateUser(user); err != nil {
			return "", nil, err
		}
	}

	_, err = s.repo.FindOrCreateOAuthAccount(
		user.ID,
		provider,
		providerAccountID,
	)
	if err != nil {
		return "", nil, err
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
