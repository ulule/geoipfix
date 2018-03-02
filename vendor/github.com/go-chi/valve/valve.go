package valve

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

var (
	ValveCtxKey     = &contextKey{"Valve"}
	ErrTimedout     = errors.New("valve: shutdown timed out")
	ErrShuttingdown = errors.New("valve: shutdown in progress")
	ErrOff          = errors.New("valve: valve already shutdown")
)

type Valve struct {
	stopCh chan struct{}
	wg     sync.WaitGroup

	shutdown bool
	mu       sync.Mutex
}

type LeverControl interface {
	Stop() <-chan struct{}
	Add(delta int) error
	Done()
	Open() error
	Close()
}

func New() *Valve {
	return &Valve{
		stopCh: make(chan struct{}, 0),
	}
}

// Context returns a fresh context with the Lever value set.
//
// It is useful as the base context in a server, that provides shutdown
// signaling across a context tree.
func (v *Valve) Context() context.Context {
	return context.WithValue(context.Background(), ValveCtxKey, LeverControl(v))
}

// Lever returns the lever controls from a context object.
func Lever(ctx context.Context) LeverControl {
	valveCtx, ok := ctx.Value(ValveCtxKey).(LeverControl)
	if !ok {
		panic("valve: ValveCtxKey has not been set on the context.")
	}
	return valveCtx
}

// Shutdown will signal to the context to stop all processing, and will
// give a grace period of `timeout` duration. If `timeout` is 0 then it will
// wait indefinitely until all valves are closed.
func (v *Valve) Shutdown(timeout time.Duration) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.shutdown {
		return ErrOff
	}
	close(v.stopCh)
	v.shutdown = true

	if timeout == 0 {
		v.wg.Wait()
	} else {
		tc := make(chan struct{})
		go func() {
			defer close(tc)
			v.wg.Wait()
		}()
		select {
		case <-tc:
			return nil
		case <-time.After(timeout):
			return ErrTimedout
		}
	}

	return nil
}

// Stop returns a channel that will be closed once the system is supposed to
// be stopped. It mimics the behaviou of the ctx.Done() method in "context".
func (v *Valve) Stop() <-chan struct{} {
	return v.stopCh
}

// Add increments by `delta` (should be 1), to a waitgroup on the valve that
// signifies that a block of code must complete before we exit the system.
func (v *Valve) Add(delta int) error {
	select {
	case <-v.stopCh:
		return ErrShuttingdown
	default:
		v.wg.Add(delta)
		return nil
	}
}

// Done decrements the valve waitgroup that informs the lever control that the
// non-preemptive app code is finished.
func (v *Valve) Done() {
	v.wg.Done()
}

// Open is an alias for Add(1) intended to read better for opening a valve.
func (v *Valve) Open() error {
	return v.Add(1)
}

// Close is an alias for Done() intended to read better for closing a valve.
func (v *Valve) Close() {
	v.Done()
}

// ShutdownHandler is an optional HTTP middleware handler that will stop
// accepting new connections if the server is in a shutting-down state.
//
// If you're using something that github.com/tylerb/graceful which stops
// accepting new connections on the socket anyways, then this handler
// wouldnt be necessary, but it is handy otherwise.
func (v *Valve) ShutdownHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		lever := Lever(r.Context())
		lever.Open()
		defer lever.Close()

		select {
		case <-lever.Stop():
			// Shutdown in progress - don't accept new requests
			http.Error(w, ErrShuttingdown.Error(), http.StatusServiceUnavailable)

		default:
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "valve context value " + k.name
}
