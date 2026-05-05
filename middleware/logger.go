package middleware

import (
	"log/slog"
	"time"

	"github.com/alex-cos/restc"
)

type MiddlewareLogger struct {
	logger *slog.Logger
}

func NewMiddlewareLogger() *MiddlewareLogger {
	return &MiddlewareLogger{
		logger: slog.Default(),
	}
}

func NewMiddlewareLoggerWith(logger *slog.Logger) *MiddlewareLogger {
	return &MiddlewareLogger{
		logger: logger,
	}
}

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
