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
	subs      []chan Chunk
	closeFunc CloseFunc

	// new fields for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

type Subscription struct {
	ch     chan Chunk
	cancel func()
}

func (s *Subscription) Cancel() {
	s.cancel()
}

func (s *Subscription) Read() <-chan Chunk {
	return s.ch
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
	ch := make(chan Chunk, buffer)
	s.mu.Lock()
	ch <- s.cache
	s.subs = append(s.subs, ch)
	s.mu.Unlock()
	return &Subscription{ch, func() { s.unsubscribe(ch) }}
}

// unsubscribe removes ch from s.subs and closes it.
func (s *Stream) unsubscribe(ch chan Chunk) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// find & remove
	for i, c := range s.subs {
		if c == ch {
			s.subs = slices.Delete(s.subs, i, i+1)
			close(ch)
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
	for _, ch := range subs {
		select {
		case ch <- chunk:
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
		s.done = true
		close(s.pub)
		for _, ch := range s.subs {
			close(ch)
		}
		s.subs = nil
		s.closeFunc(s.cache, s.err)
	})
}
