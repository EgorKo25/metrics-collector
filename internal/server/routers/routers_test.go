package routers

import (
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/server/handlers"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {

	mem := storage.NewStorage()
	handler := handlers.NewHandler(mem)

	r := NewRouter(handler)
	ts := httptest.NewServer(r)
	defer ts.Close()

	mem.SetGaugeStat("Alloc", storage.Gauge(123), "gauge")

	statusCode, body := testRequest(t, ts, "GET", "/")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "> Alloc:  123.000000\n", body)

	statusCode, body = testRequest(t, ts, "GET", "/value/gauge/Alloc")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "123\n", body)

	statusCode, body = testRequest(t, ts, "POST", "/update/gauge/Alloc/123")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "", body)

}
func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(respBody)
}
