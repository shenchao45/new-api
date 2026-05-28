package relay

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseBodyRecorder wraps gin.ResponseWriter to capture bytes written to the client.
// It implements relaycommon.CapturedWriter via the Captured() method.
type ResponseBodyRecorder struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func NewResponseBodyRecorder(w gin.ResponseWriter) *ResponseBodyRecorder {
	return &ResponseBodyRecorder{ResponseWriter: w, buf: &bytes.Buffer{}}
}

func newResponseBodyRecorder(w gin.ResponseWriter) *ResponseBodyRecorder {
	return NewResponseBodyRecorder(w)
}

func (r *ResponseBodyRecorder) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *ResponseBodyRecorder) WriteString(s string) (int, error) {
	r.buf.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

func (r *ResponseBodyRecorder) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseBodyRecorder) WriteHeaderNow() {
	r.ResponseWriter.WriteHeaderNow()
}

func (r *ResponseBodyRecorder) Status() int {
	return r.ResponseWriter.Status()
}

func (r *ResponseBodyRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *ResponseBodyRecorder) Captured() string {
	return r.buf.String()
}
