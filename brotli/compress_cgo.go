package brotli

import (
	"io"

	"github.com/google/brotli/go/cbrotli"

	"github.com/acoshift/middleware"
)

// BrCompressor is the brotli compressor for compress middleware
var BrCompressor = middleware.CompressConfig{
	Skipper: middleware.DefaultSkipper,
	New: func() middleware.Compressor {
		return &brWriter{quality: 4}
	},
	Encoding:  "br",
	Vary:      middleware.DefaultCompressVary,
	Types:     middleware.DefaultCompressTypes,
	MinLength: middleware.DefaultCompressMinLength,
}

type brWriter struct {
	quality int
	*cbrotli.Writer
}

func (w *brWriter) Reset(p io.Writer) {
	w.Writer = cbrotli.NewWriter(p, cbrotli.WriterOptions{Quality: w.quality})
}
