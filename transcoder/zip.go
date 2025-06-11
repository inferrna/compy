package transcoder

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/barnacs/compy/proxy"
	brotlidec "gopkg.in/kothar/brotli-go.v0/dec"
	brotlienc "gopkg.in/kothar/brotli-go.v0/enc"
)

type Zip struct {
	proxy.Transcoder
	BrotliCompressionLevel int
	GzipCompressionLevel   int
	SkipCompressed         bool
}

func (t *Zip) Transcode(w *proxy.ResponseWriter, r io.Reader, headers http.Header) error {
	shouldBrotli := false
	shouldGzip := false
	gzipped := w.Header().Get("Content-Encoding") == "gzip"
	brotlied := w.Header().Get("Content-Encoding") == "br"
	shouldCompress := !(gzipped || brotlied)
	for _, v := range strings.Split(headers.Get("Accept-Encoding"), ", ") {
		switch strings.SplitN(v, ";", 2)[0] {
		case "br":
			shouldBrotli = true
		case "gzip":
			shouldGzip = true
		}
	}

	log.Printf("br == %v , gzip == %v", shouldBrotli, shouldGzip)

	// always gunzip if the client supports Brotli
	if gzipped && (shouldBrotli || !t.SkipCompressed) {
		gzr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer gzr.Close()
		r = gzr
		shouldCompress = true
		log.Printf("Found gzipped page")
	}

	if brotlied && !t.SkipCompressed {
		brr := brotlidec.NewBrotliReader(r)
		defer brr.Close()
		r = brr
		shouldCompress = true
	}

	if shouldBrotli && shouldCompress {
		params := brotlienc.NewBrotliParams()
		params.SetQuality(t.BrotliCompressionLevel)
		brw := brotlienc.NewBrotliWriter(params, w.Writer)
		defer brw.Close()
		w.Writer = brw
		w.Header().Set("Content-Encoding", "br")
		log.Printf("Compress to br")
	} else if shouldGzip && shouldCompress {
		gzw, err := gzip.NewWriterLevel(w.Writer, t.GzipCompressionLevel)
		if err != nil {
			return err
		}
		defer gzw.Close()
		w.Writer = gzw
		w.Header().Set("Content-Encoding", "gzip")
		log.Printf("Compress to gzip")
	}
	//return nil
	return t.Transcoder.Transcode(w, r, headers)
}
