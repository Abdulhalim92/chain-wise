package interceptors

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const requestIDMetaKey = "x-request-id"

type requestIDCtxKey struct{}

// UnaryRequestID внедряет или пробрасывает request_id в metadata и context.
func UnaryRequestID(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	id := getOrCreateRequestID(ctx)
	ctx = context.WithValue(ctx, requestIDCtxKey{}, id)
	if err := grpc.SetHeader(ctx, metadata.Pairs(requestIDMetaKey, id)); err != nil {
		// ignore
	}
	return handler(ctx, req)
}

// StreamRequestID то же для stream.
func StreamRequestID(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	id := getOrCreateRequestID(ctx)
	ctx = context.WithValue(ctx, requestIDCtxKey{}, id)
	if err := ss.SendHeader(metadata.Pairs(requestIDMetaKey, id)); err != nil {
		// ignore
	}
	return handler(srv, &streamWithContext{ServerStream: ss, ctx: ctx})
}

type streamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *streamWithContext) Context() context.Context { return s.ctx }

func getOrCreateRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v := md.Get(requestIDMetaKey); len(v) > 0 {
			return v[0]
		}
	}
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}

// UnaryRecovery ловит panic, логирует и возвращает gRPC error.
func UnaryRecovery(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "error", r, "stack", string(debug.Stack()), "method", info.FullMethod)
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(ctx, req)
	}
}

// StreamRecovery то же для stream.
func StreamRecovery(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "error", r, "stack", string(debug.Stack()), "method", info.FullMethod)
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(srv, ss)
	}
}

// UnaryLogging логирует вызов: method, status, duration.
func UnaryLogging(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}
		logger.Info("grpc call", "method", info.FullMethod, "status", code.String(), "duration_ms", time.Since(start).Milliseconds())
		return resp, err
	}
}

// StreamLogging то же для stream.
func StreamLogging(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}
		logger.Info("grpc stream", "method", info.FullMethod, "status", code.String(), "duration_ms", time.Since(start).Milliseconds())
		return err
	}
}
