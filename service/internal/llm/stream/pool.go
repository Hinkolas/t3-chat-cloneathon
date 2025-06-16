package stream

import (
	"fmt"
	"sync"
)

// Stream pool for managing active streams
type StreamPool struct {
	streams map[string]*Stream
	mu      sync.RWMutex
}

func NewStreamPool() *StreamPool {
	return &StreamPool{streams: make(map[string]*Stream)}
}

func (p *StreamPool) Add(id string, s *Stream) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.streams[id] = s
}

func (p *StreamPool) Get(id string) (*Stream, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok := p.streams[id]
	return s, ok
}

func (p *StreamPool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.streams, id)
}

func (p *StreamPool) Cancel(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if s, ok := p.streams[id]; ok {
		s.Fail(fmt.Errorf("stream cancelled by user"))
		delete(p.streams, id)
	}
}

func (p *StreamPool) Subscribe(id string) (sub *Subscription, ok bool) {
	var s *Stream
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok = p.streams[id]
	if ok {
		return s.Subscribe(10), true // TODO: find a buffer size that works well
	}
	return nil, false
}
