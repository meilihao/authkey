package svc

import (
	"authkey/cmd/apiserver/internal/config"
	"authkey/pkg/lib/db"

	"github.com/meilihao/layer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	l      *layer.Layer
	tracer trace.Tracer
)

func Init() error {
	var err error

	tracer = otel.Tracer("svc")
	if l, err = db.InitMySQL2Layer(config.GlobalConfig.DB); err != nil {
		return err
	}

	return nil
}
