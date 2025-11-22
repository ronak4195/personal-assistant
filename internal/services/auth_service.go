package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(ctx context.Context, name, email, password string) (*models.User, string, error)
	Login(ctx context.Context, email, password string) (*models.User, string, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
}

type authService struct {
	users     repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		users:     userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Signup(ctx context.Context, name, email, password string) (*models.User, string, error) {
	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	u := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if u == nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	return u, token, nil
}

func (s *authService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.users.FindByID(ctx, id)
}

func (s *authService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
