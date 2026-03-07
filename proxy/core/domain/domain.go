package domain

type HttpPackage struct {
	StatusLine string
	Headers    map[string]string
	Body       string
}
