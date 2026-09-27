package gormx_test

import (
	"testing"

	"github.com/sanad/gormx"
	"github.com/stretchr/testify/assert"
)

func TestQ_SimpleOr(t *testing.T) {
	q := gormx.QOr(
		gormx.Eq("status", "active"),
		gormx.Eq("status", "pending"),
	)

	sql, args := q.ToSQL()
	assert.Contains(t, sql, "status = ?")
	assert.Contains(t, sql, " OR ")
	assert.Len(t, args, 2)
}

func TestQ_SimpleAnd(t *testing.T) {
	q := gormx.QAnd(
		gormx.Eq("status", "active"),
		gormx.Gt("age", 18),
	)

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, " AND ")
}

func TestQ_Nested(t *testing.T) {
	q := gormx.Qb().And(
		gormx.QOr(
			gormx.Eq("status", "active"),
			gormx.Eq("status", "pending"),
		),
		gormx.Gt("age", 18),
	)

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, "(")
	assert.Contains(t, sql, ")")
	assert.Contains(t, sql, " AND ")
	assert.Contains(t, sql, " OR ")
}

func TestQ_Not(t *testing.T) {
	q := gormx.QOr(
		gormx.Not(gormx.Eq("status", "deleted")),
		gormx.Eq("status", "active"),
	)

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, "NOT (")
}

func TestQ_NotWithNestedQ(t *testing.T) {
	inner := gormx.QOr(
		gormx.Eq("a", 1),
		gormx.Eq("b", 2),
	)

	q := gormx.Qb().And(
		gormx.Not(inner),
		gormx.Eq("c", 3),
	)

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, "NOT (")
}

func TestQ_Immutability(t *testing.T) {
	base := gormx.Qb().And(gormx.Eq("a", 1))
	extended := base.And(gormx.Eq("b", 2))

	sqlBase, _ := base.ToSQL()
	sqlExt, _ := extended.ToSQL()

	assert.NotEqual(t, sqlBase, sqlExt)
	assert.NotContains(t, sqlBase, "b = ?")
}

func TestQ_OrAfterAnd(t *testing.T) {
	q := gormx.Qb().
		And(gormx.Eq("a", 1)).
		Or(gormx.Eq("b", 2))

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, " OR ")
}

func TestQ_Empty(t *testing.T) {
	q := gormx.Qb()
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Len())

	sql, _ := q.ToSQL()
	assert.Equal(t, "1=1", sql)
}

func TestQ_AndGroup(t *testing.T) {
	q := gormx.Qb().
		And(gormx.Eq("x", 1)).
		AndGroup(
			gormx.Eq("a", 2),
			gormx.Eq("b", 3),
		)

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, " AND ")
	assert.Contains(t, sql, "(a = ? AND b = ?)")
}

func TestQ_OrGroup(t *testing.T) {
	q := gormx.Qb().
		And(gormx.Eq("x", 1)).
		OrGroup(
			gormx.Eq("a", 2),
			gormx.Eq("b", 3),
		)

	sql, _ := q.ToSQL()
	assert.Contains(t, sql, "(a = ? OR b = ?)")
}