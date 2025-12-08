package service

import (
	"errors"
	"time"

	"shosha-finance/internal/models"
	"shosha-finance/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user is not active")
	ErrUserNotFound       = errors.New("user not found")
)

type JWTClaims struct {
	UserID   string          `json:"user_id"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(identifier, password string) (*models.User, string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	CreateDefaultUsers() error
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: 24 * time.Hour,
	}
}

func (s *authService) Login(identifier, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByIdentifier(identifier)
	if err != nil {
		log.Error().Err(err).Str("identifier", identifier).Msg("User not found")
		return nil, "", ErrInvalidCredentials
	}

	log.Info().Str("username", user.Username).Str("role", string(user.Role)).Msg("User found")

	if !user.IsActive {
		log.Warn().Str("username", user.Username).Msg("User not active")
		return nil, "", ErrUserNotActive
	}

	if !user.CheckPassword(password) {
		log.Warn().Str("username", user.Username).Msg("Invalid password")
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, "", err
	}

	log.Info().Str("username", user.Username).Msg("Login successful")
	return user, token, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *authService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *authService) CreateDefaultUsers() error {
	count, err := s.userRepo.Count()
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	defaultUsers := []struct {
		username string
		password string
		name     string
		role     models.UserRole
	}{
		{"admin", "admin123", "Administrator", models.RoleAdmin},
		{"adminShosha", "password123*", "Admin Shosha", models.RoleAdmin},
		{"adminCabang", "password123*", "Admin Cabang", models.RoleAdmin},
		{"officialShosha", "password123*", "Official Shosha", models.RoleManager},
		{"officialCabang", "password123*", "Official Cabang", models.RoleManager},
	}

	for _, u := range defaultUsers {
		user := &models.User{
			Username: u.username,
			Name:     u.name,
			Role:     u.role,
			IsActive: true,
		}
		if err := user.SetPassword(u.password); err != nil {
			return err
		}
		if err := s.userRepo.Create(user); err != nil {
			return err
		}
	}

	return nil
}
