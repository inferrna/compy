package transcoder

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/barnacs/compy/proxy"
	"github.com/chai2010/webp"
	"github.com/pixiv/go-libjpeg/jpeg"
)

type Jpeg struct {
	decOptions *jpeg.DecoderOptions
	encOptions *jpeg.EncoderOptions
}

func NewJpeg(quality int) *Jpeg {
	return &Jpeg{
		decOptions: &jpeg.DecoderOptions{},
		encOptions: &jpeg.EncoderOptions{
			Quality:        quality,
			OptimizeCoding: true,
		},
	}
}

func (t *Jpeg) Transcode(w *proxy.ResponseWriter, r io.Reader, headers http.Header) error {
	img, err := jpeg.Decode(r, t.decOptions)
	if err != nil {
		return err
	}

	encOptions := t.encOptions
	qualityString := headers.Get("X-Compy-Quality")
	if qualityString != "" {
		if quality, err := strconv.Atoi(qualityString); err != nil {
			encOptions.Quality = 73
		} else {
			encOptions.Quality = quality
		}
	}

	options := webp.Options{
		Lossless: false,
		Quality:  float32(encOptions.Quality),
	}
	if webperr := webp.Encode(w, img, &options); webperr != nil {
		if jpegerr := jpeg.Encode(w, img, encOptions); jpegerr != nil {
			return errors.Join(jpegerr, jpegerr)
		}
	} else {
		w.Header().Set("Content-Type", "image/webp")
	}
	return nil
}
