package gormx_test

import (
	"strings"
	"testing"

	"github.com/abdallah-elngar/gormx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserWithHooks نموذج مع hooks
type UserWithHooks struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:255"`
	Email string `gorm:"size:255"`
}

func (u *UserWithHooks) BeforeCreate(tx *gorm.DB) error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

func TestHooks_BeforeCreate(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, db.AutoMigrate(&UserWithHooks{}))
	gormx.SetDB(db)

	user := &UserWithHooks{
		Name:  "  Ali  ",
		Email: "  ALI@TEST.COM  ",
	}

	require.NoError(t, gormx.New[UserWithHooks]().Create(user))

	got, _ := gormx.New[UserWithHooks]().Get(user.ID)
	assert.Equal(t, "Ali", got.Name)
	assert.Equal(t, "ali@test.com", got.Email)
}