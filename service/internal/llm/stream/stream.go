package stream

import (
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
	// protected by mu
	mu        sync.RWMutex
	wg        sync.WaitGroup
	closeOnce sync.Once

	cache Chunk // accumulation of all chunks received so far

	done bool  // true once the stream finishes
	err  error // any terminal error

	pub  chan Chunk   // where callers Publish
	subs []chan Chunk // subscriber channels

	closeFunc CloseFunc // callback for cleanup after finish
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

// New returns a Stream that’s ready to Start() / Publish().
func New() *Stream {
	s := &Stream{
		cache:     Chunk{},
		pub:       make(chan Chunk),
		closeFunc: func(Chunk, error) {},
	}
	// start a goroutine that fans-out anything sent on inCh
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

func (s *Stream) OnClose(closeFunc func(Chunk, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeFunc = closeFunc
}

// Publish sends one chunk into the stream. Blocks if no buffer on inCh.
func (s *Stream) Publish(c Chunk) {
	s.pub <- c
}

// Subscribe returns a channel on which the caller will receive all past and future chunks
func (s *Stream) Subscribe(buffer int) Subscription {
	ch := make(chan Chunk, buffer)
	s.mu.Lock()
	ch <- s.cache
	s.subs = append(s.subs, ch)
	s.mu.Unlock()
	return Subscription{ch, func() { s.unsubscribe(ch) }}
}

// unsubscribe removes ch from s.subs and closes it.
func (s *Stream) unsubscribe(ch chan Chunk) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// find & remove
	for i, c := range s.subs {
		if c == ch {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)
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
	s.finish()
}

// Closes the stream and all subscriber channels. Should be called ofter stream is done.
func (s *Stream) Close() {
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
