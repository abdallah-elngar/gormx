# Q Builder — الشروط المعقدة

> بناء شروط AND/OR/NOT معقدة.

---

## 📖 الأساسيات

### Constructors

| Function | SQL |
|----------|-----|
| `Qb()` | `AND` افتراضي |
| `QOr(a, b)` | `a OR b` |
| `QAnd(a, b)` | `a AND b` |

### Conditions

| Function | SQL |
|----------|-----|
| `Eq(field, v)` | `field = v` |
| `Ne(field, v)` | `field != v` |
| `Gt(field, v)` | `field > v` |
| `Gte(field, v)` | `field >= v` |
| `Lt(field, v)` | `field < v` |
| `Lte(field, v)` | `field <= v` |
| `Contains(field, s)` | `field LIKE '%s%'` |
| `StartsWith(field, s)` | `field LIKE 's%'` |
| `EndsWith(field, s)` | `field LIKE '%s'` |
| `In(field, values)` | `field IN (...)` |
| `IsNull(field)` | `field IS NULL` |
| `NotNull(field)` | `field IS NOT NULL` |

### Negation

```go
Not(condition)
```

---

## 💡 أمثلة

### مثال 1: OR بسيط

```go
q := gormx.QOr(
    gormx.Eq("status", "active"),
    gormx.Eq("status", "pending"),
)

users, _ := gormx.New[User]().Q(q).All()
// SELECT * FROM users WHERE status = 'active' OR status = 'pending'
```

### مثال 2: AND

```go
q := gormx.QAnd(
    gormx.Eq("active", true),
    gormx.Gt("age", 18),
)

// active = true AND age > 18
```

### مثال 3: معقّد

```go
// (status = 'active' OR status = 'pending') AND age > 18
q := gormx.Qb().And(
    gormx.QOr(
        gormx.Eq("status", "active"),
        gormx.Eq("status", "pending"),
    ),
    gormx.Gt("age", 18),
)

users, _ := gormx.New[User]().Q(q).All()
```

### مثال 4: NOT

```go
q := gormx.QAnd(
    gormx.Not(gormx.Eq("status", "deleted")),
    gormx.Eq("active", true),
)

// NOT (status = 'deleted') AND active = true
// status != 'deleted' AND active = true
```

### مثال 5: NOT مع Q

```go
inner := gormx.QOr(
    gormx.Eq("a", 1),
    gormx.Eq("b", 2),
)

q := gormx.QAnd(
    gormx.Not(inner),
    gormx.Eq("c", 3),
)

// NOT (a = 1 OR b = 2) AND c = 3
```

### مثال 6: تداخل عميق

```go
q := gormx.QAnd(
    gormx.QOr(
        gormx.Eq("type", "A"),
        gormx.QAnd(
            gormx.Eq("type", "B"),
            gormx.Gt("value", 100),
        ),
    ),
    gormx.Eq("active", true),
)

// (type = 'A' OR (type = 'B' AND value > 100)) AND active = true
```

### مثال 7: دمج مع Filter

```go
q := gormx.QOr(
    gormx.Eq("role", "admin"),
    gormx.Eq("role", "moderator"),
)

users, _ := gormx.New[User]().
    Filter("active", true).       // AND
    Q(q).                          // AND (role = 'admin' OR role = 'moderator')
    Filter("age__gte", 18).
    All()
// active = true AND (role = 'admin' OR role = 'moderator') AND age >= 18
```

---

## 🔧 Methods

### `And(children...)`

```go
q := gormx.Qb().
    And(gormx.Eq("a", 1)).
    And(gormx.Eq("b", 2))
// a = 1 AND b = 2
```

### `Or(children...)`

```go
q := gormx.Qb().
    Or(gormx.Eq("a", 1)).
    Or(gormx.Eq("b", 2))
// a = 1 OR b = 2
```

⚠️ **مهم**: كلاهما **immutable**:

```go
base := gormx.Qb().And(gormx.Eq("a", 1))
extended := base.And(gormx.Eq("b", 2))

// base لم يتغير
```

### `AndGroup(children...)`

```go
q := gormx.Qb().
    And(gormx.Eq("a", 1)).
    AndGroup(
        gormx.Eq("b", 2),
        gormx.Eq("c", 3),
    )
// a = 1 AND (b = 2 AND c = 3)
```

### `OrGroup(children...)`

```go
q := gormx.Qb().
    And(gormx.Eq("a", 1)).
    OrGroup(
        gormx.Eq("b", 2),
        gormx.Eq("c", 3),
    )
// a = 1 OR (b = 2 OR c = 3)
```

### `ToSQL()`

```go
sql, args := q.ToSQL()
fmt.Println(sql)
```

---

## 🎯 مثال واقعي

### Search API

```go
func SearchUsers(filters SearchFilters) ([]User, error) {
    q := gormx.New[User]().Filter("active", true)

    // بحث نصي (OR على عدة أعمدة)
    if filters.Query != "" {
        searchQ := gormx.QOr(
            gormx.Contains("name", filters.Query),
            gormx.Contains("email", filters.Query),
            gormx.Contains("phone", filters.Query),
        )
        q = q.Q(searchQ)
    }

    // فلترة حسب الدور
    if len(filters.Roles) > 0 {
        q = q.Filter("role__in", filters.Roles)
    }

    // فلترة حسب العمر
    if filters.MinAge > 0 {
        q = q.Filter("age__gte", filters.MinAge)
    }
    if filters.MaxAge > 0 {
        q = q.Filter("age__lte", filters.MaxAge)
    }

    // حالة خاصة: admin أو verified
    if filters.SpecialOnly {
        specialQ := gormx.QOr(
            gormx.Eq("role", "admin"),
            gormx.Eq("verified", true),
        )
        q = q.Q(specialQ)
    }

    return q.OrderBy("-created_at").All()
}
```

---

## ⚠️ Performance

- `OR` مع عدة أعمدة → لا يستخدم الفهارس جيدًا
- `NOT` قد يكون بطيئًا
- استخدم `IN` بدل `OR` عند الإمكان:

```go
// ❌ بطيء
q := gormx.QOr(
    gormx.Eq("status", "active"),
    gormx.Eq("status", "pending"),
    gormx.Eq("status", "verified"),
)

// ✅ أسرع
q := gormx.In("status", []any{"active", "pending", "verified"})
```

---

## 📝 Best Practices

### ✅ Do

- استخدم `IN` بدل `OR` متكرر على نفس العمود
- استخدم `Q` للشروط المترابطة
- اختبر SQL النهائي بـ `ToSQL()`

### ❌ Don't

- لا تخلط `And`/`Or` بدون grouping
- لا تنسَ parenthesization
- لا تستخدم `NOT` مع فهارس