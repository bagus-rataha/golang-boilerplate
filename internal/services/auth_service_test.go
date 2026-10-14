package services

import (
	"errors"
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/dto"
	"fiber-api-boilerplate/internal/models"
	"fiber-api-boilerplate/internal/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func newTestConfig() *config.Config {
	return &config.Config{
		JWTAccessSecret:  "test-access-secret",
		JWTRefreshSecret: "test-refresh-secret",
		JWTAccessExpire:  time.Hour,
		JWTRefreshExpire: 24 * time.Hour,
	}
}

func newAuthService(userRepo *MockUserRepository, rtRepo *MockRefreshTokenRepository) *AuthService {
	return NewAuthService(userRepo, rtRepo, newTestConfig())
}

func TestRegister_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	userRepo.On("FindByEmail", "test@test.com").Return(nil, gorm.ErrRecordNotFound)
	userRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	rtRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.Register(dto.RegisterInput{
		Email:    "test@test.com",
		Password: "123456",
		Name:     "Test",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "test@test.com", result.User.Email)
}

func TestRegister_EmailExists(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	existingUser := &models.User{Email: "test@test.com"}
	userRepo.On("FindByEmail", "test@test.com").Return(existingUser, nil)

	svc := newAuthService(userRepo, rtRepo)
	_, err := svc.Register(dto.RegisterInput{
		Email:    "test@test.com",
		Password: "123456",
		Name:     "Test",
	})

	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
}

func TestLogin_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	hashedPassword, _ := utils.HashPassword("123456")
	user := &models.User{
		BaseModel: models.BaseModel{ID: uuid.Must(uuid.NewV7())},
		Email:     "test@test.com",
		Password:  hashedPassword,
		Name:      "Test",
		Role:      "user",
	}

	userRepo.On("FindByEmail", "test@test.com").Return(user, nil)
	rtRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.Login(dto.LoginInput{
		Email:    "test@test.com",
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	userRepo.On("FindByEmail", "test@test.com").Return(nil, gorm.ErrRecordNotFound)

	svc := newAuthService(userRepo, rtRepo)
	_, err := svc.Login(dto.LoginInput{
		Email:    "test@test.com",
		Password: "123456",
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	hashedPassword, _ := utils.HashPassword("123456")
	user := &models.User{
		Email:    "test@test.com",
		Password: hashedPassword,
	}

	userRepo.On("FindByEmail", "test@test.com").Return(user, nil)

	svc := newAuthService(userRepo, rtRepo)
	_, err := svc.Login(dto.LoginInput{
		Email:    "test@test.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogout_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	rtRepo.On("DeleteByToken", "some-token").Return(int64(1), nil)

	svc := newAuthService(userRepo, rtRepo)
	err := svc.Logout("some-token")

	assert.NoError(t, err)
}

func TestLogoutAll_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	userID := uuid.Must(uuid.NewV7())

	rtRepo.On("DeleteAllByUserID", userID).Return(nil)

	svc := newAuthService(userRepo, rtRepo)
	err := svc.LogoutAll(userID)

	assert.NoError(t, err)
}

func newRefreshTokenFor(t *testing.T, userID uuid.UUID) string {
	t.Helper()
	token, err := utils.GenerateToken(userID, "test@test.com", "user", "test-refresh-secret", time.Hour)
	assert.NoError(t, err)
	return token
}

func newTestUser() *models.User {
	return &models.User{
		BaseModel: models.BaseModel{ID: uuid.Must(uuid.NewV7())},
		Email:     "test@test.com",
		Name:      "Test",
		Role:      "user",
	}
}

func TestRefreshToken_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	user := newTestUser()
	oldToken := newRefreshTokenFor(t, user.ID)

	rtRepo.On("DeleteByToken", oldToken).Return(int64(1), nil)
	userRepo.On("FindByID", user.ID).Return(user, nil)
	rtRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.RefreshToken(oldToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEqual(t, oldToken, result.RefreshToken)
	rtRepo.AssertExpectations(t)
}

func TestRefreshToken_AlreadyConsumed(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	user := newTestUser()
	oldToken := newRefreshTokenFor(t, user.ID)

	rtRepo.On("DeleteByToken", oldToken).Return(int64(0), nil)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.RefreshToken(oldToken)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
	userRepo.AssertNotCalled(t, "FindByID", mock.Anything)
	rtRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRefreshToken_MalformedToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.RefreshToken("not-a-jwt")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
	rtRepo.AssertNotCalled(t, "DeleteByToken", mock.Anything)
}

func TestRefreshToken_AccessTokenRejected(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	accessToken, err := utils.GenerateToken(uuid.Must(uuid.NewV7()), "test@test.com", "user", "test-access-secret", time.Hour)
	assert.NoError(t, err)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.RefreshToken(accessToken)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
	rtRepo.AssertNotCalled(t, "DeleteByToken", mock.Anything)
}

func TestRefreshToken_UserGone(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	user := newTestUser()
	oldToken := newRefreshTokenFor(t, user.ID)

	rtRepo.On("DeleteByToken", oldToken).Return(int64(1), nil)
	userRepo.On("FindByID", user.ID).Return(nil, gorm.ErrRecordNotFound)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.RefreshToken(oldToken)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
	rtRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRefreshToken_StoreFailure(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	user := newTestUser()
	oldToken := newRefreshTokenFor(t, user.ID)
	dbErr := errors.New("connection refused")

	rtRepo.On("DeleteByToken", oldToken).Return(int64(0), dbErr)

	svc := newAuthService(userRepo, rtRepo)
	result, err := svc.RefreshToken(oldToken)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NotErrorIs(t, err, ErrInvalidRefreshToken)
}

func TestLogout_AlreadyGone(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)

	rtRepo.On("DeleteByToken", "gone-token").Return(int64(0), nil)

	svc := newAuthService(userRepo, rtRepo)
	err := svc.Logout("gone-token")

	assert.NoError(t, err)
}
