package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildLookup_Equality(t *testing.T) {
	l, err := BuildLookup("name", "Ali", false)
	require.NoError(t, err)
	assert.Equal(t, "name = ?", l.SQL)
	assert.Equal(t, []any{"Ali"}, l.Args)
	assert.False(t, l.Negate)
}

func TestBuildLookup_GreaterThan(t *testing.T) {
	l, err := BuildLookup("age__gt", 18, false)
	require.NoError(t, err)
	assert.Equal(t, "age > ?", l.SQL)
	assert.Equal(t, []any{18}, l.Args)
}

func TestBuildLookup_Contains(t *testing.T) {
	l, err := BuildLookup("name__contains", "ali", false)
	require.NoError(t, err)
	assert.Equal(t, "name LIKE ?", l.SQL)
	assert.Equal(t, []any{"%ali%"}, l.Args)
}

func TestBuildLookup_IContains(t *testing.T) {
	l, err := BuildLookup("name__icontains", "ALI", false)
	require.NoError(t, err)
	assert.Contains(t, l.SQL, "LOWER(name)")
	assert.Contains(t, l.SQL, "LIKE")
	assert.Equal(t, []any{"%ali%"}, l.Args)
}

func TestBuildLookup_In(t *testing.T) {
	l, err := BuildLookup("id__in", []int{1, 2, 3}, false)
	require.NoError(t, err)
	assert.Equal(t, "id IN (?,?,?)", l.SQL)
	assert.Len(t, l.Args, 3)
}

func TestBuildLookup_InEmpty(t *testing.T) {
	l, err := BuildLookup("id__in", []int{}, false)
	require.NoError(t, err)
	assert.Equal(t, "1=0", l.SQL)
}

func TestBuildLookup_IsNull(t *testing.T) {
	l, err := BuildLookup("deleted_at__isnull", true, false)
	require.NoError(t, err)
	assert.Equal(t, "deleted_at IS NULL", l.SQL)

	l, err = BuildLookup("deleted_at__isnull", false, false)
	require.NoError(t, err)
	assert.Equal(t, "deleted_at IS NOT NULL", l.SQL)
}

func TestBuildLookup_Between(t *testing.T) {
	l, err := BuildLookup("age__between", []int{18, 65}, false)
	require.NoError(t, err)
	assert.Equal(t, "age BETWEEN ? AND ?", l.SQL)
	assert.Equal(t, []any{18, 65}, l.Args)
}

func TestBuildLookup_InvalidField(t *testing.T) {
	_, err := BuildLookup("bad;field", "x", false)
	assert.Error(t, err)
}

func TestBuildLookup_InvalidLookup(t *testing.T) {
	_, err := BuildLookup("name__bad", "x", false)
	assert.Error(t, err)
}

func TestBuildLookup_DateYear(t *testing.T) {
	l, err := BuildLookupWithDialect("created_at__year", 2025, false, DialectSQLite)
	require.NoError(t, err)
	assert.Contains(t, l.SQL, "strftime")
	assert.Equal(t, []any{2025}, l.Args)
}

func TestBuildLookupWithDialect_Postgres(t *testing.T) {
	l, err := BuildLookupWithDialect("created_at__year", 2025, false, DialectPostgres)
	require.NoError(t, err)
	assert.Contains(t, l.SQL, "EXTRACT")
	assert.Contains(t, l.SQL, "YEAR")
}

func TestBuildLookupWithDialect_MySQL(t *testing.T) {
	l, err := BuildLookupWithDialect("created_at__year", 2025, false, DialectMySQL)
	require.NoError(t, err)
	assert.Contains(t, l.SQL, "YEAR(")
}