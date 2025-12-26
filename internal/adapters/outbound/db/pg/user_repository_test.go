package pg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thealiakbari/task-pool-system/internal/domain/task/entity"
	"github.com/thealiakbari/task-pool-system/pkg/common/config"
	"github.com/thealiakbari/task-pool-system/pkg/common/db"
)

func setupTestDB(t *testing.T) db.DBWrapper {
	conf := config.LoadConfig("../../../../../config/config.yml")
	gormDB, err := db.NewPostgresConn(context.Background(), conf.DB.Postgres)
	assert.NoError(t, err)
	dbw := db.NewDBWrapper(gormDB)

	return dbw
}

func TestUserRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	testDB := setupTestDB(t)
	repo := NewUserRepository(testDB)

	// Create
	item := entity.User{
		Username:    "Test",
		PhoneNumber: "989121111111",
	}
	created, err := repo.Create(ctx, item)
	assert.NoError(t, err)
	assert.Equal(t, item.PhoneNumber, created.Username)

	// FindByIdOrEmpty
	found, err := repo.FindByIdOrEmpty(ctx, created.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, created.Id, found.Id)

	// Update
	created.Username = "UpdatedTest"
	err = repo.Update(ctx, created)
	assert.NoError(t, err)

	updated, err := repo.FindByIdOrEmpty(ctx, created.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedTask", updated.Username)

	// FindByIds
	list, err := repo.FindByIds(ctx, []string{created.Id.String()})
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// FilterFind
	results, err := repo.FilterFind(ctx, []any{"username LIKE ?", "%Task%"}, "created_at desc", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	// FilterCount
	count, err := repo.FilterCount(ctx, []any{"username LIKE ?", "%Task%"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Delete
	err = repo.Delete(ctx, created.Id.String())
	assert.NoError(t, err)

	_, err = repo.FindByIdOrEmpty(ctx, created.Id.String())
	assert.NoError(t, err) // should return empty entity, not fail

	// Purge
	// Re-create and then purge
	item2 := entity.User{
		Username:    "TempUser",
		PhoneNumber: "989121111111",
	}
	created2, _ := repo.Create(ctx, item2)
	err = repo.Purge(ctx, created2.Id.String())
	assert.NoError(t, err)
}
