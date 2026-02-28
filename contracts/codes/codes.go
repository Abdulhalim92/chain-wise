// Package codes — доменные коды ошибок (один источник правды для бэкенда и клиентов).
package codes

// Доменные коды ошибок.
const (
	OrderNotFound                 = "ORDER_NOT_FOUND"
	InsufficientStock             = "INSUFFICIENT_STOCK"
	ReservationFailed             = "RESERVATION_FAILED"
	InvalidOrderStatusTransition  = "INVALID_ORDER_STATUS_TRANSITION"
	Unauthenticated               = "UNAUTHENTICATED"
	InvalidCredentials           = "INVALID_CREDENTIALS"
	ValidationError              = "VALIDATION_ERROR"
	Conflict                      = "CONFLICT"
	NotFound                      = "NOT_FOUND"
	Internal                      = "INTERNAL"
)
