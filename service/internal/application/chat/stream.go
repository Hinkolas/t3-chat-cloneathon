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

	w.WriteHeader(http.StatusOK)

	// Subscribe to the stream
	sub, ok := s.sp.Subscribe(streamID)
	if !ok {
		s.log.Debug("stream not found", "stream_id", streamID)
		http.Error(w, "stream not found", http.StatusNotFound)
		return
	}

	// Flush all received chunks to the client
	for c := range sub {
		fmt.Fprint(w, "event: message_delta\ndata: ")
		if err := json.NewEncoder(w).Encode(c); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		fmt.Fprint(w, "\n")
		flusher.Flush()
	}
	fmt.Fprintf(w, "event: message_end\ndata: {\"done\":true}\n\n")
	flusher.Flush()

	s.log.Debug("streaming completed successfully", "stream_id", streamID)

}

func (s *Service) CancelStream(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented yet", http.StatusNotImplemented)
}
