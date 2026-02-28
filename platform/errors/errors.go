// Package errors — кастомные типы ошибок и конструкторы для маппинга в транспорт (gRPC/REST).
package errors

import (
	"errors"
	"fmt"
)

// DomainError — ошибка с доменным кодом (для маппинга в envelope.error.code и gRPC details).
type DomainError struct {
	code    string
	Message string
}

func (e *DomainError) Error() string { return e.Message }

// Code возвращает доменный код для маппинга в транспорт.
func (e *DomainError) Code() string { return e.code }

// NotFound — сущность не найдена.
type NotFound struct{ DomainError }

func NewNotFound(code, message string) *NotFound {
	return &NotFound{DomainError{code: code, Message: message}}
}

// Conflict — конфликт (например, дубликат, неверное состояние).
type Conflict struct{ DomainError }

func NewConflict(code, message string) *Conflict {
	return &Conflict{DomainError{code: code, Message: message}}
}

// Unauthenticated — не авторизован.
type Unauthenticated struct{ DomainError }

func NewUnauthenticated(code, message string) *Unauthenticated {
	return &Unauthenticated{DomainError{code: code, Message: message}}
}

// Detail — одна запись ошибки валидации по полю.
type Detail struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Validation — ошибка валидации (error.details в envelope).
type Validation struct {
	DomainError
	Details []Detail
}

func NewValidation(code, message string, details []Detail) *Validation {
	return &Validation{DomainError: DomainError{code: code, Message: message}, Details: details}
}

func (e *Validation) Error() string {
	if len(e.Details) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s: %d field(s)", e.Message, len(e.Details))
}

// Code возвращает доменный код (для Validation — общий код, напр. VALIDATION_ERROR).
func (e *Validation) Code() string { return e.DomainError.Code() }

// coder — интерфейс для извлечения доменного кода (DomainError, NotFound, Conflict, Unauthenticated, Validation).
type coder interface{ Code() string }

// CodeFromError извлекает доменный код из err; ok == false, если err не содержит кода.
func CodeFromError(err error) (code string, ok bool) {
	var c coder
	if errors.As(err, &c) {
		return c.Code(), true
	}
	return "", false
}
