package advanced

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
// Savepoints — نقاط الحفظ
// ═══════════════════════════════════════════════

// Savepoint يمثل نقطة حفظ.
type Savepoint struct {
	Name string
	tx   *gorm.DB
}

// NewSavepoint ينشئ savepoint.
func NewSavepoint(name string, tx *gorm.DB) *Savepoint {
	if name == "" {
		name = fmt.Sprintf("sp_%d", time.Now().UnixNano())
	}
	return &Savepoint{Name: name, tx: tx}
}

// Create ينشئ savepoint.
func (sp *Savepoint) Create() error {
	if sp.tx == nil {
		return errors.New("gormx/advanced: nil transaction")
	}
	return sp.tx.SavePoint(sp.Name).Error
}

// Rollback يتراجع إلى savepoint.
func (sp *Savepoint) Rollback() error {
	if sp.tx == nil {
		return errors.New("gormx/advanced: nil transaction")
	}
	return sp.tx.RollbackTo(sp.Name).Error
}

// Release يحذف savepoint.
func (sp *Savepoint) Release() error {
	if sp.tx == nil {
		return errors.New("gormx/advanced: nil transaction")
	}
	return sp.tx.Exec("RELEASE SAVEPOINT " + sp.Name).Error
}

// ═══════════════════════════════════════════════
// SavepointStack
// ═══════════════════════════════════════════════

// SavepointStack يدير مجموعة savepoints.
type SavepointStack struct {
	tx    *gorm.DB
	stack []*Savepoint
}

// NewSavepointStack ينشئ stack جديد.
func NewSavepointStack(tx *gorm.DB) *SavepointStack {
	return &SavepointStack{tx: tx}
}

// Push يضيف savepoint.
func (s *SavepointStack) Push(name string) (*Savepoint, error) {
	sp := NewSavepoint(name, s.tx)
	if err := sp.Create(); err != nil {
		return nil, err
	}
	s.stack = append(s.stack, sp)
	return sp, nil
}

// Pop يزيل ويحرّر آخر savepoint.
func (s *SavepointStack) Pop() error {
	if len(s.stack) == 0 {
		return errors.New("gormx/advanced: empty stack")
	}
	sp := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return sp.Release()
}

// RollbackLast يتراجع إلى آخر savepoint.
func (s *SavepointStack) RollbackLast() error {
	if len(s.stack) == 0 {
		return errors.New("gormx/advanced: empty stack")
	}
	sp := s.stack[len(s.stack)-1]
	return sp.Rollback()
}

// RollbackTo يتراجع إلى savepoint معين.
func (s *SavepointStack) RollbackTo(name string) error {
	for i := len(s.stack) - 1; i >= 0; i-- {
		if s.stack[i].Name == name {
			return s.stack[i].Rollback()
		}
	}
	return fmt.Errorf("gormx/advanced: savepoint %q not found", name)
}

// Depth يرجّع عدد savepoints.
func (s *SavepointStack) Depth() int {
	return len(s.stack)
}