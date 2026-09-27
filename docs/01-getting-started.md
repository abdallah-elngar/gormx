# البدء السريع مع gormx

> من الصفر إلى أول استعلام في 5 دقائق.

---

## 📋 المتطلبات

- **Go 1.22+** — [تحميل](https://go.dev/dl/)
- **قاعدة بيانات** — SQLite (افتراضي) أو PostgreSQL أو MySQL

---

## 🚀 خطوات البدء

### 1. إنشاء مشروع

```bash
mkdir myapp && cd myapp
go mod init myapp
```

### 2. تثبيت gormx

```bash
go get github.com/sanad/gormx
go get gorm.io/driver/sqlite
go get gorm.io/gorm
```

### 3. أول برنامج

**`main.go`**:

```go
package main

import (
    "fmt"
    "log"

    "github.com/sanad/gormx"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// 1. عرّف الموديل
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:255;not null"`
    Email string `gorm:"uniqueIndex;size:255"`
    Age   int
}

func main() {
    // 2. الاتصال بقاعدة البيانات
    db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 3. ربط gormx
    gormx.SetDB(db)

    // 4. ترحيل الموديلات
    if err := gormx.Migrate[User](); err != nil {
        log.Fatal(err)
    }

    // 5. إنشاء
    user := &User{Name: "علي", Email: "ali@test.com", Age: 30}
    if err := gormx.New[User]().Create(user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user ID: %d\n", user.ID)

    // 6. قراءة
    users, err := gormx.New[User]().
        Filter("age__gte", 18).
        OrderBy("name").
        All()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d users\n", len(users))

    // 7. تعديل
    if err := gormx.New[User]().Update(user.ID, "age", 31); err != nil {
        log.Fatal(err)
    }

    // 8. حذف
    if err := gormx.New[User]().Delete(user.ID); err != nil {
        log.Fatal(err)
    }

    fmt.Println("✅ Done!")
}
```

### 4. شغّل

```bash
go run main.go
```

**النتيجة**:
```
Created user ID: 1
Found 1 users
✅ Done!
```

---

## 🎯 ما تعلمته

- ✅ الاتصال بقاعدة البيانات
- ✅ ربط gormx
- ✅ ترحيل الموديلات
- ✅ CRUD الأساسي
- ✅ الفلاتر
- ✅ الترتيب

---

## 📖 الخطوة التالية

- [QuerySet الكامل](../core/queryset.md)
- [كل الـ Lookups](../core/lookups.md)
- [الميزات المتقدمة](../advanced/README.md)
- [أمثلة عملية](../examples/README.md)

---

## 🔍 مثال أعمق

**CRUD كامل**:

```go
// Create
user := &User{Name: "Ali", Email: "ali@test.com"}
err := gormx.New[User]().Create(user)

// Read — واحد
u, err := gormx.New[User]().Get(1)

// Read — قائمة
users, err := gormx.New[User]().
    Filter("active", true).
    OrderBy("-created_at").
    Limit(10).
    All()

// Update
err := gormx.New[User]().Update(1, "name", "New Name")

// Delete
err := gormx.New[User]().Delete(1)

// Count
count, err := gormx.New[User]().Count()

// Exists
exists, err := gormx.New[User]().Filter("email", "ali@test.com").Exists()
```