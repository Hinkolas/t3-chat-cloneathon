package stream

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

type CloseFunc func(Chunk, error)

type Chunk struct {
	Reasoning string `json:"reasoning,omitempty"`
	Content   string `json:"content,omitempty"`
}

func (c *Chunk) append(c2 Chunk) {
	c.Reasoning += c2.Reasoning
	c.Content += c2.Content
}

// Stream represents one ongoing streaming process.
type Stream struct {
	mu        sync.RWMutex
	wg        sync.WaitGroup
	closeOnce sync.Once

	cache Chunk
	done  bool
	err   error

	pub       chan Chunk
	subs      []subscriber
	closeFunc CloseFunc

	// new fields for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

type subscriber struct {
	dataCh chan Chunk
	errCh  chan error
}

// Subscribe returns data- and error-channels.
type Subscription struct {
	Data   <-chan Chunk
	Err    <-chan error
	cancel func()
}

func (s *Subscription) Cancel() {
	s.cancel()
}

// New creates a Stream with a background context.
func New() *Stream {
	return NewWithContext(context.Background())
}

// NewWithContext creates a Stream using parent as its base context.
func NewWithContext(parent context.Context) *Stream {
	ctx, cancel := context.WithCancel(parent)
	s := &Stream{
		cache:     Chunk{},
		pub:       make(chan Chunk),
		closeFunc: func(Chunk, error) {},
		ctx:       ctx,
		cancel:    cancel,
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for chunk := range s.pub {
			s.emit(chunk)
		}
		s.finish()
	}()
	return s
}

// Context returns the Stream's context.
func (s *Stream) Context() context.Context {
	return s.ctx
}

func (s *Stream) OnClose(fn CloseFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeFunc = fn
}

// Publish sends a chunk unless the stream has been canceled.
func (s *Stream) Publish(c Chunk) {
	select {
	case <-s.ctx.Done():
		// stream was canceled, drop this chunk
		return
	case s.pub <- c:
	}
}

// Subscribe returns a channel on which the caller will receive all past and future chunks
func (s *Stream) Subscribe(buffer int) *Subscription {
	dataCh := make(chan Chunk, buffer)
	errCh := make(chan error, 1) // buffer=1 so we can send final err
	s.mu.Lock()
	dataCh <- s.cache
	s.subs = append(s.subs, subscriber{dataCh, errCh})
	s.mu.Unlock()

	return &Subscription{
		Data: dataCh,
		Err:  errCh,
		cancel: func() {
			s.unsubscribe(dataCh, errCh)
		},
	}
}

// unsubscribe removes ch from s.subs and closes it.
func (s *Stream) unsubscribe(dataCh chan Chunk, errCh chan error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, sb := range s.subs {
		if sb.dataCh == dataCh && sb.errCh == errCh {
			// drop from list
			s.subs = slices.Delete(s.subs, i, i+1)
			close(errCh)
			close(dataCh)
			return
		}
	}
}

// Wait blocks until the stream is done. It returns any error.
func (s *Stream) Wait() error {
	s.wg.Wait()
	return s.err
}

func (s *Stream) Fail(err error) {
	s.setError(err)
	s.cancel() // unblock any upstream readers
	s.finish()
}

// Closes the stream and all subscriber channels. Should be called ofter stream is done.
func (s *Stream) Close() {
	s.cancel()
	s.finish()
}

// emit is called for each incoming chunk.
// It appends to the buffer, fans out to all
// subscriber chans, and invokes handlers.
func (s *Stream) emit(chunk Chunk) {
	s.mu.Lock()
	// s.chunks = append(s.chunks, chunk)
	s.cache.append(chunk)        // accumulate into cache for Close
	subs := slices.Clone(s.subs) // snapshot subscribers
	s.mu.Unlock()

	// fan-out to subscribers (best-effort)
	for _, sub := range subs {
		select {
		case sub.dataCh <- chunk:
		default:
			// TODO: Handle subscriber not keeping up (e.g. cancel subscription)
			fmt.Println("subscriber is slow, blocking on chunk")
		}
	}
}

func (s *Stream) setError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err == nil {
		s.err = err
	}
}

// finish is called exactly once when the read loop terminates.
// It marks done and closes every subscriber channel.
func (s *Stream) finish() {
	s.closeOnce.Do(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		close(s.pub)
		// first send the terminal error (or nil) to each subscriber
		for _, sb := range s.subs {
			sb.errCh <- s.err
			close(sb.errCh)
			close(sb.dataCh)
		}
		s.subs = nil
		s.closeFunc(s.cache, s.err)
	})
}
