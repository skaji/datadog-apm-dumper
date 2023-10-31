package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	stdout = os.Stdout
	stderr = os.Stderr
)

type Span struct {
	// https://github.com/DataDog/dd-trace-go/blob/v1/ddtrace/tracer/span.go#L51
	Name     string             `msgpack:"name" json:"name"`
	Service  string             `msgpack:"service" json:"service"`
	Resource string             `msgpack:"resource" json:"resource"`
	Type     string             `msgpack:"type" json:"type"`
	Start    int64              `msgpack:"start" json:"start"`
	Duration int64              `msgpack:"duration" json:"duration"`
	Meta     map[string]string  `msgpack:"meta,omitempty" json:"meta,omitempty"`
	Metrics  map[string]float64 `msgpack:"metrics,omitempty" json:"metrics,omitempty"`
	SpanID   uint64             `msgpack:"span_id" json:"span_id"`
	TraceID  uint64             `msgpack:"trace_id" json:"trace_id"`
	ParentID uint64             `msgpack:"parent_id" json:"parent_id"`
	Error    int32              `msgpack:"error" json:"error"`
}

const apmPath = "/v0.4/traces"

func apmServer(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != apmPath {
			w.WriteHeader(200)
			return
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintln(stderr, err)
			w.WriteHeader(400)
			return
		}
		var spans [][]Span
		if err := msgpack.Unmarshal(b, &spans); err != nil {
			d, _ := httputil.DumpRequest(r, true)
			fmt.Fprintln(stderr, "=== ERROR ===", string(d))
			w.WriteHeader(400)
			return
		}
		fmt.Fprintln(stderr, "---", time.Now().Format(time.RFC3339), r.RemoteAddr)
		out, _ := json.MarshalIndent(spans, "", "  ")
		for _, line := range strings.Split(string(out), "\n") {
			if len(line) > 0 {
				fmt.Fprintln(stdout, line)
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
		_ = addr
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		// fmt.Fprintln(stdout, "\033[1;34mstat\033[m| ---", addr)
		for _, line := range strings.Split(string(buf[:n]), "\n") {
			if len(line) > 0 {
				// fmt.Fprintln(stdout, "\033[1;34mstat\033[m|", line)
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
