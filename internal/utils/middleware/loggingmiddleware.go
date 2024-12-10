package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

type loggingMiddleware struct {
	Next   http.RoundTripper
	Logger *slog.Logger
}

func (m *loggingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoints := req.Method + " " + req.URL.Host + req.URL.Path
	if req.URL.RawQuery != "" {
		endpoints += "?" + req.URL.RawQuery
	}

	if req.Body != nil {
		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		m.Logger.Info("Elasticsearch request", slog.String("endpoints", endpoints),
			slog.String("request body", string(body)),
		)
	} else {
		m.Logger.Info("Elasticsearch request", slog.String("endpoints", endpoints))
	}

	res, err := m.Next.RoundTrip(req)
	if res.Body != nil {
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		m.Logger.Info("response", slog.String("endpoints", endpoints),
			slog.String("response status", res.Status),
			slog.String("response body", string(body)),
		)
		res.Body = io.NopCloser(bytes.NewReader(body))
	} else {
		m.Logger.Info("response", slog.String("endpoints", endpoints),
			slog.String("response status", res.Status),
		)
	}

	return res, err
}

func NewLoggingMiddleware(logger *slog.Logger) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return &loggingMiddleware{
			Next:   next,
			Logger: logger,
		}
	}
}
