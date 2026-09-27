// Package internal provides internal helpers for gormx.
package internal

import (
	"fmt"
	"regexp"
	"strings"
)

// ═══════════════════════════════════════════════
// Field Validation — حماية ضد SQL injection
// ═══════════════════════════════════════════════

var (
	// fieldRegex يطابق أسماء الحقول الصحيحة.
	//
	// يسمح بـ:
	//   - snake_case: user_id
	//   - dotted: users.id
	//   - prefixed: t1.name
	//   - numbers: field_1
	fieldRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)*$`)

	// lookupRegex يطابق أسماء lookups.
	lookupRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

	// sqlKeywords كلمات محجوزة (لا يمكن استخدامها كأسماء حقول).
	sqlKeywords = map[string]bool{
		"select": true, "from": true, "where": true, "insert": true,
		"update": true, "delete": true, "drop": true, "table": true,
		"union": true, "or": true, "and": true, "not": true,
		"null": true, "true": true, "false": true,
		"alter": true, "create": true, "truncate": true, "grant": true,
		"revoke": true, "exec": true, "execute": true, "declare": true,
		"primary": true, "foreign": true, "references": true, "index": true,
	}
)

// ═══════════════════════════════════════════════
// Field Validation
// ═══════════════════════════════════════════════

// ValidateField يتحقق من صحة اسم حقل واحد.
//
// يعيد ValidationError إذا كان الحقل غير صالح.
func ValidateField(field string) error {
	if field == "" {
		return NewValidationError(field, "empty field name")
	}
	if len(field) > 128 {
		return NewValidationError(field, "field name too long (max 128)")
	}
	if !fieldRegex.MatchString(field) {
		return NewValidationError(field, "invalid characters (allowed: a-z, A-Z, 0-9, _, .)")
	}

	// فحص الكلمات المحجوزة (للجزء الأول فقط)
	base := field
	if i := strings.Index(field, "."); i > 0 {
		base = field[:i]
	}
	if sqlKeywords[strings.ToLower(base)] {
		return NewValidationError(field, "reserved SQL keyword")
	}
	return nil
}

// ValidateFields يتحقق من عدة حقول.
func ValidateFields(fields ...string) error {
	for _, f := range fields {
		if err := ValidateField(f); err != nil {
			return err
		}
	}
	return nil
}

// MustValidateField يتحقق وpanic عند الخطأ.
//
// للاستخدام الداخلي عندما نعرف أن الحقل صالح.
func MustValidateField(field string) {
	if err := ValidateField(field); err != nil {
		panic(err)
	}
}

// ═══════════════════════════════════════════════
// Lookup Validation
// ═══════════════════════════════════════════════

// validLookups قائمة lookups الصحيحة.
var validLookups = map[string]bool{
	"gt": true, "gte": true, "lt": true, "lte": true, "ne": true,
	"contains": true, "icontains": true,
	"startswith": true, "istartswith": true,
	"endswith": true, "iendswith": true,
	"in": true, "notin": true, "isnull": true, "between": true,
	"year": true, "month": true, "day": true,
	"regex": true, "iregex": true,
}

// ValidateLookup يتحقق من صحة اسم lookup.
func ValidateLookup(lookup string) error {
	if lookup == "" {
		return NewValidationError(lookup, "empty lookup")
	}
	if !lookupRegex.MatchString(lookup) {
		return NewValidationError(lookup, "invalid lookup name")
	}
	if !validLookups[lookup] {
		return NewValidationError(lookup, "unknown lookup")
	}
	return nil
}

// IsValidLookup يفحص إذا كان lookup صالحًا.
func IsValidLookup(lookup string) bool {
	return validLookups[lookup]
}

// ═══════════════════════════════════════════════
// Operator Validation
// ═══════════════════════════════════════════════

// validOperators قائمة operators الصحيحة.
var validOperators = map[string]bool{
	"=": true, "!=": true, "<>": true,
	">": true, ">=": true, "<": true, "<=": true,
	"LIKE": true, "NOT LIKE": true, "ILIKE": true,
	"IN": true, "NOT IN": true,
	"IS NULL": true, "IS NOT NULL": true,
	"BETWEEN": true, "NOT BETWEEN": true,
}

// ValidateOperator يتحقق من صحة operator.
func ValidateOperator(op string) error {
	if !validOperators[strings.ToUpper(op)] {
		return NewValidationError(op, "invalid operator")
	}
	return nil
}

// ═══════════════════════════════════════════════
// SplitFieldLookup
// ═══════════════════════════════════════════════

// SplitFieldLookup يقسم "field__lookup" إلى جزأيه.
//
//	SplitFieldLookup("age__gt")    → ("age", "gt")
//	SplitFieldLookup("name")        → ("name", "")
//	SplitFieldLookup("user__email__contains") → ("user__email", "contains")
func SplitFieldLookup(s string) (field, lookup string) {
	// ابحث عن آخر "__"
	if i := strings.LastIndex(s, "__"); i > 0 {
		return s[:i], s[i+2:]
	}
	return s, ""
}

// ═══════════════════════════════════════════════
// Direction Sanitization
// ═══════════════════════════════════════════════

// SanitizeDirection يتحقق من صحة direction.
//
//	SanitizeDirection("asc")  → ("ASC", nil)
//	SanitizeDirection("DESC") → ("DESC", nil)
//	SanitizeDirection("bad")  → ("", error)
func SanitizeDirection(dir string) (string, error) {
	d := strings.ToUpper(strings.TrimSpace(dir))
	switch d {
	case "ASC", "DESC", "":
		return d, nil
	}
	return "", NewValidationError(dir, "invalid direction (allowed: ASC, DESC)")
}

// ═══════════════════════════════════════════════
// Table Name Validation
// ═══════════════════════════════════════════════

// tableRegex يطابق أسماء الجداول الصحيحة.
var tableRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ValidateTableName يتحقق من صحة اسم جدول.
func ValidateTableName(table string) error {
	if table == "" {
		return NewValidationError(table, "empty table name")
	}
	if len(table) > 64 {
		return NewValidationError(table, "table name too long (max 64)")
	}
	if !tableRegex.MatchString(table) {
		return NewValidationError(table, "invalid table name")
	}
	if sqlKeywords[strings.ToLower(table)] {
		return NewValidationError(table, "reserved SQL keyword")
	}
	return nil
}

// ═══════════════════════════════════════════════
// ValidationError Type
// ═══════════════════════════════════════════════

// ValidationError خطأ تحقق.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("gormx: invalid field %q: %s", e.Field, e.Reason)
}

// NewValidationError ينشئ خطأ تحقق.
func NewValidationError(field, reason string) *ValidationError {
	return &ValidationError{Field: field, Reason: reason}
}