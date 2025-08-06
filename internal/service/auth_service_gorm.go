package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"exit-debt/internal/models"
)

type AuthServiceGORM struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpiry time.Duration
	debtSharingService *DebtSharingService
}

func NewAuthServiceGORM(db *gorm.DB, jwtSecret string, jwtExpiry string) (*AuthServiceGORM, error) {
	duration, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		return nil, err
	}

	return &AuthServiceGORM{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: duration,
		debtSharingService: NewDebtSharingService(db),
	}, nil
}

func (s *AuthServiceGORM) Register(req *models.CreateUserRequest) (*models.RegisterResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Share debt lists if user email matches any contacts
	var sharingSummary *models.SharingSummary
	if err := s.debtSharingService.ShareDebtListsWithUser(user.ID, user.Email); err != nil {
		// Log the error but don't fail registration
		// TODO: Add proper logging here
	} else {
		// Get sharing summary
		summary, err := s.debtSharingService.GetSharingSummary(user.ID, user.Email)
		if err == nil && summary != nil {
			sharingSummary = &models.SharingSummary{
				ContactsFound:     summary.ContactsFound,
				DebtListsShared:   summary.DebtListsShared,
				TotalAmountShared: summary.TotalAmountShared,
			}
		}
	}

	return &models.RegisterResponse{
		User:    *user,
		SharingSummary: sharingSummary,
	}, nil
}

func (s *AuthServiceGORM) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthServiceGORM) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthServiceGORM) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, err
		}

		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token")
} 