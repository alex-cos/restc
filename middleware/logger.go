package middleware

import (
	"log/slog"
	"time"

	"github.com/alex-cos/restc"
)

// MiddlewareLogger is a middleware that logs HTTP requests and responses.
type MiddlewareLogger struct {
	logger *slog.Logger
}

// NewMiddlewareLogger creates a new MiddlewareLogger with the default log.Logger.
func NewMiddlewareLogger() *MiddlewareLogger {
	return &MiddlewareLogger{
		logger: slog.Default(),
	}
}

// NewMiddlewareLoggerWith creates a new MiddlewareLogger with the given slog.Logger.
func NewMiddlewareLoggerWith(logger *slog.Logger) *MiddlewareLogger {
	return &MiddlewareLogger{
		logger: logger,
	}
}

// Handler returns a Middleware that logs HTTP requests and responses.
// It logs method, URL, status code, and duration.
// For errors, it logs at error level.
// For 4xx/5xx responses, it logs at warning level.
// For successful responses, it logs at info level.
func (l *MiddlewareLogger) Handler() restc.Middleware {
	return func(req *restc.Request, next restc.HandlerFunc) (*restc.Response, error) {
		start := time.Now()

		resp, err := next(req)

		attrs := []any{
			slog.String("method", req.Method()),
			slog.String("url", req.URL()),
			slog.Int("status", resp.StatusCode()),
			slog.Duration("duration", time.Since(start)),
		}

		switch {
		case err != nil:
			attrs = append(attrs, slog.String("error", err.Error()))
			l.logger.Error("HTTP request failed", attrs...)
		case resp.IsError():
			l.logger.Warn("HTTP error response", attrs...)
		default:
			l.logger.Info("HTTP request completed", attrs...)
		}

		return resp, err
	}
}
