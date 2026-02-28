package errors

import (
	"testing"

	contractscodes "chainwise/contracts/codes"
	"google.golang.org/grpc/codes"
)

func TestDomainCodeToGRPC(t *testing.T) {
	tests := []struct {
		domainCode string
		want       codes.Code
	}{
		{contractscodes.OrderNotFound, codes.NotFound},
		{contractscodes.NotFound, codes.NotFound},
		{contractscodes.InsufficientStock, codes.FailedPrecondition},
		{contractscodes.ReservationFailed, codes.FailedPrecondition},
		{contractscodes.InvalidOrderStatusTransition, codes.FailedPrecondition},
		{contractscodes.Conflict, codes.FailedPrecondition},
		{contractscodes.Unauthenticated, codes.Unauthenticated},
		{contractscodes.InvalidCredentials, codes.Unauthenticated},
		{contractscodes.ValidationError, codes.InvalidArgument},
		{contractscodes.Internal, codes.Internal},
		{"UNKNOWN_CODE", codes.Unknown},
	}
	for _, tt := range tests {
		got := DomainCodeToGRPC(tt.domainCode)
		if got != tt.want {
			t.Errorf("DomainCodeToGRPC(%q) = %v, want %v", tt.domainCode, got, tt.want)
		}
	}
}
