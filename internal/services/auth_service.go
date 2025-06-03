package services

import (
	"bad_boyes/internal/models"
	"bad_boyes/internal/repository"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrJWTSecretMissing  = errors.New("JWT secret is not configured")
	ErrDatabaseError     = errors.New("database error occurred")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	auditRepo *repository.AuditRepository
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
}

type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Birthday  time.Time `json:"birthday"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAuthService(userRepo *repository.UserRepository, auditRepo *repository.AuditRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) error {
	log.Printf("Attempting to register user with email: %s", req.Email)

	exists, err := s.userRepo.UserExists(req.Email)
	if err != nil {
		log.Printf("Database error while checking user existence: %v", err)
		return ErrDatabaseError
	}
	if exists {
		log.Printf("Registration failed: user with email %s already exists", req.Email)
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return errors.New("error processing password")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Birthday: req.Birthday,
		Roles:    []models.Role{{Name: "user"}}, // Default role
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return ErrDatabaseError
	}

	log.Printf("User registered successfully: %s", req.Email)
	return nil
}

func (s *AuthService) Login(req models.LoginRequest) (*LoginResponse, error) {
	log.Printf("Attempting login for email: %s", req.Email)

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Login failed: user not found for email %s", req.Email)
			return nil, ErrUserNotFound
		}
		log.Printf("Database error during login: %v", err)
		return nil, ErrDatabaseError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("Login failed: invalid password for user %s", req.Email)
		return nil, ErrInvalidPassword
	}

	expiresAt := time.Now().Add(time.Hour * 24)
	log.Printf("Generating JWT token for user %s, expires at: %s", user.Email, expiresAt)

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Printf("Login failed: JWT_SECRET environment variable is not set")
		return nil, ErrJWTSecretMissing
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"name":     user.Name,
		"exp":      expiresAt.Unix(),
	}

	// Add role if user has roles
	if len(user.Roles) > 0 {
		claims["role"] = user.Roles[0].Name
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Error signing JWT token: %v", err)
		return nil, errors.New("error generating authentication token")
	}

	// Create response
	response := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Name:      user.Name,
	}

	// Add role to response if user has roles
	if len(user.Roles) > 0 {
		response.Role = user.Roles[0].Name
	}

	log.Printf("Login successful for user: %s", user.Email)
	return &response, nil
}

func (s *AuthService) GetUserProfile(userID uint) (*UserProfileResponse, error) {
	log.Printf("Fetching profile for user ID: %d", userID)

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Profile fetch failed: user not found with ID %d", userID)
			return nil, ErrUserNotFound
		}
		log.Printf("Database error while fetching profile: %v", err)
		return nil, ErrDatabaseError
	}

	log.Printf("User found: %+v", user)

	// Convert roles to string slice
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	response := &UserProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Name:      user.Name,
		Birthday:  time.Time(user.Birthday),
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	log.Printf("Created profile response: %+v", response)

	// Create audit log
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    "view_profile",
		TableName: "users",
		RecordID:  user.ID,
		NewValues: models.JSON{
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"birthday": time.Time(user.Birthday).Format("2006-01-02"),
			"roles":    roles,
		},
	}

	if err := s.auditRepo.CreateLog(auditLog); err != nil {
		log.Printf("Failed to create audit log for profile view: %v", err)
		// Don't return error, just log it
	}

	log.Printf("Profile fetched successfully for user: %s", user.Email)
	return response, nil
}
