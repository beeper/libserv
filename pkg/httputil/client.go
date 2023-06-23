package httputil

import (
	"net"
	"net/http"
	"time"
)

const DefaultTimeout = time.Second * 15

// NoHTTP2Transport is a http.Transport with HTTP/2 support disabled.
var NoHTTP2Transport = &http.Transport{
	Dial:                (&net.Dialer{Timeout: DefaultTimeout}).Dial,
	TLSHandshakeTimeout: DefaultTimeout,
	ForceAttemptHTTP2:   false, // never use HTTP/2
}

// DefaultClient is a http.Client with a 15 second timeout with HTTP/2 support
// disabled.
var DefaultClient = &http.Client{
	Transport: NoHTTP2Transport,
	Timeout:   DefaultTimeout,
}
