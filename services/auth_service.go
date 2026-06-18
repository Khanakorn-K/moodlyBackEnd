package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"moodly/utils"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo *repositories.AuthRepository
}

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
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

func (s *AuthService) GenerateOAuthState() (string, error) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(stateBytes), nil
}

func (s *AuthService) GetGoogleOAuthURL(state string) string {
	return utils.GetGoogleOAuthConfig().AuthCodeURL(state)
}

func (s *AuthService) FindOrCreateOAuthUser(
	email string,
	name string,
	provider string,
	providerAccountID string,
) (*models.User, error) {
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)
	provider = strings.TrimSpace(provider)
	providerAccountID = strings.TrimSpace(providerAccountID)

	if email == "" {
		return nil, errors.New("email is required")
	}

	if provider == "" {
		return nil, errors.New("provider is required")
	}

	if providerAccountID == "" {
		return nil, errors.New("provider account id is required")
	}

	if name == "" {
		name = email
	}

	account, err := s.repo.FindOAuthAccountWithUser(provider, providerAccountID)
	if err == nil {
		return &account.User, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user, err := s.repo.FindByEmail(email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &models.User{
			Name:     name,
			Email:    email,
			Password: nil,
		}

		if err := s.repo.CreateUser(user); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	oauthAccount := &models.OAuthAccount{
		UserID:            user.ID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
	}

	if err := s.repo.CreateOAuthAccount(oauthAccount); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginWithOAuthGoogle(
	email string,
	name string,
	provider string,
	providerAccountID string,
) (string, *models.User, error) {
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

func (s *AuthService) LoginWithGoogle(code string) (string, error) {
	oauthConfig := utils.GetGoogleOAuthConfig()

	googleToken, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}

	client := oauthConfig.Client(context.Background(), googleToken)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch google user info")
	}

	var googleUser GoogleUserInfo

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", err
	}

	user, err := s.FindOrCreateOAuthUser(
		googleUser.Email,
		googleUser.Name,
		"google",
		googleUser.ID,
	)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
