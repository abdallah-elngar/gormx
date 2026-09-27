// Package gormx provides a Django-inspired, type-safe ORM layer on top of GORM.
//
// # Philosophy
//
// gormx keeps GORM's power while offering a cleaner, more intuitive API:
//
//	// Instead of GORM's verbose queries:
//	var users []User
//	db.Model(&User{}).Where("age > ?", 18).Order("name").Find(&users)
//
//	// Use gormx's fluent, type-safe API:
//	users, _ := gormx.New[User]().Filter("age__gt", 18).OrderBy("name").All()
//
// # Key Features
//
//   - Type-safe generic QuerySet[T]
//   - Django-style lookups (`__gt`, `__in`, `__contains`, ...)
//   - Multi-DB support via Instance
//   - Immutable, thread-safe queries
//   - Context-first design
//   - Built-in pagination
//   - Model registry to avoid import cycles
//   - Advanced features: CTEs, Window Functions, Locks, Retries
//
// # Quick Start
//
//	// 1. Set up the global DB
//	db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
//	gormx.SetDB(db)
//
//	// 2. Migrate
//	gormx.MustMigrate[User]()
//
//	// 3. Query
//	users, err := gormx.New[User]().
//	    Filter("active", true).
//	    Filter("age__gte", 18).
//	    OrderBy("-created_at").
//	    Limit(10).
//	    All()
//
// # Multi-DB
//
//	app := gormx.NewInstance(db)
//	q := app.Query[User]()
//	users, _ := q.Filter("active", true).All()
//
// # Registry
//
//	var Objects = gormx.Register[User]("user")
package gormx

import (
	"context"
	"errors"

	"github.com/abdallah-elngar/gormx/internal"
	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
// Version & Metadata
// ═══════════════════════════════════════════════

// Version هو إصدار gormx الحالي.
const Version = "0.1.0"

// Author معلومات المؤلف.
const (
	Author  = "Sanad Team"
	License = "MIT"
	URL     = "https://github.com/abdallah-elngar/gormx"
)

// ═══════════════════════════════════════════════
// Global DB — Backward Compatible
// ═══════════════════════════════════════════════

var globalDB *gorm.DB

// SetDB يربط GORM عالميًا. يجب أن يُنادى قبل أي استعلام.
//
//	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	gormx.SetDB(db)
//
// يpanic إذا كان db == nil.
func SetDB(d *gorm.DB) {
	if d == nil {
		panic(ErrNilDB)
	}
	globalDB = d
}

// DB يرجّع *gorm.DB الأصلي.
//
// يpanic إذا لم يُنادى SetDB.
func DB() *gorm.DB {
	if globalDB == nil {
		panic(ErrNotInitialized)
	}
	return globalDB
}

// IsReady يفحص إذا كان الاتصال جاهزًا.
func IsReady() bool {
	return globalDB != nil
}

// ResetDB يعيد تعيين الاتصال العام (مفيد في الاختبارات).
//
// ⚠️ للاختبارات فقط. لا تستخدمه في الإنتاج.
func ResetDB() {
	globalDB = nil
}

// ═══════════════════════════════════════════════
// Migrate — الترحيل
// ═══════════════════════════════════════════════

// Migrate يهاجر موديل T.
//
//	err := gormx.Migrate[User]()
func Migrate[T any]() error {
	var zero T
	return DB().AutoMigrate(&zero)
}

// MustMigrate مثل Migrate لكن يpanic عند الخطأ.
//
//	gormx.MustMigrate[User]()
func MustMigrate[T any]() {
	if err := Migrate[T](); err != nil {
		panic(err)
	}
}

// MigrateAll يهاجر عدة موديلات.
//
//	gormx.MigrateAll(&User{}, &Order{}, &Product{})
func MigrateAll(models ...any) error {
	if len(models) == 0 {
		return nil
	}
	return DB().AutoMigrate(models...)
}

// MustMigrateAll مثل MigrateAll لكن يpanic.
func MustMigrateAll(models ...any) {
	if err := MigrateAll(models...); err != nil {
		panic(err)
	}
}

// DropTable يحذف جدول موديل T.
//
// ⚠️ يحذف كل البيانات!
//
//	gormx.DropTable[User]()
func DropTable[T any]() error {
	var zero T
	return DB().Migrator().DropTable(&zero)
}

// HasTable يفحص وجود جدول موديل T.
func HasTable[T any]() bool {
	var zero T
	return DB().Migrator().HasTable(&zero)
}

// ═══════════════════════════════════════════════
// Helpers — أدوات مساعدة
// ═══════════════════════════════════════════════

// TableNameOf يرجّع اسم الجدول لموديل T.
//
//	fmt.Println(gormx.TableNameOf[User]())  // "users"
func TableNameOf[T any]() string {
	return internal.TableName[T]()
}

// ModelNameOf يرجّع اسم الموديل T.
//
//	fmt.Println(gormx.ModelNameOf[User]())  // "User"
func ModelNameOf[T any]() string {
	return internal.ModelName[T]()
}

// FieldNamesOf يرجّع أسماء الأعمدة لموديل T.
//
//	fmt.Println(gormx.FieldNamesOf[User]())  // ["id", "name", "email"]
func FieldNamesOf[T any]() []string {
	return internal.FieldNames[T]()
}

// ColumnNameOf يرجّع اسم العمود لموديل T.
//
//	col, err := gormx.ColumnNameOf[User]("Email")
//	// col = "email"
func ColumnNameOf[T any](goField string) (string, error) {
	info := internal.GetModelInfo[T]()
	for _, f := range info.Fields {
		if f.GoName == goField {
			return f.Column, nil
		}
	}
	return "", ErrInvalidField
}

// ═══════════════════════════════════════════════
// Type Aliases — للاستخدام في advanced
// ═══════════════════════════════════════════════

// Context هو alias لـ context.Context.
type Context = context.Context

// Tx هو alias لـ *gorm.DB داخل transaction.
type Tx = gorm.DB

// ═══════════════════════════════════════════════
// Transaction — معاملة بسيطة
// ═══════════════════════════════════════════════

// Transaction ينفّذ عملية داخل transaction.
//
//	err := gormx.Transaction(ctx, func(tx *gorm.DB) error {
//	    if err := tx.Create(&user).Error; err != nil {
//	        return err
//	    }
//	    return tx.Create(&order).Error
//	})
func Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if fn == nil {
		return errors.New("gormx: nil transaction callback")
	}
	return DB().WithContext(ctx).Transaction(fn)
}

// ═══════════════════════════════════════════════
// Connection Pool
// ═══════════════════════════════════════════════

// Configure يضبط إعدادات الاتصال العام.
func Configure(cfg Config) error {
	return GlobalInstance().Configure(cfg)
}

// Ping يفحص الاتصال العام.
func Ping() error {
	return GlobalInstance().Ping()
}

// Close يغلق الاتصال العام.
func Close() error {
	return GlobalInstance().Close()
}