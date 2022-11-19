package metrichandler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isGzipBody := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		isGzipResponse := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		if isGzipBody {
			gzReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gzReader.Close()

			r.Body = gzReader
		}

		if isGzipResponse {
			gzWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gzWriter.Close()

			w.Header().Set("Content-Encoding", "gzip")

			w = gzipWriter{
				ResponseWriter: w,
				Writer:         gzWriter,
			}
		}

		next.ServeHTTP(w, r)
	})
}
