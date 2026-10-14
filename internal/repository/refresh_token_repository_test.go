//go:build integration

package repository

import (
	"fiber-api-boilerplate/internal/models"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRefreshTokenRepository(db)

	user := &models.User{Email: "test@test.com", Password: "h", Name: "Test", Role: "user"}
	db.Create(user)

	rt := &models.RefreshToken{
		Token:     "some-token",
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	err := repo.Create(rt)

	assert.NoError(t, err)
	assert.NotEmpty(t, rt.ID)
}

func TestRefreshTokenRepository_FindByToken_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRefreshTokenRepository(db)

	user := &models.User{Email: "test@test.com", Password: "h", Name: "Test", Role: "user"}
	db.Create(user)
	db.Create(&models.RefreshToken{Token: "my-token", UserID: user.ID, ExpiredAt: time.Now().Add(time.Hour)})

	found, err := repo.FindByToken("my-token")

	assert.NoError(t, err)
	assert.Equal(t, "my-token", found.Token)
}

func TestRefreshTokenRepository_DeleteByToken(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRefreshTokenRepository(db)

	user := &models.User{Email: "test@test.com", Password: "h", Name: "Test", Role: "user"}
	db.Create(user)
	db.Create(&models.RefreshToken{Token: "to-delete", UserID: user.ID, ExpiredAt: time.Now().Add(time.Hour)})

	deleted, err := repo.DeleteByToken("to-delete")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), deleted)

	_, err = repo.FindByToken("to-delete")
	assert.Error(t, err)
}

func TestRefreshTokenRepository_DeleteAllByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRefreshTokenRepository(db)

	user := &models.User{Email: "test@test.com", Password: "h", Name: "Test", Role: "user"}
	db.Create(user)

	db.Create(&models.RefreshToken{Token: "token-1", UserID: user.ID, ExpiredAt: time.Now().Add(time.Hour)})
	db.Create(&models.RefreshToken{Token: "token-2", UserID: user.ID, ExpiredAt: time.Now().Add(time.Hour)})

	err := repo.DeleteAllByUserID(user.ID)
	assert.NoError(t, err)

	_, err = repo.FindByToken("token-1")
	assert.Error(t, err)
}

func TestRefreshTokenRepository_DeleteByToken_Missing(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRefreshTokenRepository(db)

	deleted, err := repo.DeleteByToken("never-existed")

	assert.NoError(t, err)
	assert.Equal(t, int64(0), deleted)
}

func TestRefreshTokenRepository_DeleteByToken_ConcurrentCallsOnlyOneWins(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRefreshTokenRepository(db)

	user := &models.User{Email: "test@test.com", Password: "h", Name: "Test", Role: "user"}
	db.Create(user)
	db.Create(&models.RefreshToken{Token: "contested", UserID: user.ID, ExpiredAt: time.Now().Add(time.Hour)})

	const callers = 10
	var wg sync.WaitGroup
	var removed atomic.Int64
	start := make(chan struct{})
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			deleted, err := repo.DeleteByToken("contested")
			assert.NoError(t, err)
			removed.Add(deleted)
		}()
	}
	close(start)
	wg.Wait()

	assert.Equal(t, int64(1), removed.Load())
}
