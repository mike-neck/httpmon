package httpmon

import (
	"bytes"
	"io"
	"net/url"
)

type getRequest struct {
	method    HttpMethod
	targetURL url.URL
	headers   []HttpHeader
}

func (r *getRequest) URL() URL {
	return URL(r.targetURL.String())
}

func (r *getRequest) Body() io.Reader {
	empty := make([]byte, 0)
	return bytes.NewReader(empty)
}

func (r *getRequest) Headers() HttpHeaders {
	hs := make(HttpHeaders, 0)
	for _, h := range r.headers {
		if _, ok := hs[h.Name]; !ok {
			hs[h.Name] = make([]string, 0)
		}
		hs[h.Name] = append(hs[h.Name], h.Value)
	}
	return hs
}

func (r *getRequest) Method() HttpMethod {
	return r.method
}

func NewGetRequestDetails(method HttpMethod, URL string, headers ...HttpHeader) (HttpRequestDetails, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return &getRequest{
		method:    method,
		targetURL: *u,
		headers:   headers,
	}, nil
}
