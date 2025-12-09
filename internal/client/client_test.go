package client

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLinkAvailable(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		handlerFunc   http.HandlerFunc
		wantAvailable bool
		wantErr       bool
		errContains   string
	}{
		{
			name:          "Status 200 OK returns true",
			statusCode:    http.StatusOK,
			wantAvailable: true,
			wantErr:       false,
		},
		{
			name:          "Status 301 Redirect returns true",
			statusCode:    http.StatusMovedPermanently,
			wantAvailable: true,
			wantErr:       false,
		},
		{
			name:          "Status 404 Not Found returns false",
			statusCode:    http.StatusNotFound,
			wantAvailable: false,
			wantErr:       false,
		},
		{
			name:          "Status 500 Internal Server Error returns false",
			statusCode:    http.StatusInternalServerError,
			wantAvailable: false,
			wantErr:       false,
		},
		{
			// Special case for network errors (simulated by closing the connection)
			name: "Network error returns error",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				// Hijack the connection and close it to simulate a network error
				hj, ok := w.(http.Hijacker)
				require.True(t, ok, "webserver doesn't support hijacking")
				conn, _, err := hj.Hijack()
				require.NoError(t, err)
				func() {
					if closeErr := conn.Close(); closeErr != nil {
						log.Printf("cannot close response body: %v", closeErr)
					}
				}()

			},
			wantAvailable: false,
			wantErr:       true,
			errContains:   "network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodHead, r.Method)

				if tt.handlerFunc != nil {
					// Use custom handler for network error cases
					tt.handlerFunc(w, r)
					return
				}
				// Default handler: set the status code
				w.WriteHeader(tt.statusCode)
				// Need to write something to make Go's client happy during HEAD request mocks
				if _, err := io.WriteString(w, "body content"); err != nil {
					t.Logf("Error writing mock response body: %v", err)
				}
			}))
			defer server.Close()

			// Initialize the client under test
			// We create a new client that uses our test server's URL.
			client := NewCustomHTTPClient()
			// Crucially, we disable automatic redirects in the test client,
			// otherwise the httptest server will follow redirects automatically
			// and we can't test our 301/302 logic correctly.
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Call the function under test using the test server's URL
			available, err := client.IsLinkAvailable(ctx, server.URL)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.wantAvailable, available)
		})
	}
}
