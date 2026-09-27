package advanced

import (
	"fmt"
	"strings"

	"github.com/abdallah-elngar/gormx"
)

// ═══════════════════════════════════════════════
// CTE — Common Table Expressions
// ═══════════════════════════════════════════════

// CTE يمثل Common Table Expression.
type CTE struct {
	Name      string
	SQL       string
	Args      []any
	Recursive bool
	Columns   []string
}

// NewCTE ينشئ CTE من QuerySet.
func NewCTE[T any](name string, q *gormx.QuerySet[T]) *CTE {
	if q == nil {
		return &CTE{Name: name}
	}
	sql, args := q.ToSQL()
	return &CTE{
		Name: name,
		SQL:  sql,
		Args: args,
	}
}

// NewRecursiveCTE ينشئ recursive CTE.
func NewRecursiveCTE[T any](name string, base *gormx.QuerySet[T]) *CTE {
	if base == nil {
		return &CTE{Name: name, Recursive: true}
	}
	sql, args := base.ToSQL()
	return &CTE{
		Name:      name,
		SQL:       sql,
		Args:      args,
		Recursive: true,
	}
}

// WithColumns يحدد أسماء الأعمدة.
func (c *CTE) WithColumns(cols ...string) *CTE {
	c.Columns = cols
	return c
}

// UnionRaw يضيف UNION ALL مع SQL خام.
func (c *CTE) UnionRaw(sql string, args ...any) *CTE {
	c.SQL = c.SQL + " UNION ALL " + sql
	c.Args = append(c.Args, args...)
	return c
}

// ToSQL يحوّل CTE إلى SQL.
func (c *CTE) ToSQL() (string, []any) {
	var b strings.Builder

	if c.Recursive {
		b.WriteString("WITH RECURSIVE ")
	} else {
		b.WriteString("WITH ")
	}

	b.WriteString(c.Name)

	if len(c.Columns) > 0 {
		b.WriteString(" (")
		b.WriteString(strings.Join(c.Columns, ", "))
		b.WriteString(")")
	}

	b.WriteString(" AS (")
	b.WriteString(c.SQL)
	b.WriteString(")")

	return b.String(), c.Args
}

// String للتصحيح.
func (c *CTE) String() string {
	sql, _ := c.ToSQL()
	return sql
}

// ═══════════════════════════════════════════════
// CTEBuilder
// ═══════════════════════════════════════════════

// CTEBuilder يمثل استعلامًا مع CTEs.
type CTEBuilder[T any] struct {
	ctes      []*CTE
	mainQuery *gormx.QuerySet[T]
}

// With يبدأ استعلامًا مع CTE.
func With[T any](ctes ...*CTE) *CTEBuilder[T] {
	return &CTEBuilder[T]{ctes: ctes}
}

// Query يحدد الاستعلام الرئيسي.
func (cb *CTEBuilder[T]) Query(q *gormx.QuerySet[T]) *CTEBuilder[T] {
	cb.mainQuery = q
	return cb
}

// All ينفّذ الاستعلام.
func (cb *CTEBuilder[T]) All() ([]T, error) {
	if cb.mainQuery == nil {
		return nil, fmt.Errorf("gormx/advanced: CTEBuilder requires Query")
	}

	cteSQL, cteArgs := cb.buildCTEs()
	mainSQL, mainArgs := cb.mainQuery.ToSQL()

	fullSQL := cteSQL + " " + mainSQL
	allArgs := append(cteArgs, mainArgs...)

	var zero T
	var results []T
	err := gormx.DB().Model(&zero).Raw(fullSQL, allArgs...).Scan(&results).Error
	return results, err
}

// ScanInto ينفّذ ويقرأ في struct.
func (cb *CTEBuilder[T]) ScanInto(dest any) error {
	if cb.mainQuery == nil {
		return fmt.Errorf("gormx/advanced: CTEBuilder requires Query")
	}
	if dest == nil {
		return fmt.Errorf("gormx/advanced: nil destination")
	}

	cteSQL, cteArgs := cb.buildCTEs()
	mainSQL, mainArgs := cb.mainQuery.ToSQL()

	fullSQL := cteSQL + " " + mainSQL
	allArgs := append(cteArgs, mainArgs...)

	return gormx.DB().Raw(fullSQL, allArgs...).Scan(dest).Error
}

// buildCTEs يبني سلسلة CTEs.
func (cb *CTEBuilder[T]) buildCTEs() (string, []any) {
	if len(cb.ctes) == 0 {
		return "", nil
	}

	var parts []string
	var allArgs []any

	for i, cte := range cb.ctes {
		var part string
		if i == 0 {
			part, _ = cte.ToSQL()
		} else {
			part = cte.Name + " AS (" + cte.SQL + ")"
		}
		parts = append(parts, part)
		allArgs = append(allArgs, cte.Args...)
	}

	return strings.Join(parts, ", "), allArgs
}

// ═══════════════════════════════════════════════
// RecursiveTree
// ═══════════════════════════════════════════════

// RecursiveTree يبني CTE عودي للشجرة.
func RecursiveTree[T any](
	cteName string,
	idCol, parentCol string,
	allCols []string,
	recursiveQuery func(cteName string) string,
	rootID any,
) *CTE {
	cols := strings.Join(allCols, ", ")

	baseSQL := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = ?",
		cols, gormx.TableNameOf[T](), idCol,
	)

	recursive := recursiveQuery(cteName)

	return &CTE{
		Name:      cteName,
		SQL:       baseSQL + " UNION ALL " + recursive,
		Args:      []any{rootID},
		Recursive: true,
		Columns:   allCols,
	}
}