package transcoder

import (
	"errors"
	"image/png"
	"io"
	"net/http"

	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
)

type Png struct{}

func (t *Png) Transcode(w *proxy.ResponseWriter, r io.Reader, headers http.Header) error {
	img, err := png.Decode(r)
	if err != nil {
		return err
	}

	options := webp.Options{
		Lossless: true,
	}
	if webperr := webp.Encode(w, img, &options); webperr != nil {
		if pngerr := png.Encode(w, img); pngerr != nil {
			return errors.Join(webperr, pngerr)
		}
	} else {
		w.Header().Set("Content-Type", "image/webp")
	}
	return nil
}
