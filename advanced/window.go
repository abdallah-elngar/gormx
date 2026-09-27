package advanced

import (
	"strconv"
	"strings"
)

// ═══════════════════════════════════════════════
// Window Functions — دوال النافذة
// ═══════════════════════════════════════════════

// RowNumber → ROW_NUMBER() OVER (...)
func RowNumber(alias string, partition ...string) string {
	return buildWindow("ROW_NUMBER()", alias, "", partition, "")
}

// Rank → RANK() OVER (...)
func Rank(alias string, orderBy string, partition ...string) string {
	return buildWindow("RANK()", alias, orderBy, partition, "")
}

// DenseRank → DENSE_RANK() OVER (...)
func DenseRank(alias string, orderBy string, partition ...string) string {
	return buildWindow("DENSE_RANK()", alias, orderBy, partition, "")
}

// Lag → LAG(field, offset) OVER (...)
func Lag(field string, offset int, alias string, partition ...string) string {
	expr := "LAG(" + field
	if offset > 0 {
		expr += ", " + strconv.Itoa(offset)
	}
	expr += ")"
	return buildWindow(expr, alias, "", partition, "")
}

// Lead → LEAD(field, offset) OVER (...)
func Lead(field string, offset int, alias string, partition ...string) string {
	expr := "LEAD(" + field
	if offset > 0 {
		expr += ", " + strconv.Itoa(offset)
	}
	expr += ")"
	return buildWindow(expr, alias, "", partition, "")
}

// RunningSum → SUM(field) OVER (ORDER BY ...)
func RunningSum(field, alias, orderBy string, partition ...string) string {
	return buildWindow("SUM("+field+")", alias, orderBy, partition, "")
}

// RunningCount → COUNT(*) OVER (ORDER BY ...)
func RunningCount(alias, orderBy string, partition ...string) string {
	return buildWindow("COUNT(*)", alias, orderBy, partition, "")
}

// RunningAvg → AVG(field) OVER (ORDER BY ...)
func RunningAvg(field, alias, orderBy string, partition ...string) string {
	return buildWindow("AVG("+field+")", alias, orderBy, partition, "")
}

// AvgOver → AVG(field) OVER (...)
func AvgOver(field, alias string, partition ...string) string {
	return buildWindow("AVG("+field+")", alias, "", partition, "")
}

// SumOver → SUM(field) OVER (...)
func SumOver(field, alias string, partition ...string) string {
	return buildWindow("SUM("+field+")", alias, "", partition, "")
}

// CountOver → COUNT(*) OVER (...)
func CountOver(alias string, partition ...string) string {
	return buildWindow("COUNT(*)", alias, "", partition, "")
}

// NTile → NTILE(n) OVER (...)
func NTile(n int, alias, orderBy string, partition ...string) string {
	return buildWindow("NTILE("+strconv.Itoa(n)+")", alias, orderBy, partition, "")
}

// FirstValue → FIRST_VALUE(field) OVER (...)
func FirstValue(field, alias, orderBy string, partition ...string) string {
	return buildWindow("FIRST_VALUE("+field+")", alias, orderBy, partition, "")
}

// LastValue → LAST_VALUE(field) OVER (...)
func LastValue(field, alias, orderBy string, partition ...string) string {
	return buildWindow("LAST_VALUE("+field+")", alias, orderBy, partition, "")
}

// NthValue → NTH_VALUE(field, n) OVER (...)
func NthValue(field string, n int, alias, orderBy string, partition ...string) string {
	expr := "NTH_VALUE(" + field + ", " + strconv.Itoa(n) + ")"
	return buildWindow(expr, alias, orderBy, partition, "")
}

// PercentRank → PERCENT_RANK() OVER (...)
func PercentRank(alias, orderBy string, partition ...string) string {
	return buildWindow("PERCENT_RANK()", alias, orderBy, partition, "")
}

// CumeDist → CUME_DIST() OVER (...)
func CumeDist(alias, orderBy string, partition ...string) string {
	return buildWindow("CUME_DIST()", alias, orderBy, partition, "")
}

// ═══════════════════════════════════════════════
// Internal
// ═══════════════════════════════════════════════

func buildWindow(expr, alias, orderBy string, partition []string, frame string) string {
	var b strings.Builder
	b.WriteString(expr)
	b.WriteString(" OVER (")

	if len(partition) > 0 {
		b.WriteString("PARTITION BY ")
		b.WriteString(strings.Join(partition, ", "))
	}

	if orderBy != "" {
		if len(partition) > 0 {
			b.WriteString(" ")
		}
		b.WriteString("ORDER BY ")
		b.WriteString(orderBy)
	}

	if frame != "" {
		b.WriteString(" ")
		b.WriteString(frame)
	}

	b.WriteString(")")

	if alias != "" {
		b.WriteString(" AS ")
		b.WriteString(alias)
	}

	return b.String()
}

// WithFrame يضيف frame مخصص.
func WithFrame(expr, alias, orderBy string, partition []string, frame string) string {
	return buildWindow(expr, alias, orderBy, partition, frame)
}

// RowsBetween يبني frame clause.
func RowsBetween(start, end string) string {
	return "ROWS BETWEEN " + start + " AND " + end
}

// RangeBetween يبني RANGE frame clause.
func RangeBetween(start, end string) string {
	return "RANGE BETWEEN " + start + " AND " + end
}