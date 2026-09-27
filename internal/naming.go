package internal

import (
	"strings"
	"unicode"
)

// ═══════════════════════════════════════════════
// Naming — تحويلات أسماء
// ═══════════════════════════════════════════════

// ToSnake يحوّل CamelCase إلى snake_case.
//
//	ToSnake("FirstName")  → "first_name"
//	ToSnake("UserID")     → "user_id"
//	ToSnake("HTTPServer") → "http_server"
//	ToSnake("APIKey")     → "api_key"
//	ToSnake("ID")         → "id"
//	ToSnake("")           → ""
func ToSnake(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			// أضف _ قبل حرف كبير إذا:
			// - ليس الأول
			// - الحرف السابق صغير أو رقم
			// - أو الحرف التالي صغير (لكلمات مثل "HTTPServer")
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					b.WriteRune('_')
				} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					b.WriteRune('_')
				}
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}


// commonInitialisms قائمة الاختصارات الشائعة (من golint).
//
// تُستخدم لتحويل snake_case → camelCase / PascalCase بشكل صحيح:
//
//	user_id  → userID  (وليس userId)
//	api_key  → apiKey  (وليس apikey - كلمة عادية)
//
// المرجع: https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

// capitalizeWord يُعيد الكلمة بحرف كبير أول، مع احترام الاختصارات.
//
//	capitalizeWord("id")   → "ID"
//	capitalizeWord("name") → "Name"
//	capitalizeWord("api")  → "API"
func capitalizeWord(s string) string {
	if s == "" {
		return ""
	}
	upper := strings.ToUpper(s)
	if commonInitialisms[upper] {
		return upper
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// ToCamel يحوّل snake_case إلى camelCase.
//
//	ToCamel("first_name") → "firstName"
//	ToCamel("user_id")    → "userID"
//	ToCamel("")           → ""
func ToCamel(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	var b strings.Builder
	b.WriteString(parts[0])

	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		b.WriteString(capitalizeWord(p))
	}

	return b.String()
}

// ToPascal يحوّل snake_case إلى PascalCase.
//
//	ToPascal("first_name") → "FirstName"
//	ToPascal("user_id")    → "UserID"
func ToPascal(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	var b strings.Builder

	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(capitalizeWord(p))
	}

	return b.String()
}

// ═══════════════════════════════════════════════
// Pluralization
// ═══════════════════════════════════════════════

// irregularPlurals كلمات شاذة.
var irregularPlurals = map[string]string{
	"person": "people",
	"man":    "men",
	"woman":  "women",
	"child":  "children",
	"tooth":  "teeth",
	"foot":   "feet",
	"mouse":  "mice",
	"goose":  "geese",
	"ox":     "oxen",
	"datum":  "data",
	"medium": "media",
	"index":  "indices",
	"matrix": "matrices",
	"vertex": "vertices",
}

// uncountable كلمات غير معدودة.
var uncountable = map[string]bool{
	"information": true,
	"equipment":   true,
	"money":       true,
	"rice":        true,
	"series":      true,
	"species":     true,
	"fish":        true,
	"sheep":       true,
	"deer":        true,
	"data":        true,
	"media":       true,
}

// Pluralize يضيف "s" أو "es" حسب الكلمة.
//
//	Pluralize("user")   → "users"
//	Pluralize("box")    → "boxes"
//	Pluralize("city")   → "cities"
//	Pluralize("person") → "people"
//	Pluralize("sheep")  → "sheep"
func Pluralize(s string) string {
	if s == "" {
		return ""
	}

	lower := strings.ToLower(s)

	// شاذ
	if p, ok := irregularPlurals[lower]; ok {
		return preserveCase(s, p)
	}

	// غير معدود
	if uncountable[lower] {
		return s
	}

	// قواعد
	switch {
	case strings.HasSuffix(lower, "y") && len(s) > 1 && !isVowel(lower[len(lower)-2]):
		return s[:len(s)-1] + "ies"
	case strings.HasSuffix(lower, "s"),
		strings.HasSuffix(lower, "x"),
		strings.HasSuffix(lower, "z"),
		strings.HasSuffix(lower, "ch"),
		strings.HasSuffix(lower, "sh"):
		return s + "es"
	case strings.HasSuffix(lower, "f"):
		return s[:len(s)-1] + "ves"
	case strings.HasSuffix(lower, "fe"):
		return s[:len(s)-2] + "ves"
	default:
		return s + "s"
	}
}

// Singularize يحوّل للفرد (تقريبي).
//
//	Singularize("users")   → "user"
//	Singularize("boxes")   → "box"
//	Singularize("cities")  → "city"
//	Singularize("people")  → "person"
func Singularize(s string) string {
	if s == "" {
		return ""
	}

	lower := strings.ToLower(s)

	// ابحث في irregular
	for k, v := range irregularPlurals {
		if strings.ToLower(v) == lower {
			return preserveCase(s, k)
		}
	}

	switch {
	case strings.HasSuffix(lower, "ies"):
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(lower, "ves"):
		return s[:len(s)-3] + "f"
	case strings.HasSuffix(lower, "es"):
		return s[:len(s)-2]
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss"):
		return s[:len(s)-1]
	}
	return s
}

// ═══════════════════════════════════════════════
// Helpers
// ═══════════════════════════════════════════════

func isVowel(b byte) bool {
	return b == 'a' || b == 'e' || b == 'i' || b == 'o' || b == 'u'
}

// preserveCase يحافظ على حالة الأحرف الأصلية.
func preserveCase(original, replacement string) string {
	if original == "" {
		return replacement
	}
	if unicode.IsUpper(rune(original[0])) {
		return strings.ToUpper(replacement[:1]) + replacement[1:]
	}
	return replacement
}