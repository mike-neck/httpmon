package httpmon

type DefaultHttpClient struct {
}

func (client *DefaultHttpClient) Run(request HttpRequest) (HttpTest, error) {
	return &DefaultHttpTest{
		Status: 0,
		Header: nil,
	}, nil
}
