package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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

func apmServer(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	handler := func(w http.ResponseWriter, r *http.Request) {
		var spans [][]Span
		if err := msgpack.NewDecoder(r.Body).UseJSONTag(true).Decode(&spans); err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}
		fmt.Println("\033[1;31mapm\033[m | ---", r.RemoteAddr)
		b, _ := json.MarshalIndent(spans, "", "  ")
		for _, line := range strings.Split(string(b), "\n") {
			if len(line) > 0 {
				fmt.Println("\033[1;31mapm\033[m |", line)
			}
		}
		w.WriteHeader(200)
	}
	svc := &http.Server{Addr: addr, Handler: http.HandlerFunc(handler)}
	done := make(chan error, 1)
	go func() {
		<-ctx.Done()
		if err := svc.Shutdown(context.Background()); err != nil {
			done <- err
		}
		close(done)
	}()
	if err := svc.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		cancel()
		<-done
		return err
	}
	return <-done
}

func statsdServer(ctx context.Context, addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	var buf [1500]byte
	for {
		n, addr, err := conn.ReadFrom(buf[:])
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		fmt.Println("\033[1;34mstat\033[m| ---", addr)
		for _, line := range strings.Split(string(buf[:n]), "\n") {
			if len(line) > 0 {
				fmt.Println("\033[1;34mstat\033[m|", line)
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("start apm server at http://localhost:8126")
		if err := apmServer(ctx, "localhost:8126"); err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		fmt.Println("start statsd server at udp://localhost:8125")
		if err := statsdServer(ctx, "localhost:8125"); err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
	fmt.Println("finish")
}
