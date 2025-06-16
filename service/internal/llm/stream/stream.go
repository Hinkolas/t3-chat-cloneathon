package stream

import (
	"fmt"
	"slices"
	"sync"
)

type CloseFunc func([]Chunk) error

type Chunk struct {
	Thinking string `json:"thinking,omitempty"`
	Content  string `json:"content,omitempty"`
}

// Stream represents one ongoing streaming process.
type Stream struct {
	// protected by mu
	mu        sync.RWMutex
	wg        sync.WaitGroup
	closeOnce sync.Once

	chunks []Chunk // all chunks received so far
	done   bool    // true once the stream finishes
	err    error   // any terminal error

	pub  chan Chunk   // where callers Publish
	subs []chan Chunk // subscriber channels

	closeFunc CloseFunc // callback for cleanup after finish
}

// New returns a Stream that’s ready to Start() / Publish().
func New() *Stream {
	s := &Stream{
		pub:       make(chan Chunk),
		closeFunc: func([]Chunk) error { return nil },
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

func (s *Stream) OnClose(closeFunc func([]Chunk) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeFunc = closeFunc
}

// Publish sends one chunk into the stream. Blocks if no buffer on inCh.
func (s *Stream) Publish(c Chunk) {
	s.pub <- c
}

// Subscribe returns a channel on which the caller will receive _all_ past and future chunks
func (s *Stream) Subscribe(buffer int) <-chan Chunk {
	ch := make(chan Chunk, buffer)
	s.mu.Lock()
	// replay buffered chunks
	for _, c := range s.chunks {
		select {
		case ch <- c:
		default:
			// TODO: Handle subscriber not keeping up (e.g. cancel subscription)
			fmt.Println("subscriber is slow, blocking on replay")
		}
	}
	s.subs = append(s.subs, ch)
	s.mu.Unlock()
	return ch
}

// Wait blocks until the stream is done. It returns any error.
func (s *Stream) Wait() error {
	s.wg.Wait()
	return s.err
}

func (s *Stream) Fail(err error) {
	s.setError(err)
	close(s.pub)
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
	s.chunks = append(s.chunks, chunk)
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
		for _, ch := range s.subs {
			close(ch)
		}
		s.subs = nil
		if err := s.closeFunc(s.chunks); err != nil {
			s.setError(err)
		}
	})
}
