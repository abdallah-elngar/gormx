# Registry — سجل الموديلات

> تسجيل QuerySets لتجنّب import cycles.

---

## 🎯 لماذا Registry؟

### المشكلة: Import Cycle

```
pkg/user   →  import  →  pkg/order
pkg/order  →  import  →  pkg/user  ❌ cycle!
```

### الحل: Registry

```go
// في pkg/user
var Users = gormx.Register[User]("user")

// في pkg/order
q, _ := gormx.Lookup[User]("user")  // ← لا import!
```

---

## 📖 API

### `Register[T](name)`

```go
var Users = gormx.Register[User]("user")
```

**يرفع panic** عند التكرار.

### `TryRegister[T](name)`

```go
q, err := gormx.TryRegister[User]("user")
if err != nil {
    log.Fatal(err)
}
```

### `Lookup[T](name)`

```go
q, ok := gormx.Lookup[User]("user")
if !ok {
    log.Fatal("user not registered")
}
```

### `MustLookup[T](name)`

```go
q := gormx.MustLookup[User]("user")  // panic if not found
```

### Helpers

| Function | الوصف |
|----------|-------|
| `Has(name)` | هل مسجّل؟ |
| `Unregister(name)` | إزالة |
| `RegisteredNames()` | كل الأسماء |
| `RegisteredCount()` | العدد |
| `ClearRegistry()` | مسح كل شيء (للاختبارات) |

---

## 💡 أمثلة

### مثال 1: Basic

```go
// apps/user/objects.go
package user

import "github.com/abdallah-elngar/gormx"

var Users = gormx.Register[User]("user")

// apps/order/service.go
package order

import "github.com/abdallah-elngar/gormx"

func FindUserOrders(userID uint) ([]Order, error) {
    // احصل على QuerySet من Registry
    users := gormx.MustLookup[User]("user")

    // استخدمه
    user, err := users.Get(userID)
    if err != nil {
        return nil, err
    }

    return gormx.New[Order]().Filter("user_id", user.ID).All()
}
```

### مثال 2: Cross-Package

```
pkg/
├── user/
│   └── objects.go    (يُسجّل)
├── order/
│   └── service.go    (يستخدم)
└── report/
    └── service.go    (يستخدم)
```

```go
// pkg/user/objects.go
var Users = gormx.Register[User]("user")

// pkg/order/service.go
func CreateOrder(userID uint, total float64) (*Order, error) {
    // احصل على Users من registry
    users := gormx.MustLookup[User]("user")
    user, err := users.Get(userID)
    if err != nil {
        return nil, err
    }

    order := &Order{UserID: user.ID, Total: total}
    return order, gormx.New[Order]().Create(order)
}

// pkg/report/service.go
func TopUsers(limit int) ([]User, error) {
    users := gormx.MustLookup[User]("user")
    return users.OrderBy("-created_at").Limit(limit).All()
}
```

### مثال 3: Init Pattern

```go
// apps/user/manifest.go
package user

import "github.com/abdallah-elngar/gormx"

var Objects *gormx.QuerySet[User]

func init() {
    Objects = gormx.Register[User]("user")
}

// استخدام لاحقًا
func Find(email string) (*User, error) {
    return Objects.Find("email", email)
}
```

### مثال 4: Testing

```go
func TestSomething(t *testing.T) {
    setupTestDB(t)
    gormx.ClearRegistry()  // ← نظّف قبل

    // سجّل
    Users := gormx.Register[User]("user")

    // استخدم
    Users.Create(&User{Name: "Ali"})

    // تأكد
    assert.Equal(t, 1, gormx.RegisteredCount())
}
```

---

## ⚠️ تحذيرات

### 1. Duplicate Registration

```go
gormx.Register[User]("user")  // ok
gormx.Register[User]("user")  // panic!
```

**الحل**: استخدم `TryRegister` أو `Has()`:

```go
if !gormx.Has("user") {
    gormx.Register[User]("user")
}
```

### 2. Type Mismatch

```go
gormx.Register[User]("user")

// ❌ خطأ نوع
q, ok := gormx.Lookup[Order]("user")
// q = nil, ok = false

// ✅ صحيح
q, ok := gormx.Lookup[User]("user")
```

### 3. Thread Safety

Registry **thread-safe**:

```go
// آمن من goroutines متعددة
go gormx.Register[User]("user")
go gormx.Register[Order]("order")
```

### 4. Global State

Registry **global**. للاختبارات:

```go
func TestMain(m *testing.M) {
    // Setup
    code := m.Run()

    // Cleanup
    gormx.ClearRegistry()
    os.Exit(code)
}
```

---

## 🎯 Best Practices

### ✅ Do

- اسم الموديل بصيغة `singular` (`"user"` بدل `"users"`)
- استخدم `init()` للتسجيل التلقائي
- `ClearRegistry` في الاختبارات
- `Has` قبل `Register` إن أمكن

### ❌ Don't

- لا تسجّل نفس الاسم مرتين
- لا تخلط الأنواع
- لا تعتمد على ترتيب التسجيل