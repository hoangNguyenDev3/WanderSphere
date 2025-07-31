package middleware

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CircuitState represents the current state of the circuit breaker
type CircuitState int

const (
	StateClosed   CircuitState = iota // Normal operation, requests pass through
	StateOpen                         // Circuit is open, requests fail fast
	StateHalfOpen                     // Testing if service recovered
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig holds the configuration for a circuit breaker
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of consecutive failures before opening the circuit
	FailureThreshold int
	// SuccessThreshold is the number of consecutive successes in half-open to close the circuit
	SuccessThreshold int
	// OpenTimeout is how long to wait in open state before transitioning to half-open
	OpenTimeout time.Duration
	// HalfOpenMaxCalls is the max concurrent test calls allowed in half-open state
	HalfOpenMaxCalls int
}

// DefaultCircuitBreakerConfig returns sensible defaults
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenTimeout:      30 * time.Second,
		HalfOpenMaxCalls: 3,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu              sync.Mutex
	state           CircuitState
	failureCount    int
	successCount    int
	halfOpenCalls   int
	lastFailureTime time.Time
	config          CircuitBreakerConfig
	name            string
	logger          *zap.Logger
}

// NewCircuitBreaker creates a new circuit breaker with the given config
func NewCircuitBreaker(name string, config CircuitBreakerConfig, logger *zap.Logger) *CircuitBreaker {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return &CircuitBreaker{
		state:  StateClosed,
		config: config,
		name:   name,
		logger: logger,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := cb.canExecute(); err != nil {
		return err
	}

	err := fn(ctx)

	cb.recordResult(err)
	return err
}

// canExecute checks if a request is allowed through
func (cb *CircuitBreaker) canExecute() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		// Check if enough time has passed to transition to half-open
		if time.Since(cb.lastFailureTime) > cb.config.OpenTimeout {
			cb.transitionTo(StateHalfOpen)
			cb.halfOpenCalls++
			return nil
		}
		return status.Errorf(codes.Unavailable,
			"circuit breaker '%s' is OPEN; failing fast (last failure: %s ago)",
			cb.name, time.Since(cb.lastFailureTime).Round(time.Second))

	case StateHalfOpen:
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			return status.Errorf(codes.Unavailable,
				"circuit breaker '%s' is HALF_OPEN; max test calls (%d) reached",
				cb.name, cb.config.HalfOpenMaxCalls)
		}
		cb.halfOpenCalls++
		return nil
	}

	return nil
}

// recordResult records the result of a call and transitions state if needed
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil && isCircuitBreakerError(err) {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

// isCircuitBreakerError checks if an error should count toward the circuit breaker failure threshold
func isCircuitBreakerError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		// Non-gRPC error (e.g., context canceled, network error) — count as failure
		return true
	}
	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Aborted, codes.Internal:
		return true
	default:
		// Application-level errors (NotFound, InvalidArgument, etc.) are NOT circuit breaker failures
		return false
	}
}

func (cb *CircuitBreaker) onFailure() {
	switch cb.state {
	case StateClosed:
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.failureCount >= cb.config.FailureThreshold {
			cb.transitionTo(StateOpen)
		}

	case StateHalfOpen:
		// Any failure in half-open goes back to open
		cb.lastFailureTime = time.Now()
		cb.transitionTo(StateOpen)
	}
}

func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		// Reset failure count on success
		cb.failureCount = 0

	case StateHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.config.SuccessThreshold {
			cb.transitionTo(StateClosed)
		}
	}
}

func (cb *CircuitBreaker) transitionTo(newState CircuitState) {
	oldState := cb.state
	cb.state = newState
	cb.failureCount = 0
	cb.successCount = 0
	cb.halfOpenCalls = 0

	cb.logger.Info("Circuit breaker state transition",
		zap.String("name", cb.name),
		zap.String("from", oldState.String()),
		zap.String("to", newState.String()))
}

// State returns the current state (for monitoring/testing)
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// UnaryClientInterceptor creates a gRPC unary client interceptor that wraps calls with circuit breaker logic
func UnaryClientInterceptor(cb *CircuitBreaker) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return cb.Execute(ctx, func(ctx context.Context) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}
