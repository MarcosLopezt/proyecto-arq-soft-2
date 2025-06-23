package user_test

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	models "users/models"
// 	"users/services/cache"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // Mock del repositorio de Memcached
// type MockCacheRepo struct {
// 	mock.Mock
// }

// func (m *MockCacheRepo) Create(ctx context.Context, user models.User) (int64, error) {
// 	args := m.Called(ctx, user)
// 	return args.Get(0).(int64), args.Error(1)
// }

// func (m *MockCacheRepo) GetUserByID(ctx context.Context, id string) (models.User, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(models.User), args.Error(1)
// }

// func (m *MockCacheRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
// 	args := m.Called(ctx, email)
// 	return args.Get(0).(models.User), args.Error(1)
// }

// func TestCacheService(t *testing.T) {
// 	mockCache := new(MockCacheRepo)
// 	cacheConfig := cache.MemcachedConfig{
// 		Host:     "localhost",
// 		Port:     "11211",
// 		Duration: 5 * time.Minute,
// 	}
// 	cacheService, _ := cache.NewCacheWithRepo(cacheConfig, mockCache) // Usa el mock en lugar del repo real

// 	t.Run("Create - Success", func(t *testing.T) {
// 		testUser := models.User{
// 			Email: "pablo.lopez@gmail.com",
// 			PasswordHash: "passowrd123",
// 		}

// 		mockCache.On("Create", mock.Anything, testUser).Return(int64(1), nil).Once()

// 		_, err := cacheService.Create(context.Background(), testUser)
// 		assert.NoError(t, err)

// 		mockCache.AssertExpectations(t)
// 	})

// 	t.Run("GetUserByID - Success", func(t *testing.T) {
// 		testUser := models.User{
// 			ID:    1,
// 			Email: "test@example.com",
// 		}

// 		mockCache.On("GetUserByID", mock.Anything, "1").Return(testUser, nil).Once()

// 		retrievedUser, err := cacheService.GetUserByID(context.Background(), "1")
// 		assert.NoError(t, err)
// 		assert.Equal(t, "test@example.com", retrievedUser.Email)

// 		mockCache.AssertExpectations(t)
// 	})

// 	t.Run("GetUserByID - Not Found", func(t *testing.T) {
// 		mockCache.On("GetUserByID", mock.Anything, "2").Return(models.User{}, errors.New("user not found")).Once()

// 		retrievedUser, err := cacheService.GetUserByID(context.Background(), "2")
// 		assert.Error(t, err)
// 		assert.Empty(t, retrievedUser)

// 		mockCache.AssertExpectations(t)
// 	})

// 	t.Run("GetUserByEmail - Success", func(t *testing.T) {
// 		testUser := models.User{
// 			ID:    1,
// 			Email: "test@example.com",
// 		}

// 		mockCache.On("GetUserByEmail", mock.Anything, "test@example.com").Return(testUser, nil).Once()

// 		retrievedUser, err := cacheService.GetUserByEmail(context.Background(), "test@example.com")
// 		assert.NoError(t, err)
// 		assert.Equal(t, 1, retrievedUser.ID)

// 		mockCache.AssertExpectations(t)
// 	})

// 	t.Run("GetUserByEmail - Not Found", func(t *testing.T) {
// 		mockCache.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(models.User{}, errors.New("user not found")).Once()

// 		retrievedUser, err := cacheService.GetUserByEmail(context.Background(), "unknown@example.com")
// 		assert.Error(t, err)
// 		assert.Empty(t, retrievedUser)

// 		mockCache.AssertExpectations(t)
// 	})
// }
