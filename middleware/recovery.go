package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/alex-cos/restc"
)

type RecoveryMiddleware struct {
	logger *slog.Logger
}

func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{logger: slog.Default()}
}

func NewRecoveryMiddlewareWith(logger *slog.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{logger: logger}
}

func (rm *RecoveryMiddleware) Handler() restc.Middleware {
	return func(req *restc.Request, next restc.HandlerFunc) (*restc.Response, error) {
		var (
			resp *restc.Response
			err  error
		)

		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic recovered: %v", r)
					rm.logger.Error("request panicked",
						slog.String("method", req.Method()),
						slog.String("url", req.URL()),
						slog.Any("panic", r),
						slog.String("stack", string(debug.Stack())),
					)
				}
			}()
			resp, err = next(req)
		}()

		return resp, err
	}
}
