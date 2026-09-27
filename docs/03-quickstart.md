# Quickstart — 10 دقائق

> كل ما تحتاجه لتصبح منتجًا مع gormx.

---

## 🎯 الفكرة الأساسية

**بدلًا من**:
```go
var users []User
db.Model(&User{}).
    Where("age > ?", 18).
    Where("status = ?", "active").
    Order("name").
    Limit(10).
    Find(&users)
```

**استخدم**:
```go
users, _ := gormx.New[User]().
    Filter("age__gt", 18).
    Filter("status", "active").
    OrderBy("name").
    Limit(10).
    All()
```

---

## 📖 1. الموديلات

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:255;not null"`
    Email     string    `gorm:"uniqueIndex;size:255"`
    Age       int
    Active    bool      `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**ملاحظات**:
- `gorm.DeletedAt` يفعّل soft delete
- `uniqueIndex` لمنع التكرار
- `size:255` لتحديد حجم VARCHAR

---

## 📖 2. الإعداد

```go
db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
gormx.SetDB(db)
gormx.MustMigrate[User]()
```

---

## 📖 3. CRUD

### Create

```go
user := &User{Name: "Ali", Email: "ali@test.com", Age: 30}
err := gormx.New[User]().Create(user)
// user.ID = 1
```

### Create Many

```go
users := []User{
    {Name: "Ali", Email: "ali@test.com"},
    {Name: "Sara", Email: "sara@test.com"},
}
err := gormx.New[User]().CreateMany(users)
```

### Read — واحد

```go
user, err := gormx.New[User]().Get(1)
// user = *User
```

### Read — قائمة

```go
users, err := gormx.New[User]().All()
// users = []User
```

### Read — بفلتر

```go
users, err := gormx.New[User]().
    Filter("active", true).
    Filter("age__gte", 18).
    All()
```

### Update

```go
err := gormx.New[User]().Update(1, "age", 31)
```

### Delete

```go
err := gormx.New[User]().Delete(1)
// Soft delete
```

---

## 📖 4. Lookups

```go
// مقارنة
.Filter("age__gt", 18)      // age > 18
.Filter("age__gte", 18)     // age >= 18
.Filter("age__lt", 65)      // age < 65
.Filter("age__lte", 65)     // age <= 65
.Filter("age__ne", 0)       // age != 0

// نصوص
.Filter("name__contains", "Ali")     // name LIKE '%Ali%'
.Filter("name__startswith", "A")     // name LIKE 'A%'
.Filter("name__endswith", "com")     // name LIKE '%com'

// قوائم
.Filter("status__in", []string{"active", "pending"})

// NULL
.Filter("deleted_at__isnull", true)

// مدى
.Filter("age__between", []int{18, 65})
```

---

## 📖 5. Q Builder

```go
// OR
q := gormx.QOr(
    gormx.Eq("status", "active"),
    gormx.Eq("status", "pending"),
)
users, _ := gormx.New[User]().Q(q).All()

// AND
q := gormx.QAnd(
    gormx.Eq("active", true),
    gormx.Gt("age", 18),
)

// معقد
q := gormx.Qb().And(
    gormx.QOr(
        gormx.Eq("status", "active"),
        gormx.Eq("status", "pending"),
    ),
    gormx.Gt("age", 18),
)

// NOT
q := gormx.Qb().And(
    gormx.Not(gormx.Eq("status", "deleted")),
    gormx.Eq("active", true),
)
```

---

## 📖 6. Pagination

```go
page, _ := gormx.New[User]().
    Filter("active", true).
    OrderBy("-created_at").
    Paginate(1, 20)

fmt.Println(page.Items)      // []User
fmt.Println(page.Total)      // 150
fmt.Println(page.Page)       // 1
fmt.Println(page.PerPage)    // 20
fmt.Println(page.TotalPages) // 8
fmt.Println(page.HasNext)    // true
fmt.Println(page.HasPrev)    // false
```

---

## 📖 7. Hooks

```go
type User struct {
    ID    uint
    Name  string
    Email string
}

func (u *User) BeforeCreate() error {
    u.Name = strings.TrimSpace(u.Name)
    u.Email = strings.ToLower(u.Email)
    return nil
}
```

---

## 📖 8. Transactions

```go
err := gormx.Transaction(context.Background(), func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    return tx.Create(&order).Error
})
```

---

## 🎯 الخلاصة

الآن أنت تعرف:
- ✅ CRUD
- ✅ Lookups
- ✅ Q Builder
- ✅ Pagination
- ✅ Hooks
- ✅ Transactions

**انتقل إلى**:
- [Advanced](../advanced/README.md) — ميزات متقدمة
- [Examples](../examples/README.md) — أمثلة كاملة