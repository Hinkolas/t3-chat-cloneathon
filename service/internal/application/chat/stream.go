package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) OpenStream(w http.ResponseWriter, r *http.Request) {

	streamID := mux.Vars(r)["id"]

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.log.Debug("streaming not supported", "stream_id", streamID)
		http.Error(w, "streaming not supported", http.StatusNotFound)
		return
	}

	// Subscribe to the stream
	sub, ok := s.sp.Subscribe(streamID)
	if !ok {
		s.log.Debug("stream not found", "stream_id", streamID)
		http.Error(w, "stream not found", http.StatusNotFound)
		return
	}
	defer sub.Cancel()

	w.WriteHeader(http.StatusOK)
	for {
		select {
		case <-r.Context().Done():
			// client disconnected
			s.log.Debug("stream: client closed connection", "stream_id", streamID)
			return

		case chunk, more := <-sub.Read():
			if !more {
				// publisher closed the stream
				s.log.Debug("stream: provider closed the stream", "stream_id", streamID)
				// write the SSE event
				if _, err := fmt.Fprint(w,
					"event: message_end\n",
					"data: {}\n\n",
				); err != nil {
					s.log.Debug("stream: write failed", "err", err)
					return
				}
				flusher.Flush()
				return
			}
			// write the SSE event
			if _, err := fmt.Fprint(w,
				"event: message_delta\n",
				"data: ",
			); err != nil {
				s.log.Debug("stream: write failed", "err", err)
				return
			}
			if err := json.NewEncoder(w).Encode(chunk); err != nil {
				s.log.Debug("stream: json encoding failed", "err", err)
				return
			}
			if _, err := fmt.Fprint(w, "\n"); err != nil {
				s.log.Debug("stream: write failed", "err", err)
				return
			}
			flusher.Flush()
		}
	}

}

// TODO: Consider creating a database table to keep track of ongoing streams
func (s *Service) CancelStream(w http.ResponseWriter, r *http.Request) {
	streamID := mux.Vars(r)["id"]
	s.sp.Cancel(streamID)
}
