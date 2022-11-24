package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetricList(t *testing.T) {
	type args struct {
		MetricList  *map[string]handlers.Gauge
		CounterList *map[string]handlers.Counter
	}
	tests := []struct {
		name     string
		args     args
		url      string
		wantCode int
	}{
		{
			name: "test code 200",
			args: args{
				MetricList:  &map[string]handlers.Gauge{},
				CounterList: &map[string]handlers.Counter{},
			},
			url:      "http://127.0.0.1:8080/update/gauge/Alloc/12",
			wantCode: 200,
		},
		{
			name: "test code 400",
			args: args{
				MetricList:  &map[string]handlers.Gauge{},
				CounterList: &map[string]handlers.Counter{},
			},
			url:      "/update/gauge/Alloc/qdqw",
			wantCode: 400,
		},
		{
			name: "test code 400",
			args: args{
				MetricList:  &map[string]handlers.Gauge{},
				CounterList: &map[string]handlers.Counter{},
			},
			url:      "/update/counter/Alloc/qdqw",
			wantCode: 400,
		},
		{
			name: "test code 501",
			args: args{
				MetricList:  &map[string]handlers.Gauge{},
				CounterList: &map[string]handlers.Counter{},
			},
			url:      "/update/khjhk/Alloc/125",
			wantCode: 501,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(handlers.GetMetricList(tt.args.MetricList, tt.args.CounterList))
			h.ServeHTTP(w, req)
			resp := w.Result()

			if statusCode := resp.StatusCode; statusCode != tt.wantCode {
				t.Errorf("want %d, got %d", tt.wantCode, statusCode)
			}
		})
	}
}
