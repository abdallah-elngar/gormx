// Package internal provides internal helpers for gormx.
package internal

// ═══════════════════════════════════════════════
// Clauses — أنواع الشروط الداخلية
// ═══════════════════════════════════════════════

// RawClause يمثل شرط SQL خام.
//
// يُستخدم داخل gormx لتمرير SQL + args بدون validation.
//
// ⚠️ للاستخدام الداخلي فقط.
type RawClause struct {
	SQL  string
	Args []any
}

// NewRawClause ينشئ RawClause.
func NewRawClause(sql string, args ...any) RawClause {
	return RawClause{SQL: sql, Args: args}
}

// IsEmpty يفحص إذا كان الشرط فارغًا.
func (r RawClause) IsEmpty() bool {
	return r.SQL == ""
}

// String يرجّع تمثيل نصي.
func (r RawClause) String() string {
	if len(r.Args) == 0 {
		return r.SQL
	}
	return r.SQL
}