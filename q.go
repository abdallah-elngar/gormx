package gormx

import (
	"strings"

	"github.com/abdallah-elngar/gormx/internal"
)

// ═══════════════════════════════════════════════
// Q Builder — لبناء شروط معقدة (AND / OR / NOT)
// ═══════════════════════════════════════════════

// Q يمثل مجموعة شروط.
//
// مثال:
//
//	q := gormx.QOr(
//	    gormx.Eq("status", "active"),
//	    gormx.Eq("status", "pending"),
//	)
//	users, _ := gormx.New[User]().Q(q).All()
type Q struct {
	op       string
	children []any
}

// ToSQL يحوّل Q إلى SQL + args (للتصحيح).
func (q *Q) ToSQL() (string, []any) {
	return q.toSQL()
}

// toSQL داخلي.
func (q *Q) toSQL() (string, []any) {
	if q == nil || len(q.children) == 0 {
		return "1=1", nil
	}

	parts := make([]string, 0, len(q.children))
	var args []any

	for _, c := range q.children {
		switch v := c.(type) {
		case whereClause:
			parts = append(parts, v.sql)
			args = append(args, v.args...)
		case notClause:
			parts = append(parts, "NOT ("+v.sql+")")
			args = append(args, v.args...)
		case internal.RawClause:
			parts = append(parts, v.SQL)
			args = append(args, v.Args...)
		case *Q:
			sql, a := v.toSQL()
			parts = append(parts, "("+sql+")")
			args = append(args, a...)
		}
	}

	sep := " " + q.op + " "
	return strings.Join(parts, sep), args
}

// deepCopy نسخة عميقة للـ Q.
func (q *Q) deepCopy() *Q {
	if q == nil {
		return nil
	}
	nq := &Q{op: q.op, children: make([]any, len(q.children))}
	for i, c := range q.children {
		if sub, ok := c.(*Q); ok {
			nq.children[i] = sub.deepCopy()
		} else {
			nq.children[i] = c
		}
	}
	return nq
}

// ═══════════════════════════════════════════════
// Clauses — خاصة بـ gormx
// ═══════════════════════════════════════════════

type whereClause struct {
	sql  string
	args []any
}

type notClause struct {
	sql  string
	args []any
}

// ═══════════════════════════════════════════════
// Public Lookup Constructors
// ═══════════════════════════════════════════════
//
// ملاحظة: كل الدوال panic على حقل غير صحيح.
// استخدم الإصدارات *Err للتحكم.

// Eq → field = value. Panics on invalid field.
func Eq(field string, value any) whereClause {
	c, _ := EqErr(field, value)
	return c
}

// EqErr مثل Eq لكن يرجّع خطأ.
func EqErr(field string, value any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " = ?", args: []any{value}}, nil
}

// Ne → field != value. Panics.
func Ne(field string, value any) whereClause {
	c, _ := NeErr(field, value)
	return c
}

// NeErr مثل Ne.
func NeErr(field string, value any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " != ?", args: []any{value}}, nil
}

// Gt → field > value.
func Gt(field string, value any) whereClause {
	c, _ := GtErr(field, value)
	return c
}

// GtErr مثل Gt.
func GtErr(field string, value any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " > ?", args: []any{value}}, nil
}

// Gte → field >= value.
func Gte(field string, value any) whereClause {
	c, _ := GteErr(field, value)
	return c
}

// GteErr مثل Gte.
func GteErr(field string, value any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " >= ?", args: []any{value}}, nil
}

// Lt → field < value.
func Lt(field string, value any) whereClause {
	c, _ := LtErr(field, value)
	return c
}

// LtErr مثل Lt.
func LtErr(field string, value any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " < ?", args: []any{value}}, nil
}

// Lte → field <= value.
func Lte(field string, value any) whereClause {
	c, _ := LteErr(field, value)
	return c
}

// LteErr مثل Lte.
func LteErr(field string, value any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " <= ?", args: []any{value}}, nil
}

// Contains → field LIKE '%value%'.
func Contains(field, value string) whereClause {
	c, _ := ContainsErr(field, value)
	return c
}

// ContainsErr مثل Contains.
func ContainsErr(field, value string) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " LIKE ?", args: []any{"%" + value + "%"}}, nil
}

// StartsWith → field LIKE 'value%'.
func StartsWith(field, value string) whereClause {
	c, _ := StartsWithErr(field, value)
	return c
}

// StartsWithErr مثل StartsWith.
func StartsWithErr(field, value string) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " LIKE ?", args: []any{value + "%"}}, nil
}

// EndsWith → field LIKE '%value'.
func EndsWith(field, value string) whereClause {
	c, _ := EndsWithErr(field, value)
	return c
}

// EndsWithErr مثل EndsWith.
func EndsWithErr(field, value string) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " LIKE ?", args: []any{"%" + value}}, nil
}

// In → field IN (values...).
func In(field string, values []any) whereClause {
	c, _ := InErr(field, values)
	return c
}

// InErr مثل In.
func InErr(field string, values []any) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	if len(values) == 0 {
		return whereClause{sql: "1=0"}, nil
	}
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1]
	return whereClause{sql: field + " IN (" + placeholders + ")", args: values}, nil
}

// IsNull → field IS NULL.
func IsNull(field string) whereClause {
	c, _ := IsNullErr(field)
	return c
}

// IsNullErr مثل IsNull.
func IsNullErr(field string) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " IS NULL"}, nil
}

// NotNull → field IS NOT NULL.
func NotNull(field string) whereClause {
	c, _ := NotNullErr(field)
	return c
}

// NotNullErr مثل NotNull.
func NotNullErr(field string) (whereClause, error) {
	if err := internal.ValidateField(field); err != nil {
		return whereClause{}, err
	}
	return whereClause{sql: field + " IS NOT NULL"}, nil
}

// Raw → SQL خام. لا يتحقق من الحقل.
//
// ⚠️ المسؤولية على المستخدم — لا يوجد validation.
//
// مثال:
//
//	gormx.Raw("age > ? AND status = ?", 18, "active")
func Raw(sql string, args ...any) internal.RawClause {
	return internal.RawClause{SQL: sql, Args: args}
}

// Not → نفي الشرط.
//
// يدعم: whereClause, internal.RawClause, *Q
func Not(c any) notClause {
	switch v := c.(type) {
	case whereClause:
		return notClause{sql: v.sql, args: v.args}
	case internal.RawClause:
		return notClause{sql: v.SQL, args: v.Args}
	case *Q:
		sql, args := v.toSQL()
		return notClause{sql: sql, args: args}
	}
	return notClause{sql: "1=1"}
}

// ═══════════════════════════════════════════════
// Q Constructors
// ═══════════════════════════════════════════════

// Qb ينشئ Q جديد بـ AND.
func Qb() *Q {
	return &Q{op: "AND"}
}

// QOr ينشئ Q جديد بـ OR.
//
//	q := gormx.QOr(
//	    gormx.Eq("status", "active"),
//	    gormx.Eq("status", "pending"),
//	)
func QOr(children ...any) *Q {
	return &Q{op: "OR", children: children}
}

// QAnd ينشئ Q جديد بـ AND.
func QAnd(children ...any) *Q {
	return &Q{op: "AND", children: children}
}

// And يضيف شروط بـ AND.
//
// يعيد Q جديد (immutable).
func (q *Q) And(children ...any) *Q {
	nq := q.deepCopy()
	if nq == nil {
		nq = &Q{op: "AND"}
	}
	nq.op = "AND"
	nq.children = append(nq.children, children...)
	return nq
}

// Or يضيف شروط بـ OR.
//
// يعيد Q جديد (immutable).
func (q *Q) Or(children ...any) *Q {
	nq := q.deepCopy()
	if nq == nil {
		nq = &Q{op: "OR"}
	}
	nq.op = "OR"
	nq.children = append(nq.children, children...)
	return nq
}

// AndGroup يضيف مجموعة AND متداخلة.
//
//	q.AndGroup(
//	    gormx.Eq("a", 1),
//	    gormx.Eq("b", 2),
//	)
//	→ ... AND (a = ? AND b = ?)
func (q *Q) AndGroup(children ...any) *Q {
	return q.And(QAnd(children...))
}

// OrGroup يضيف مجموعة OR متداخلة.
//
//	q.OrGroup(
//	    gormx.Eq("a", 1),
//	    gormx.Eq("b", 2),
//	)
//	→ ... OR (a = ? OR b = ?)
func (q *Q) OrGroup(children ...any) *Q {
	return q.Or(QOr(children...))
}

// Len يرجّع عدد الشروط.
func (q *Q) Len() int {
	if q == nil {
		return 0
	}
	return len(q.children)
}

// IsEmpty يفحص إذا كان Q فارغًا.
func (q *Q) IsEmpty() bool {
	return q == nil || len(q.children) == 0
}