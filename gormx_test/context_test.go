package gormx_test

import (
	"context"
	"testing"

	"github.com/abdallah-elngar/gormx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestContext_WithDB(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, db.AutoMigrate(&User{}))

	app := gormx.NewInstance(db)
	ctx := gormx.WithDB(context.Background(), app)

	require.NoError(t, gormx.FromContext[User](ctx).Create(&User{Name: "Ali"}))

	users, err := gormx.FromContext[User](ctx).All()
	require.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestContext_FallbackToGlobal(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, db.AutoMigrate(&User{}))
	gormx.SetDB(db)

	ctx := context.Background()

	require.NoError(t, gormx.FromContext[User](ctx).Create(&User{Name: "Ali"}))

	count, _ := gormx.FromContext[User](ctx).Count()
	assert.Equal(t, int64(1), count)
}

func TestContext_DBFromContext(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	app := gormx.NewInstance(db)

	_, ok := gormx.DBFromContext(context.Background())
	assert.False(t, ok)

	ctx := gormx.WithDB(context.Background(), app)
	instance, ok := gormx.DBFromContext(ctx)
	assert.True(t, ok)
	assert.Same(t, app, instance)
}

func TestContext_WithGormDB(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	ctx := gormx.WithGormDB(context.Background(), db)

	instance, ok := gormx.DBFromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, instance)
}

func TestContext_WithGormDBNilPanics(t *testing.T) {
	assert.Panics(t, func() {
		gormx.WithGormDB(context.Background(), nil)
	})
}