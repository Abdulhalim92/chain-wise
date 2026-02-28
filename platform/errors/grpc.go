// Маппинг доменный код → gRPC code (источник кодов: contracts/codes).
// При добавлении кода в contracts/codes — добавить case в switch и строку в grpc_test.go.

package errors

import (
	contractscodes "chainwise/contracts/codes"
	"google.golang.org/grpc/codes"
)

// DomainCodeToGRPC возвращает gRPC-код по доменному коду (из contracts/codes).
func DomainCodeToGRPC(domainCode string) codes.Code {
	switch domainCode {
	case contractscodes.OrderNotFound, contractscodes.NotFound:
		return codes.NotFound
	case contractscodes.InsufficientStock, contractscodes.ReservationFailed, contractscodes.InvalidOrderStatusTransition, contractscodes.Conflict:
		return codes.FailedPrecondition
	case contractscodes.Unauthenticated, contractscodes.InvalidCredentials:
		return codes.Unauthenticated
	case contractscodes.ValidationError:
		return codes.InvalidArgument
	case contractscodes.Internal:
		return codes.Internal
	default:
		return codes.Unknown
	}
}
