package internal

import (
	"fmt"
	"reflect"
	"strings"
)

// ═══════════════════════════════════════════════
// Lookups — معالجة __lookup
// ═══════════════════════════════════════════════

// Lookup يمثل lookup واحد.
type Lookup struct {
	SQL     string
	Args    []any
	Negate  bool
	Dialect Dialect
}

// BuildLookup يحوّل "field__op" إلى SQL + args.
//
// يستخدم DialectSQLite افتراضيًا.
func BuildLookup(field string, value any, negate bool) (*Lookup, error) {
	return BuildLookupWithDialect(field, value, negate, DialectSQLite)
}

// BuildLookupWithDialect مثل BuildLookup لكن مع dialect.
func BuildLookupWithDialect(field string, value any, negate bool, dialect Dialect) (*Lookup, error) {
	col, op := SplitFieldLookup(field)

	if err := ValidateField(col); err != nil {
		return nil, err
	}

	if op == "" {
		return &Lookup{
			SQL:     col + " = ?",
			Args:    []any{value},
			Negate:  negate,
			Dialect: dialect,
		}, nil
	}

	if err := ValidateLookup(op); err != nil {
		return nil, err
	}

	switch op {
	case "gt":
		return newLookup(col+" > ?", value, negate, dialect), nil
	case "gte":
		return newLookup(col+" >= ?", value, negate, dialect), nil
	case "lt":
		return newLookup(col+" < ?", value, negate, dialect), nil
	case "lte":
		return newLookup(col+" <= ?", value, negate, dialect), nil
	case "ne":
		return newLookup(col+" != ?", value, negate, dialect), nil

	case "contains":
		return newLookup(col+" LIKE ?", "%"+toString(value)+"%", negate, dialect), nil
	case "icontains":
		return newLookup(lowerFunc(dialect, col)+" LIKE ?",
			"%"+strings.ToLower(toString(value))+"%", negate, dialect), nil
	case "startswith":
		return newLookup(col+" LIKE ?", toString(value)+"%", negate, dialect), nil
	case "istartswith":
		return newLookup(lowerFunc(dialect, col)+" LIKE ?",
			strings.ToLower(toString(value))+"%", negate, dialect), nil
	case "endswith":
		return newLookup(col+" LIKE ?", "%"+toString(value), negate, dialect), nil
	case "iendswith":
		return newLookup(lowerFunc(dialect, col)+" LIKE ?",
			"%"+strings.ToLower(toString(value)), negate, dialect), nil

	case "in":
		return buildInLookup(col, value, negate, dialect)
	case "notin":
		return buildInLookup(col, value, !negate, dialect)

	case "isnull":
		isNull, ok := value.(bool)
		if !ok {
			return nil, NewValidationError(field, "isnull expects bool")
		}
		if isNull {
			return newLookup(col+" IS NULL", nil, negate, dialect), nil
		}
		return newLookup(col+" IS NOT NULL", nil, negate, dialect), nil

	case "between":
		rv := reflect.ValueOf(value)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return nil, NewValidationError(field, "between expects slice/array")
		}
		if rv.Len() != 2 {
			return nil, NewValidationError(field, "between expects exactly 2 elements")
		}
		return newLookup(
			col+" BETWEEN ? AND ?",
			[]any{rv.Index(0).Interface(), rv.Index(1).Interface()},
			negate,
			dialect,
		), nil

	// Date lookups
	case "year":
		return newLookup(dateExtract(dialect, "year", col)+" = ?", value, negate, dialect), nil
	case "month":
		return newLookup(dateExtract(dialect, "month", col)+" = ?", value, negate, dialect), nil
	case "day":
		return newLookup(dateExtract(dialect, "day", col)+" = ?", value, negate, dialect), nil

	default:
		return nil, NewValidationError(field, fmt.Sprintf("unknown lookup %q", op))
	}
}

func newLookup(sql string, args any, negate bool, dialect Dialect) *Lookup {
	l := &Lookup{SQL: sql, Negate: negate, Dialect: dialect}
	if args == nil {
		return l
	}
	switch v := args.(type) {
	case []any:
		l.Args = v
	default:
		l.Args = []any{args}
	}
	return l
}

func buildInLookup(col string, value any, negate bool, dialect Dialect) (*Lookup, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, NewValidationError("in", "expects slice/array")
	}

	n := rv.Len()
	if n == 0 {
		// IN () دائما false
		return &Lookup{SQL: "1=0", Negate: negate, Dialect: dialect}, nil
	}

	args := make([]any, n)
	for i := 0; i < n; i++ {
		args[i] = rv.Index(i).Interface()
	}

	placeholders := strings.Repeat("?,", n)
	placeholders = placeholders[:len(placeholders)-1]

	return &Lookup{
		SQL:     col + " IN (" + placeholders + ")",
		Args:    args,
		Negate:  negate,
		Dialect: dialect,
	}, nil
}

// lowerFunc يرجّع دالة lower حسب dialect.
func lowerFunc(dialect Dialect, col string) string {
	// جميع dialects تدعم LOWER()
	return "LOWER(" + col + ")"
}

// dateExtract يرجّع تعبير استخراج التاريخ حسب dialect.
func dateExtract(dialect Dialect, part, col string) string {
	switch dialect {
	case DialectSQLite:
		format := map[string]string{
			"year":  "%Y",
			"month": "%m",
			"day":   "%d",
		}[part]
		return "CAST(strftime('" + format + "', " + col + ") AS INTEGER)"

	case DialectPostgres:
		return "EXTRACT(" + strings.ToUpper(part) + " FROM " + col + ")"

	case DialectMySQL:
		switch part {
		case "year":
			return "YEAR(" + col + ")"
		case "month":
			return "MONTH(" + col + ")"
		case "day":
			return "DAY(" + col + ")"
		}
		return "EXTRACT(" + strings.ToUpper(part) + " FROM " + col + ")"

	default:
		return "EXTRACT(" + strings.ToUpper(part) + " FROM " + col + ")"
	}
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}