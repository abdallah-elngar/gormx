# Example 1: Basic

> أول برنامج مع gormx.

---

## 📦 Setup

```go
package main

import (
    "fmt"
    "log"

    "github.com/sanad/gormx"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:255"`
    Email string `gorm:"uniqueIndex"`
    Age   int
}

func main() {
    // 1. Connect
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 2. Set global
    gormx.SetDB(db)

    // 3. Migrate
    if err := gormx.Migrate[User](); err != nil {
        log.Fatal(err)
    }

    fmt.Println("✅ Ready")

    // 4. Run examples
    runBasic()
    runFilters()
    runUpdates()
}

func runBasic() {
    fmt.Println("\n--- Basic ---")

    // Create
    user := &User{Name: "Ali", Email: "ali@test.com", Age: 30}
    if err := gormx.New[User]().Create(user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created: ID=%d\n", user.ID)

    // Read
    got, err := gormx.New[User]().Get(user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Got: %s\n", got.Name)

    // List
    all, _ := gormx.New[User]().All()
    fmt.Printf("Total: %d\n", len(all))
}

func runFilters() {
    fmt.Println("\n--- Filters ---")

    // Seed
    gormx.New[User]().CreateMany([]User{
        {Name: "Sara", Email: "sara@test.com", Age: 25},
        {Name: "Omar", Email: "omar@test.com", Age: 17},
    })

    // Filter
    adults, _ := gormx.New[User]().Filter("age__gte", 18).All()
    fmt.Printf("Adults: %d\n", len(adults))

    // Contains
    ali, _ := gormx.New[User]().Filter("name__icontains", "ali").All()
    fmt.Printf("Ali-like: %d\n", len(ali))
}

func runUpdates() {
    fmt.Println("\n--- Updates ---")

    // Update
    err := gormx.New[User]().Update(1, "age", 31)
    if err != nil {
        log.Fatal(err)
    }

    user, _ := gormx.New[User]().Get(1)
    fmt.Printf("New age: %d\n", user.Age)

    // Delete
    err = gormx.New[User]().Delete(1)
    if err != nil {
        log.Fatal(err)
    }

    count, _ := gormx.New[User]().Count()
    fmt.Printf("Remaining: %d\n", count)
}
```

---

## 📊 Output

```
✅ Ready

--- Basic ---
Created: ID=1
Got: Ali
Total: 1

--- Filters ---
Adults: 2
Ali-like: 1

--- Updates ---
New age: 31
Remaining: 2
```

---

## 🎯 النقاط الرئيسية

- `SetDB` قبل أي شيء
- `Migrate` لإنشاء الجداول
- `New[T]()` للاستعلامات
- `Create`, `Get`, `Update`, `Delete` أساسيات