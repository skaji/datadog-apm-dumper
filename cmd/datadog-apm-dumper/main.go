package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vmihailenco/msgpack"
)

type Span struct {
	// https://github.com/DataDog/dd-trace-go/blob/v1/ddtrace/tracer/span.go#L51
	Name     string             `json:"name"`
	Service  string             `json:"service"`
	Resource string             `json:"resource"`
	Type     string             `json:"type"`
	Start    int64              `json:"start"`
	Duration int64              `json:"duration"`
	Meta     map[string]string  `json:"meta,omitempty"`
	Metrics  map[string]float64 `json:"metrics,omitempty"`
	SpanID   uint64             `json:"span_id"`
	TraceID  uint64             `json:"trace_id"`
	ParentID uint64             `json:"parent_id"`
	Error    int32              `json:"error"`
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var spans [][]Span
		if err := msgpack.NewDecoder(r.Body).UseJSONTag(true).Decode(&spans); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}
		b, _ := json.MarshalIndent(spans, "", "  ")
		fmt.Println(string(b))
		w.WriteHeader(200)
	}
	fmt.Println("Accepting connections at http://localhost:8126/")
	http.ListenAndServe("localhost:8126", http.HandlerFunc(handler))
}
