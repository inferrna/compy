package transcoder

import (
	"io"
	"net/http"

	"github.com/barnacs/compy/proxy"
)

type Identity struct{}

func (i *Identity) Transcode(w *proxy.ResponseWriter, r io.Reader, headers http.Header) error {
	return w.ReadFrom(r)
}
