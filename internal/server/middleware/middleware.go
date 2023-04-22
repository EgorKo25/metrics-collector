package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	config "github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"log"
	"net"
	"net/http"
)

type Compressor struct {
}

func NewCompressor() *Compressor {
	return &Compressor{}
}

func (c *Compressor) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed write data to middleware temporary buffer: %v", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed middleware data: %v", err)
	}

	return b.Bytes(), nil
}

func (c *Compressor) Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed init decompress writer: %v", err)
	}

	var b bytes.Buffer

	if _, err = b.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	if err = r.Close(); err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}

type Middle struct {
	ipNet *net.IPNet
	cfg   *config.ConfigurationServer
}

func NewMiddle(cfg *config.ConfigurationServer) *Middle {
	return &Middle{
		cfg: cfg,
	}
}

func (m *Middle) IpChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error

		ipStr := r.Header.Get("X-Real-IP")

		if m.cfg.TrustedSubnet == "" {
			next.ServeHTTP(w, r)
			w.WriteHeader(500)
			return
		}

		_, m.ipNet, err = net.ParseCIDR(m.cfg.TrustedSubnet)
		if err != nil {
			log.Printf("%s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if m.ipNet.Contains(net.ParseIP(ipStr)) {
			next.ServeHTTP(w, r)
		}

		w.WriteHeader(http.StatusForbidden)
		return
	})
}
