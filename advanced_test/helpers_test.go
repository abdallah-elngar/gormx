package advanced_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/abdallah-elngar/gormx"
	"github.com/abdallah-elngar/gormx/tests/fixtures"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User = fixtures.User
type Order = fixtures.Order
type Product = fixtures.Product

// setupTestDB creates a new DB + migrations + cleanup.
//
// Each call creates a UNIQUE in-memory database to prevent
// data leakage between tests.
func setupTestDB(t *testing.T) {
	t.Helper()

	// ✅ اسم فريد لكل اختبار — يمنع تسرب البيانات
	dsn := fmt.Sprintf("file:test_adv_%d?mode=memory&cache=shared", time.Now().UnixNano())

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&User{}, &Order{}, &Product{}))

	gormx.SetDB(db)
	gormx.ClearRegistry()

	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil && sqlDB != nil {
			_ = sqlDB.Close()
		}
		gormx.ResetDB()
		gormx.ClearRegistry()
	})
}

// seedUsers inserts 5 users.
func seedUsers(t *testing.T) []User {
	t.Helper()

	users := []User{
		{Name: "Ali", Email: "ali@test.com", Age: 30, Active: true},
		{Name: "Sara", Email: "sara@test.com", Age: 25, Active: true},
		{Name: "Omar", Email: "omar@test.com", Age: 17, Active: false},
		{Name: "Layla", Email: "layla@test.com", Age: 35, Active: true},
		{Name: "Hassan", Email: "hassan@test.com", Age: 40, Active: true},
	}

	for i := range users {
		require.NoError(t, gormx.New[User]().Create(&users[i]))
	}

	return users
}

// seedOrders inserts 5 orders.
func seedOrders(t *testing.T) []Order {
	t.Helper()

	orders := []Order{
		{UserID: 1, Total: 100, Status: "paid"},
		{UserID: 1, Total: 200, Status: "paid"},
		{UserID: 2, Total: 150, Status: "paid"},
		{UserID: 2, Total: 50, Status: "pending"},
		{UserID: 3, Total: 300, Status: "paid"},
	}

	for i := range orders {
		require.NoError(t, gormx.New[Order]().Create(&orders[i]))
	}

	return orders
}
