package health

import (
	"net"
	"net/http"

	"github.com/beeper/libserv/pkg/response"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		if ip := net.ParseIP(host); ip != nil && !ip.IsPrivate() {
			response.Forbidden(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
