package middleware

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkCircuitBreaker_Execute_Closed(b *testing.B) {
	logger, _ := zap.NewProduction()
	cb := NewCircuitBreaker("bench-test", DefaultCircuitBreakerConfig(), logger)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cb.Execute(ctx, func(ctx context.Context) error {
			return nil
		})
	}
}

func BenchmarkCircuitBreaker_Execute_Parallel(b *testing.B) {
	logger, _ := zap.NewProduction()
	cb := NewCircuitBreaker("bench-parallel", DefaultCircuitBreakerConfig(), logger)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Execute(ctx, func(ctx context.Context) error {
				return nil
			})
		}
	})
}

func BenchmarkCircuitBreaker_CanExecute(b *testing.B) {
	logger, _ := zap.NewProduction()
	cb := NewCircuitBreaker("bench-can-exec", DefaultCircuitBreakerConfig(), logger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cb.canExecute()
	}
}

// BenchmarkCircuitBreaker_Execute_WithFailures benchmarks execution when
// some calls fail, triggering state transitions
func BenchmarkCircuitBreaker_Execute_WithFailures(b *testing.B) {
	logger := zap.NewNop()
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 1000 // high threshold to avoid opening during bench
	cb := NewCircuitBreaker("bench-failures", cfg, logger)
	ctx := context.Background()
	errFail := errors.New("simulated failure")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
			cb.Execute(ctx, func(ctx context.Context) error {
				return errFail
			})
		} else {
			cb.Execute(ctx, func(ctx context.Context) error {
				return nil
			})
		}
	}
}

// BenchmarkCircuitBreaker_RecordResult benchmarks result recording in isolation
func BenchmarkCircuitBreaker_RecordResult(b *testing.B) {
	logger := zap.NewNop()
	cb := NewCircuitBreaker("bench-record", DefaultCircuitBreakerConfig(), logger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cb.recordResult(nil)
	}
}
