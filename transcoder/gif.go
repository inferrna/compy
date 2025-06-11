package transcoder

import (
	"errors"
	"image/gif"
	"io"
	"net/http"

	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
)

type Gif struct{}

func (t *Gif) Transcode(w *proxy.ResponseWriter, r io.Reader, headers http.Header) error {
	img, err := gif.Decode(r)
	if err != nil {
		return err
	}
	if SupportsWebP(headers) {
		w.Header().Set("Content-Type", "image/webp")
		options := webp.Options{
			Lossless: true,
		}
		if err = webp.Encode(w, img, &options); err != nil {
			return err
		}
	} else {
		if err = gif.Encode(w, img, nil); err != nil {
			return err
		}
	}

	options := webp.Options{
		Lossless: true,
	}
	gifopts := gif.Options{
		NumColors: 16,
	}
	if webperr := webp.Encode(w, img, &options); webperr != nil {
		if giferr := gif.Encode(w, img, &gifopts); giferr != nil {
			return errors.Join(webperr, giferr)
		}
	} else {
		w.Header().Set("Content-Type", "image/webp")
	}

	return nil
}
