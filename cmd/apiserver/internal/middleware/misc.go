package middleware

import (
	"net/http/httputil"
	"time"

	"authkey/cmd/apiserver/internal/types"
	"authkey/pkg/lib"

	"github.com/meilihao/water"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func DebugReq() water.HandlerFunc {
	return func(c *water.Context) {
		startTime := time.Now()
		dump, _ := httputil.DumpRequest(c.Request, true)

		c.Next()

		fs := make([]attribute.KeyValue, 0, 4)
		fs = append(fs,
			attribute.String("req", string(dump)),
			attribute.Int("status", c.Status()),
			attribute.String("duration", time.Since(startTime).String()),
		)

		lib.SpanLog(c.Request.Context(), trace.SpanFromContext(c.Request.Context()), zap.DebugLevel, "req", fs...)
	}
}

func IsAdmin() water.HandlerFunc {
	return func(c *water.Context) {
		tk := GetJWT(c)

		if tk.Role&types.RoleAdmin == 0 {
			c.Abort(401)
			return
		}

		c.Next()
	}
}
