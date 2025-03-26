package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	"task-manager/internal/models"
	"task-manager/pkg/database"
	"task-manager/pkg/utils"
)

// UserRegisterRequest defines the data needed to register a new user
type UserRegisterRequest struct {
	Username string
	Email    string
	Password string
}

// UserLoginRequest defines the data needed to login a user
type UserLoginRequest struct {
	Email    string
	Password string
}

// AuthResponse represents the authentication response with token and user details
type AuthResponse struct {
	Token string
	User  *models.User
}

// UserService provides methods for user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

// Register creates a new user account
func (s *UserService) Register(req UserRegisterRequest) (*AuthResponse, error) {
	// Check if username already exists
	var existingUser models.User
	result := s.db.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error while checking username: %w", result.Error)
	}

	// Check if email already exists
	result = s.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error while checking email: %w", result.Error)
	}

	// Create new user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // Will be hashed by BeforeSave hook
	}

	// Save user to database
	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate session ID
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  &user,
	}, nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(req UserLoginRequest) (*AuthResponse, error) {
	// Find user by email
	var user models.User
	result := s.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Verify password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate session ID
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  &user,
	}, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID uint, updates map[string]interface{}) (*models.User, error) {
	// Get the user
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// If password is being updated, hash it
	if password, ok := updates["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password"] = string(hashedPassword)
	}

	// Apply updates
	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Refresh user data
	if err := s.db.First(user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve updated user: %w", err)
	}

	return user, nil
}