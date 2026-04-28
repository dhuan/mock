package tests_e2e

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_BaseApi_SSE_IsStreamedWithoutWaitingForUpstreamClose(t *testing.T) {
	upstreamHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected upstream response writer to support flushing")
		}

		_, err := fmt.Fprint(w, "data: first event\n\n")
		if err != nil {
			t.Fatalf("failed to write first SSE event: %v", err)
		}
		flusher.Flush()

		time.Sleep(2 * time.Second)
	})
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate upstream listener: %v", err)
	}
	upstreamServer := &http.Server{Handler: upstreamHandler}
	go func() {
		_ = upstreamServer.Serve(listener)
	}()
	defer func() {
		_ = upstreamServer.Close()
	}()

	state := NewState()
	killMock, serverOutput, _, _ := RunMockBg(
		state,
		fmt.Sprintf("serve -p {{TEST_E2E_PORT}} --base 'http://%s'", listener.Addr().String()),
		nil,
		true,
		nil,
	)
	defer killMock()

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("http://localhost:%d/sse", state.Port),
		nil,
	)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		t.Fatalf("mock did not start streaming SSE before the upstream stayed open: %v\n\nServer output:\n%s", err, serverOutput.String())
	}
	defer response.Body.Close()

	if got := response.Header.Get("Content-Type"); !strings.HasPrefix(got, "text/event-stream") {
		t.Fatalf("expected SSE content type, got %q", got)
	}

	reader := bufio.NewReader(response.Body)
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read first SSE line: %v", err)
	}

	if firstLine != "data: first event\n" {
		t.Fatalf("unexpected first SSE line: %q", firstLine)
	}
}
