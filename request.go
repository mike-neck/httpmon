package httpmon

type getTest struct {
	targetURL string
}

func (g *getTest) Method() HttpMethod {
	return "GET"
}

func (g *getTest) URL() string {
	return g.targetURL
}

func NewGetRequest(URL string) HttpTestRequest {
	return &getTest{targetURL: URL}
}
