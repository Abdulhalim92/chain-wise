package interceptors

import (
	"context"
	"log/slog"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryRequestID_PassesThrough(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}
	resp, err := UnaryRequestID(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)
	if err != nil {
		t.Fatal(err)
	}
	if resp != "ok" {
		t.Errorf("resp = %v, want ok", resp)
	}
}

func TestUnaryRecovery_ReturnsInternalOnPanic(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		panic("test panic")
	}
	log := slog.Default()
	interceptor := UnaryRecovery(log)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)
	if err == nil {
		t.Fatal("expected error on panic")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Internal {
		t.Errorf("want Internal, got %v", err)
	}
}
