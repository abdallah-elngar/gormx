package internal

import (
	"reflect"
	"strings"
	"sync"
)

// ═══════════════════════════════════════════════
// Reflection — أدوات انعكاس
// ═══════════════════════════════════════════════

var (
	modelCache   = make(map[reflect.Type]*ModelInfo)
	modelCacheMu sync.RWMutex
)

// ModelInfo معلومات مخزّنة عن موديل.
type ModelInfo struct {
	Type       reflect.Type
	Name       string // اسم الـ Go struct
	TableName  string // اسم الجدول
	Fields     []FieldInfo
	FieldIndex map[string]int // gorm column → index
}

// FieldInfo معلومات حقل.
type FieldInfo struct {
	GoName    string
	Column    string
	Type      reflect.Type
	IsPrimary bool
	IsAuto    bool
	IsTime    bool
	JSONName  string
}

// GetModelInfo يستخرج معلومات الموديل (مع caching).
func GetModelInfo[T any]() *ModelInfo {
	var zero T
	t := reflect.TypeOf(zero)

	// فك الـ pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return &ModelInfo{Type: t, Name: t.Name()}
	}

	modelCacheMu.RLock()
	if info, ok := modelCache[t]; ok {
		modelCacheMu.RUnlock()
		return info
	}
	modelCacheMu.RUnlock()

	info := extractModelInfo(t)

	modelCacheMu.Lock()
	modelCache[t] = info
	modelCacheMu.Unlock()

	return info
}

func extractModelInfo(t reflect.Type) *ModelInfo {
	info := &ModelInfo{
		Type:       t,
		Name:       t.Name(),
		TableName:  TableNameFromType(t),
		FieldIndex: make(map[string]int),
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// تخطي unexported
		if !f.IsExported() {
			continue
		}

		fi := FieldInfo{
			GoName: f.Name,
			Column: getColumnName(f),
			Type:   f.Type,
			IsTime: f.Type.String() == "time.Time",
		}

		// فحص tags
		gormTag := f.Tag.Get("gorm")
		if strings.Contains(gormTag, "primaryKey") {
			fi.IsPrimary = true
		}
		if strings.Contains(gormTag, "autoIncrement") {
			fi.IsAuto = true
		}

		// JSON tag
		jsonTag := f.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			fi.JSONName = parts[0]
		}

		info.Fields = append(info.Fields, fi)
		info.FieldIndex[fi.Column] = i
	}

	return info
}

func getColumnName(f reflect.StructField) string {
	gormTag := f.Tag.Get("gorm")
	for _, part := range strings.Split(gormTag, ";") {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ToSnake(f.Name)
}

// FieldNames يرجّع أسماء الأعمدة لموديل.
func FieldNames[T any]() []string {
	info := GetModelInfo[T]()
	names := make([]string, len(info.Fields))
	for i, f := range info.Fields {
		names[i] = f.Column
	}
	return names
}

// ModelName يرجّع اسم الموديل.
func ModelName[T any]() string {
	return GetModelInfo[T]().Name
}

// TableName يرجّع اسم الجدول.
func TableName[T any]() string {
	return GetModelInfo[T]().TableName
}

// TableName يستخرج اسم الجدول من reflect.Type.
//
// القواعد:
//  1. إذا كان هناك TableName() string method → استخدمه
//  2. وإلا → pluralize(to_snake(name))
func TableNameFromType(t reflect.Type) string {
	// فك الـ pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return ""
	}

	// فحص TableName() method
	if m, ok := t.MethodByName("TableName"); ok {
		// يجب أن يكون: func() string
		if m.Type.NumIn() == 1 && m.Type.NumOut() == 1 {
			if m.Type.Out(0).Kind() == reflect.String {
				// استدعِ Method على instance
				zero := reflect.New(t).Interface()
				if tn, ok := zero.(interface{ TableName() string }); ok {
					return tn.TableName()
				}
			}
		}
	}

	// وإلا → pluralize
	return Pluralize(ToSnake(t.Name()))
}

// ClearModelCache يمسح الـ cache (مفيد في الاختبارات).
func ClearModelCache() {
	modelCacheMu.Lock()
	modelCache = make(map[reflect.Type]*ModelInfo)
	modelCacheMu.Unlock()
}