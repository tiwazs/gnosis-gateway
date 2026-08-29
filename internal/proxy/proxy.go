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

	return proxy_instance, nil
}