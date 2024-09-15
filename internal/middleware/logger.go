package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// func New(log *slog.Logger) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		log := log.With(
// 			slog.String("component", "middleware/logger"),
// 		)

// 		log.Info("logger middleware enabled")

// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			entry := log.With(
// 				slog.String("method", r.Method),
// 				slog.String("path", r.URL.Path),
// 				slog.String("remote_addr", r.RemoteAddr),
// 				slog.String("user_agent", r.UserAgent()),
// 				slog.String("request_id", middleware.GetReqID(r.Context())),
// 			)
// 			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

// 			t1 := time.Now()
// 			defer func() {
// 				entry.Info("request completed",
// 					slog.Int("status", ww.Status()),
// 					slog.Int("bytes", ww.BytesWritten()),
// 					slog.String("duration", time.Since(t1).String()),
// 				)
// 			}()

// 			next.ServeHTTP(ww, r)
// 		}

// 		return http.HandlerFunc(fn)
// 	}
// }

type LoggerMiddlware struct {
	logger *slog.Logger
}

func NewLoggerMiddlware(l *slog.Logger) *LoggerMiddlware {
	return &LoggerMiddlware{
		logger: l,
	}
}

func (m *LoggerMiddlware) Proccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := m.logger.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
		start := time.Now()
		next.ServeHTTP(w, r)
		entry.Info("Request completed!", slog.Duration("duration", time.Since(start)))
	})
}
