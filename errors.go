package gormx

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
// Errors — الأخطاء الأساسية
// ═══════════════════════════════════════════════

var (
	// ErrNotFound يُرجع عند البحث عن سجل غير موجود.
	ErrNotFound = gorm.ErrRecordNotFound

	// ErrNotInitialized يُرجع لو SetDB لم يُنادى.
	ErrNotInitialized = errors.New("gormx: DB not initialized")

	// ErrNilDB يُرجع لو تم تمرير nil.
	ErrNilDB = errors.New("gormx: nil DB")

	// ErrInvalidField يُرجع عند حقل غير صحيح.
	ErrInvalidField = errors.New("gormx: invalid field")

	// ErrInvalidQuery يُرجع عند استعلام غير صحيح.
	ErrInvalidQuery = errors.New("gormx: invalid query")

	// ErrDangerousOperation يُرجع عند عملية خطرة بدون conditions.
	ErrDangerousOperation = errors.New("gormx: dangerous operation without conditions")

	// ErrAlreadyRegistered يُرجع عند تسجيل موديل بنفس الاسم مرتين.
	ErrAlreadyRegistered = errors.New("gormx: model already registered")

	// ErrNotFoundInRegistry يُرجع عند البحث عن موديل غير مسجّل.
	ErrNotFoundInRegistry = errors.New("gormx: model not found in registry")
)

// ═══════════════════════════════════════════════
// Typed Errors — أخطاء مع سياق
// ═══════════════════════════════════════════════

// NotFoundError يحمل معلومات عن السجل غير الموجود.
type NotFoundError struct {
	Model string
	ID    any
}

func (e *NotFoundError) Error() string {
	if e.ID == nil {
		return fmt.Sprintf("gormx: %s not found", e.Model)
	}
	return fmt.Sprintf("gormx: %s with id %v not found", e.Model, e.ID)
}

// Unwrap يسمح باستخدام errors.Is(err, ErrNotFound).
func (e *NotFoundError) Unwrap() error {
	return ErrNotFound
}

// NewNotFoundError ينشئ NotFoundError.
func NewNotFoundError(model string, id any) *NotFoundError {
	return &NotFoundError{Model: model, ID: id}
}

// ValidationError يحمل معلومات عن حقل غير صحيح.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("gormx: invalid field %q: %s", e.Field, e.Reason)
}

func (e *ValidationError) Unwrap() error {
	return ErrInvalidField
}

// NewValidationError ينشئ ValidationError.
func NewValidationError(field, reason string) *ValidationError {
	return &ValidationError{Field: field, Reason: reason}
}

// DangerousOperationError يحمل معلومات عن عملية خطرة.
type DangerousOperationError struct {
	Operation string
	Reason    string
}

func (e *DangerousOperationError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("gormx: dangerous operation %q", e.Operation)
	}
	return fmt.Sprintf("gormx: dangerous operation %q: %s", e.Operation, e.Reason)
}

func (e *DangerousOperationError) Unwrap() error {
	return ErrDangerousOperation
}

// NewDangerousError ينشئ DangerousOperationError.
func NewDangerousError(operation, reason string) *DangerousOperationError {
	return &DangerousOperationError{Operation: operation, Reason: reason}
}

// ═══════════════════════════════════════════════
// Error Checks — أدوات فحص
// ═══════════════════════════════════════════════

// IsNotFound يفحص إذا كان الخطأ ErrNotFound.
//
// يدعم NotFoundError أيضًا.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsValidation يفحص إذا كان الخطأ ValidationError.
func IsValidation(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve) || errors.Is(err, ErrInvalidField)
}

// IsDangerous يفحص إذا كان الخطأ DangerousOperationError.
func IsDangerous(err error) bool {
	var de *DangerousOperationError
	return errors.As(err, &de) || errors.Is(err, ErrDangerousOperation)
}

// IsAlreadyRegistered يفحص إذا كان الخطأ ErrAlreadyRegistered.
func IsAlreadyRegistered(err error) bool {
	return errors.Is(err, ErrAlreadyRegistered)
}

// IsNotFoundInRegistry يفحص إذا كان الخطأ ErrNotFoundInRegistry.
func IsNotFoundInRegistry(err error) bool {
	return errors.Is(err, ErrNotFoundInRegistry)
}

// AsNotFound يحوّل الخطأ إلى NotFoundError إذا أمكن.
func AsNotFound(err error) (*NotFoundError, bool) {
	var ne *NotFoundError
	if errors.As(err, &ne) {
		return ne, true
	}
	return nil, false
}

// AsValidation يحوّل الخطأ إلى ValidationError إذا أمكن.
func AsValidation(err error) (*ValidationError, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve, true
	}
	return nil, false
}

// AsDangerous يحوّل الخطأ إلى DangerousOperationError إذا أمكن.
func AsDangerous(err error) (*DangerousOperationError, bool) {
	var de *DangerousOperationError
	if errors.As(err, &de) {
		return de, true
	}
	return nil, false
}