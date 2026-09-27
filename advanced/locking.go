package advanced

import (
	"context"

	"github.com/sanad/gormx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ═══════════════════════════════════════════════
// Locking — Pessimistic & Optimistic
// ═══════════════════════════════════════════════

// LockMode وضع القفل.
type LockMode string

const (
	LockUpdate      LockMode = "UPDATE"
	LockNoKeyUpdate LockMode = "NO KEY UPDATE"
	LockShare       LockMode = "SHARE"
	LockKeyShare    LockMode = "KEY SHARE"
)

// ═══════════════════════════════════════════════
// PessimisticQuery
// ═══════════════════════════════════════════════

// PessimisticQuery يمثل استعلامًا مع قفل.
type PessimisticQuery[T any] struct {
	q    *gormx.QuerySet[T]
	mode LockMode
}

// WithLock ينشئ PessimisticQuery.
func WithLock[T any](q *gormx.QuerySet[T]) *PessimisticQuery[T] {
	return &PessimisticQuery[T]{q: q, mode: LockUpdate}
}

// ForUpdate يفعّل FOR UPDATE.
func (pq *PessimisticQuery[T]) ForUpdate() *PessimisticQuery[T] {
	pq.mode = LockUpdate
	return pq
}

// ForShare يفعّل FOR SHARE.
func (pq *PessimisticQuery[T]) ForShare() *PessimisticQuery[T] {
	pq.mode = LockShare
	return pq
}

// ForNoKeyUpdate يفعّل FOR NO KEY UPDATE.
func (pq *PessimisticQuery[T]) ForNoKeyUpdate() *PessimisticQuery[T] {
	pq.mode = LockNoKeyUpdate
	return pq
}

// ForKeyShare يفعّل FOR KEY SHARE.
func (pq *PessimisticQuery[T]) ForKeyShare() *PessimisticQuery[T] {
	pq.mode = LockKeyShare
	return pq
}

// Filter يضيف فلتر.
func (pq *PessimisticQuery[T]) Filter(field string, value any) *PessimisticQuery[T] {
	pq.q = pq.q.Filter(field, value)
	return pq
}

// Where يضيف SQL خام.
func (pq *PessimisticQuery[T]) Where(sql string, args ...any) *PessimisticQuery[T] {
	pq.q = pq.q.Where(sql, args...)
	return pq
}

// All ينفّذ.
func (pq *PessimisticQuery[T]) All() ([]T, error) {
	var results []T
	err := pq.build().Find(&results).Error
	return results, err
}

// First يرجّع أول سجل.
func (pq *PessimisticQuery[T]) First() (*T, error) {
	var result T
	err := pq.build().First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get يرجّع سجلًا بالـ ID.
func (pq *PessimisticQuery[T]) Get(id any) (*T, error) {
	var result T
	err := pq.build().Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// build يبني الاستعلام.
func (pq *PessimisticQuery[T]) build() *gorm.DB {
	db := pq.q.Build()
	strength := string(pq.mode)
	if strength == "" {
		strength = string(LockUpdate)
	}
	return db.Clauses(clause.Locking{Strength: strength})
}

// ═══════════════════════════════════════════════
// Optimistic Locking
// ═══════════════════════════════════════════════

// DefaultOptimisticField اسم الحقل الافتراضي.
const DefaultOptimisticField = "version"

// ═══════════════════════════════════════════════
// WithLockedTransaction
// ═══════════════════════════════════════════════

// WithLockedTransaction ينفّذ عملية داخل transaction مع قفل.
func WithLockedTransaction(
	ctx context.Context,
	cfg TxConfig,
	fn func(tx *Tx) error,
) error {
	return WithTransaction(ctx, cfg, fn)
}

// ═══════════════════════════════════════════════
// Manual Lock Helpers
// ═══════════════════════════════════════════════

// LockClause ينشئ clause.Locking.
func LockClause(mode LockMode) clause.Locking {
	if mode == "" {
		mode = LockUpdate
	}
	return clause.Locking{Strength: string(mode)}
}

// LockClauseWithTable ينشئ clause.Locking مع جدول محدد.
func LockClauseWithTable(mode LockMode, table string) clause.Locking {
	l := LockClause(mode)
	l.Table = clause.Table{Name: table}
	return l
}