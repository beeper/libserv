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
	Dial:                (&net.Dialer{Timeout: DefaultTimeout}).Dial,
	TLSHandshakeTimeout: DefaultTimeout,
	ForceAttemptHTTP2:   false, // never use HTTP/2
}

// DefaultClient is a http.Client using the DefaultTransport.
var DefaultClient = &http.Client{
	Transport: DefaultTransport,
}
