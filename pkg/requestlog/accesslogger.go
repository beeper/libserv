package requestlog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

const MaxRequestSizeLog = 4 * 1024
const MaxStringRequestSizeLog = MaxRequestSizeLog / 2

func AccessLogger(logOptions bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := hlog.FromRequest(r)

			crw := &CountingResponseWriter{
				ResponseWriter: w,
				ResponseLength: -1,
				StatusCode:     -1,
			}

			start := time.Now()
			defer func() {
				requestDuration := time.Since(start)

				var requestLog *zerolog.Event
				if crw.StatusCode >= 500 {
					requestLog = log.Error()
				} else if crw.StatusCode >= 400 {
					requestLog = log.Warn()
				} else {
					requestLog = log.Info()
				}

				// recover and log the error
				if rcvr := recover(); rcvr != nil && rcvr != http.ErrAbortHandler {
					crw.StatusCode = http.StatusInternalServerError
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(crw.StatusCode)

					resp := struct {
						Errcode   string `json:"errcode"`
						Error     string `json:"error"`
						RequestID string `json:"request_id"`
						Time      string `json:"time"`
					}{
						Errcode: "COM.BEEPER.PANIC",
						Error:   "Internal Server Error",
						Time:    time.Now().UTC().Format(time.RFC3339),
					}

					if requestID, ok := hlog.IDFromRequest(r); ok {
						resp.RequestID = requestID.String()
					}

					json.NewEncoder(w).Encode(&resp)

					requestLog = log.Error()
					requestLog.Interface("panic", rcvr)
					requestLog.Bytes("stack", debug.Stack())
				}

				if r.Method == http.MethodOptions && !logOptions {
					return
				}

				if userAgent := r.UserAgent(); userAgent != "" {
					requestLog.Str("user_agent", userAgent)
				}
				if referer := r.Referer(); referer != "" {
					requestLog.Str("referer", referer)
				}
				remoteAddr := r.RemoteAddr

				requestLog.Str("remote_addr", remoteAddr)
				requestLog.Str("method", r.Method)
				requestLog.Str("proto", r.Proto)
				requestLog.Int64("request_length", r.ContentLength)
				requestLog.Str("host", r.Host)
				requestLog.Str("request_uri", r.RequestURI)
				if r.Method != http.MethodGet && r.Method != http.MethodHead {
					requestLog.Str("request_content_type", r.Header.Get("Content-Type"))
					if crw.RequestBody != nil {
						logRequestMaybeJSON(requestLog, "request_body", crw.RequestBody.Bytes())
					}
				}

				// response
				requestLog.Int64("request_time_ms", requestDuration.Milliseconds())
				requestLog.Int("status_code", crw.StatusCode)
				requestLog.Int("response_length", crw.ResponseLength)
				requestLog.Str("response_content_type", crw.Header().Get("Content-Type"))
				if crw.ResponseBody != nil {
					logRequestMaybeJSON(requestLog, "response_body", crw.ResponseBody.Bytes())
				}

				// don't log successful health requests
				if r.URL.Path == "/health" && crw.StatusCode == http.StatusNoContent {
					return
				}

				requestLog.Msg("Access")
			}()

			next.ServeHTTP(crw, r)
		})
	}
}

func logRequestMaybeJSON(evt *zerolog.Event, key string, data []byte) {
	data = removeNewlines(data)
	if json.Valid(data) {
		evt.RawJSON(key, data)
	} else {
		// Logging as a string will create lots of escaping and it's not valid json anyway, so cut off a bit more
		if len(data) > MaxStringRequestSizeLog {
			data = data[:MaxStringRequestSizeLog]
		}
		evt.Bytes(key+"_invalid", data)
	}
}

func removeNewlines(data []byte) []byte {
	data = bytes.TrimSpace(data)
	if bytes.ContainsRune(data, '\n') {
		data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{})
		data = bytes.ReplaceAll(data, []byte{'\r'}, []byte{})
	}
	return data
}
