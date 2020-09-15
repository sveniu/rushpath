package main

import (
	"bytes"
	"net/http"
	"net/http/httputil"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	http.DefaultTransport = wrapRoundTripper(http.DefaultTransport)
}

func wrapRoundTripper(in http.RoundTripper) http.RoundTripper {
	return &loggingRoundTripper{inner: in}
}

type loggingRoundTripper struct {
	inner http.RoundTripper
}

func (d *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	d.dumpRequest(req)
	res, err := d.inner.RoundTrip(req)
	d.dumpResponse(res)
	return res, err
}

func (d *loggingRoundTripper) dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		log.Trace().
			Err(err).
			Msg("error dumping request")
		return
	}
	headersBody := bytes.Split(dump, []byte("\r\n\r\n"))
	headers := bytes.Split(headersBody[0], []byte("\r\n"))

	headerDict := zerolog.Dict()
	for _, header := range headers[1:] {
		keyVal := bytes.Split(header, []byte(": "))
		headerDict.Str(string(keyVal[0]), string(keyVal[1]))
	}

	log.Trace().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Dict("headers", headerDict).
		Bytes("body", headersBody[1]).
		Msg("http request")
}

func (d *loggingRoundTripper) dumpResponse(r *http.Response) {
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Trace().
			Err(err).
			Msg("error dumping response")
		return
	}
	headersBody := bytes.Split(dump, []byte("\r\n\r\n"))
	headers := bytes.Split(headersBody[0], []byte("\r\n"))

	headerDict := zerolog.Dict()
	for _, header := range headers[1:] {
		keyVal := bytes.Split(header, []byte(": "))
		headerDict.Str(string(keyVal[0]), string(keyVal[1]))
	}

	log.Trace().
		Int("status_code", r.StatusCode).
		Dict("headers", headerDict).
		Bytes("body", headersBody[1]).
		Msg("http response")
}
