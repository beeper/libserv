package httputil

import (
	"net"
	"net/http"
	"time"
)

const DefaultTimeout = time.Second * 15

// DeafultTransport is a http.Transport with a 15 second dial and TLS handshake
// timeout and HTTP/2 support disabled.
var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   DefaultTimeout,
		KeepAlive: 30 * time.Second,
	}).Dial,
	ForceAttemptHTTP2:     false, // never use HTTP/2
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   DefaultTimeout,
	ExpectContinueTimeout: 1 * time.Second,
}

// DefaultClient is a http.Client using the DefaultTransport.
var DefaultClient = &http.Client{
	Transport: DefaultTransport,
}
