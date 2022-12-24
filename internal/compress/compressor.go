package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
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
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
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
