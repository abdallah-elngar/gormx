// Package advanced provides advanced query capabilities for gormx.
package advanced

import (
	"context"
	"fmt"

	"github.com/sanad/gormx"
	"github.com/sanad/gormx/internal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ═══════════════════════════════════════════════
// Bulk Operations — عمليات جماعية
// ═══════════════════════════════════════════════

const (
	// DefaultBatchSize حجم الدفعة الافتراضي.
	DefaultBatchSize = 1000

	// MaxBatchSize أقصى حجم للدفعة.
	MaxBatchSize = 10000
)

// BulkConfig إعدادات العمليات الجماعية.
type BulkConfig struct {
	BatchSize        int
	UpdateOnConflict bool
	ConflictColumns  []string
	UpdateColumns    []string
	IgnoreOnConflict bool
}

// DefaultBulkConfig إعدادات افتراضية.
func DefaultBulkConfig() BulkConfig {
	return BulkConfig{
		BatchSize:        DefaultBatchSize,
		UpdateOnConflict: false,
		IgnoreOnConflict: false,
	}
}

// normalizeBatchSize يضبط حجم الدفعة ضمن الحدود.
func normalizeBatchSize(size int) int {
	if size <= 0 {
		return DefaultBatchSize
	}
	if size > MaxBatchSize {
		return MaxBatchSize
	}
	return size
}

// ═══════════════════════════════════════════════
// BulkInsert
// ═══════════════════════════════════════════════

// BulkInsert ينشئ سجلات كثيرة بكفاءة.
//
//	users := make([]User, 100000)
//	err := advanced.BulkInsert[User](ctx, users, advanced.DefaultBulkConfig())
func BulkInsert[T any](ctx context.Context, items []T, cfg BulkConfig) error {
	if len(items) == 0 {
		return nil
	}

	batchSize := normalizeBatchSize(cfg.BatchSize)
	db := gormx.DB().WithContext(ctx)

	if err := db.CreateInBatches(&items, batchSize).Error; err != nil {
		return fmt.Errorf("gormx/advanced: bulk insert failed: %w", err)
	}

	return nil
}

// ═══════════════════════════════════════════════
// BulkUpsert
// ═══════════════════════════════════════════════

// BulkUpsert يدخل أو يحدّث سجلات.
//
//	err := advanced.BulkUpsert[User](ctx, users, advanced.BulkConfig{
//	    ConflictColumns: []string{"email"},
//	    UpdateColumns:   []string{"name", "age"},
//	})
func BulkUpsert[T any](ctx context.Context, items []T, cfg BulkConfig) error {
	if len(items) == 0 {
		return nil
	}
	if len(cfg.ConflictColumns) == 0 {
		return fmt.Errorf("gormx/advanced: BulkUpsert requires ConflictColumns")
	}

	// Validate ConflictColumns
	for _, col := range cfg.ConflictColumns {
		if err := internal.ValidateField(col); err != nil {
			return err
		}
	}
	for _, col := range cfg.UpdateColumns {
		if err := internal.ValidateField(col); err != nil {
			return err
		}
	}

	batchSize := normalizeBatchSize(cfg.BatchSize)
	db := gormx.DB().WithContext(ctx)

	onConflict := clause.OnConflict{
		Columns: make([]clause.Column, len(cfg.ConflictColumns)),
	}
	for i, col := range cfg.ConflictColumns {
		onConflict.Columns[i] = clause.Column{Name: col}
	}

	if cfg.IgnoreOnConflict {
		onConflict.DoNothing = true
	} else if len(cfg.UpdateColumns) > 0 {
		onConflict.DoUpdates = clause.AssignmentColumns(cfg.UpdateColumns)
	} else {
		onConflict.UpdateAll = true
	}

	if err := db.Clauses(onConflict).CreateInBatches(&items, batchSize).Error; err != nil {
		return fmt.Errorf("gormx/advanced: bulk upsert failed: %w", err)
	}

	return nil
}

// ═══════════════════════════════════════════════
// BulkUpdate
// ═══════════════════════════════════════════════

// UpdateItem يمثل تحديثًا لسجل واحد.
type UpdateItem[T any] struct {
	ID     any
	Values map[string]any
}

// BulkUpdate يحدّث عدة سجلات بقيم مختلفة.
//
//	updates := []advanced.UpdateItem[User]{
//	    {ID: 1, Values: map[string]any{"age": 31}},
//	    {ID: 2, Values: map[string]any{"age": 26}},
//	}
//	err := advanced.BulkUpdate[User](ctx, "id", updates)
func BulkUpdate[T any](ctx context.Context, idColumn string, items []UpdateItem[T]) error {
	if len(items) == 0 {
		return nil
	}
	if idColumn == "" {
		idColumn = "id"
	}
	if err := internal.ValidateField(idColumn); err != nil {
		return err
	}

	db := gormx.DB().WithContext(ctx)

	type columnData struct {
		ids    []any
		values []any
	}
	columnMap := make(map[string]*columnData)
	allIDs := make([]any, 0, len(items))

	for _, item := range items {
		allIDs = append(allIDs, item.ID)
		for col, val := range item.Values {
			if err := internal.ValidateField(col); err != nil {
				return err
			}
			if columnMap[col] == nil {
				columnMap[col] = &columnData{}
			}
			columnMap[col].ids = append(columnMap[col].ids, item.ID)
			columnMap[col].values = append(columnMap[col].values, val)
		}
	}

	assignments := make(clause.Set, 0, len(columnMap))
	for col, cd := range columnMap {
		sql, args := buildCaseSQL(idColumn, col, cd.ids, cd.values)
		assignments = append(assignments, clause.Assignment{
			Column: clause.Column{Name: col},
			Value:  gorm.Expr(sql, args...),
		})
	}

	var zero T
	return db.Model(&zero).
		Where(idColumn+" IN ?", allIDs).
		Updates(assignments).Error
}

// buildCaseSQL يبني CASE WHEN ... THEN ... ELSE column END.
func buildCaseSQL(idColumn, column string, ids, values []any) (string, []any) {
	var sql string
	sql = "CASE"
	args := make([]any, 0, len(ids)*2)

	for i := range ids {
		sql += " WHEN " + idColumn + " = ? THEN ?"
		args = append(args, ids[i], values[i])
	}

	sql += " ELSE " + column + " END"
	return sql, args
}

// ═══════════════════════════════════════════════
// BulkDelete
// ═══════════════════════════════════════════════

// BulkDeleteByIDs يحذف سجلات كثيرة بالـ IDs.
//
//	deleted, err := advanced.BulkDeleteByIDs[User](ctx, ids, 1000)
func BulkDeleteByIDs[T any](ctx context.Context, ids []any, batchSize int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	batchSize = normalizeBatchSize(batchSize)

	db := gormx.DB().WithContext(ctx)
	var zero T
	var total int64

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		res := db.Delete(&zero, batch)
		if res.Error != nil {
			return total, res.Error
		}
		total += res.RowsAffected
	}

	return total, nil
}

// BulkDeleteWhere يحذف حسب شرط.
//
// ⚠️ يتطلب condition — لحماية من الحذف الكامل.
func BulkDeleteWhere[T any](ctx context.Context, where string, args ...any) (int64, error) {
	if where == "" {
		return 0, fmt.Errorf("gormx/advanced: BulkDeleteWhere requires a WHERE clause")
	}

	db := gormx.DB().WithContext(ctx)
	var zero T
	res := db.Where(where, args...).Delete(&zero)
	return res.RowsAffected, res.Error
}

// BulkCount يحسب عدد السجلات حسب شرط.
func BulkCount[T any](ctx context.Context, where string, args ...any) (int64, error) {
	db := gormx.DB().WithContext(ctx)
	var zero T
	var count int64

	q := db.Model(&zero)
	if where != "" {
		q = q.Where(where, args...)
	}

	err := q.Count(&count).Error
	return count, err
}