package httpmon

import "net/url"

type request struct {
	method    HttpMethod
	targetURL url.URL
}

func (r *request) Method() HttpMethod {
	return r.method
}

func (r *request) URL() string {
	return r.targetURL.String()
}

func NewRequest(method HttpMethod, URL string) (HttpTestRequest, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return &request{
		method:    method,
		targetURL: *u,
	}, nil
}
