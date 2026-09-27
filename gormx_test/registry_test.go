package gormx_test

import (
	"testing"

	"github.com/abdallah-elngar/gormx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ═══════════════════════════════════════════════
// Basic Registry
// ═══════════════════════════════════════════════

func TestRegistry_Register(t *testing.T) {
	setupTestDB(t)

	q := gormx.Register[User]("user")
	require.NotNil(t, q)

	got, ok := gormx.Lookup[User]("user")
	assert.True(t, ok)
	assert.Same(t, q, got)
}

func TestRegistry_RegisterDuplicate(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")

	assert.Panics(t, func() {
		gormx.Register[User]("user")
	})
}

func TestRegistry_TryRegisterDuplicate(t *testing.T) {
	setupTestDB(t)

	_, err := gormx.TryRegister[User]("user")
	require.NoError(t, err)

	_, err = gormx.TryRegister[User]("user")
	assert.Error(t, err)
	assert.True(t, gormx.IsAlreadyRegistered(err))
}

func TestRegistry_LookupNotFound(t *testing.T) {
	setupTestDB(t)

	_, ok := gormx.Lookup[User]("not-registered")
	assert.False(t, ok)
}

func TestRegistry_MustLookupPanics(t *testing.T) {
	setupTestDB(t)

	assert.Panics(t, func() {
		gormx.MustLookup[User]("not-registered")
	})
}

func TestRegistry_Unregister(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Unregister("user")

	_, ok := gormx.Lookup[User]("user")
	assert.False(t, ok)
}

func TestRegistry_RegisteredNames(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Register[Order]("order")

	names := gormx.RegisteredNames()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "user")
	assert.Contains(t, names, "order")
}

func TestRegistry_RegisteredCount(t *testing.T) {
	setupTestDB(t)

	assert.Equal(t, 0, gormx.RegisteredCount())

	gormx.Register[User]("user")
	gormx.Register[Order]("order")

	assert.Equal(t, 2, gormx.RegisteredCount())
}

func TestRegistry_Usage(t *testing.T) {
	setupTestDB(t)

	Users := gormx.Register[User]("user")
	require.NoError(t, Users.Create(&User{Name: "Ali", Email: "ali@test.com"}))

	q := gormx.MustLookup[User]("user")
	users, _ := q.All()
	assert.Len(t, users, 1)
}

// ═══════════════════════════════════════════════
// Type-based Lookup
// ═══════════════════════════════════════════════

func TestRegistry_HasType(t *testing.T) {
	setupTestDB(t)

	assert.False(t, gormx.HasType[User]())

	gormx.Register[User]("user")
	assert.True(t, gormx.HasType[User]())
}

func TestRegistry_LookupByType(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	q, ok := gormx.LookupByType[User]()
	assert.True(t, ok)
	assert.NotNil(t, q)
}

func TestRegistry_MustLookupByType(t *testing.T) {
	setupTestDB(t)

	assert.Panics(t, func() {
		gormx.MustLookupByType[User]()
	})

	gormx.Register[User]("user")
	assert.NotPanics(t, func() {
		gormx.MustLookupByType[User]()
	})
}

func TestRegistry_NamesByType(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Register[User]("admin_user")
	gormx.Register[Order]("order")

	names := gormx.NamesByType[User]()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "user")
	assert.Contains(t, names, "admin_user")
	assert.NotContains(t, names, "order")
}

// ═══════════════════════════════════════════════
// All Entries & Summary
// ═══════════════════════════════════════════════

func TestRegistry_AllEntries(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Register[Order]("order")

	entries := gormx.AllEntries()
	require.Len(t, entries, 2)

	// مرتبة حسب الاسم
	assert.Equal(t, "order", entries[0].Name)
	assert.Equal(t, "Order", entries[0].TypeName)
	assert.Equal(t, "user", entries[1].Name)
	assert.Equal(t, "User", entries[1].TypeName)
}

func TestRegistry_Summary(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Register[User]("admin_user")
	gormx.Register[Order]("order")

	summary := gormx.Summary()
	assert.Equal(t, 3, summary.Total)
	assert.Len(t, summary.ByName, 3)
	assert.Equal(t, 2, summary.ByType["User"])
	assert.Equal(t, 1, summary.ByType["Order"])
}

// ═══════════════════════════════════════════════
// Bulk Register
// ═══════════════════════════════════════════════

func TestRegistry_RegisterBatch(t *testing.T) {
	setupTestDB(t)

	entries := map[string]any{
		"user":  gormx.New[User](),
		"order": gormx.New[Order](),
	}

	err := gormx.RegisterBatch(entries)
	require.NoError(t, err)

	assert.Equal(t, 2, gormx.RegisteredCount())
}

func TestRegistry_RegisterBatchDuplicate(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")

	entries := map[string]any{
		"user":  gormx.New[User](),
		"order": gormx.New[Order](),
	}

	err := gormx.RegisterBatch(entries)
	assert.Error(t, err)
	assert.True(t, gormx.IsAlreadyRegistered(err))

	// لم يُضف أي مدخل (atomic)
	assert.Equal(t, 1, gormx.RegisteredCount())
}

func TestRegistry_UnregisterBatch(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Register[Order]("order")
	gormx.Register[User]("admin")

	count := gormx.UnregisterBatch("user", "order", "missing")
	assert.Equal(t, 2, count)
	assert.Equal(t, 1, gormx.RegisteredCount())
}

// ═══════════════════════════════════════════════
// Snapshot
// ═══════════════════════════════════════════════

func TestRegistry_Snapshot(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	gormx.Register[Order]("order")

	snap := gormx.TakeSnapshot()

	// مسح الكل
	gormx.ClearRegistry()
	assert.Equal(t, 0, gormx.RegisteredCount())

	// استعادة
	snap.Restore()
	assert.Equal(t, 2, gormx.RegisteredCount())
}

func TestRegistry_SnapshotMerge(t *testing.T) {
	setupTestDB(t)

	gormx.Register[User]("user")
	snap := gormx.TakeSnapshot()

	// إضافة order
	gormx.Register[Order]("order")

	// merge لا يستبدل
	count := snap.Merge()
	assert.Equal(t, 0, count) // user موجود مسبقًا
	assert.Equal(t, 2, gormx.RegisteredCount())
}