package munnincore

import (
	"io"
	"net/http"
)

type HttpResponse struct {
	Response *http.Response
	Error    error
}

type HttpRequest struct {
	Title   string
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}
