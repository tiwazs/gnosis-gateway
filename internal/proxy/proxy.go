package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(rawUrl string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	proxy_instance := httputil.NewSingleHostReverseProxy(target)
	proxy_instance.FlushInterval = -1
	original := proxy_instance.Director

	proxy_instance.Director = func(req *http.Request) {
		original(req)
		req.Host = target.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	}
	// FastAPI/Express also set CORS. Ingress sets it too. Duplicate
	// Access-Control-Allow-Origin makes the browser drop the response.
	proxy_instance.ModifyResponse = func(res *http.Response) error {
		stripCORS(res.Header)
		return nil
	}

	return proxy_instance, nil
}

func stripCORS(header http.Header) {
	header.Del("Access-Control-Allow-Origin")
	header.Del("Access-Control-Allow-Credentials")
	header.Del("Access-Control-Allow-Headers")
	header.Del("Access-Control-Allow-Methods")
	header.Del("Access-Control-Expose-Headers")
	header.Del("Access-Control-Max-Age")
	header.Del("Access-Control-Request-Headers")
	header.Del("Access-Control-Request-Method")
}
