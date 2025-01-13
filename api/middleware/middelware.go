package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

const millisecondConvertNumber = 1000000.0

func Logger(logger *zerolog.Logger, optSkipPaths ...[]string) func(next http.Handler) http.Handler {
	skipPaths := map[string]struct{}{} // skipPaths set for skipping specified url path from logging
	if len(optSkipPaths) > 0 {
		for _, path := range optSkipPaths[0] {
			skipPaths[path] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Skip the logger if the path is in the skip list
			if len(skipPaths) > 0 {
				_, skip := skipPaths[r.URL.Path]
				if skip {
					next.ServeHTTP(ww, r)

					return
				}
			}

			log := logger.With().Logger()

			timeStart := time.Now()

			defer func() {
				timeFinish := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				log.Info().
					Str("type", "access").
					Fields(map[string]interface{}{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(timeFinish.Sub(timeStart).Nanoseconds()) / millisecondConvertNumber,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
